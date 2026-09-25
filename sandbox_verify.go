package langsmith

import (
	"cmp"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	// SandboxUserTokenHeader carries the signed identity of a request made
	// through a LangSmith-login service URL.
	SandboxUserTokenHeader = "X-Langsmith-User-Token"
	// SandboxCallbackSignatureHeader carries the signature of a proxy callback
	// request.
	SandboxCallbackSignatureHeader = "X-LangSmith-Signature-JWT"

	sandboxCallbackSubject = "langsmith-sandbox-callback"
	sandboxJWKSPath        = "/.well-known/jwks.json"
	sandboxJWKSTTL         = 5 * time.Minute
	// Bounds refetches when a token names a kid the cached set lacks.
	sandboxJWKSMinRefresh = 30 * time.Second
	sandboxTokenLeeway    = 30 * time.Second
	sandboxJWKSMaxBytes   = 1 << 20
)

// SandboxTokenVerificationError reports a sandbox user token or proxy callback
// signature that failed verification.
type SandboxTokenVerificationError struct {
	Reason string
	Err    error
}

func (e *SandboxTokenVerificationError) Error() string {
	if e.Err != nil {
		return "sandbox token verification failed: " + e.Reason + ": " + e.Err.Error()
	}
	return "sandbox token verification failed: " + e.Reason
}

func (e *SandboxTokenVerificationError) Unwrap() error { return e.Err }

func verificationError(reason string, err error) error {
	return &SandboxTokenVerificationError{Reason: reason, Err: err}
}

// SandboxUser is the LangSmith user a service URL request was made by.
type SandboxUser struct {
	Subject   string
	Email     string
	Name      string
	ExpiresAt time.Time
}

// SandboxCallbackIdentity is the sandbox whose outbound request triggered a
// proxy callback.
type SandboxCallbackIdentity struct {
	TenantID       string `json:"tenant_id"`
	SandboxID      string `json:"sandbox_id"`
	OrganizationID string `json:"organization_id,omitempty"`
	LSUserID       string `json:"ls_user_id,omitempty"`
}

// SandboxCallbackRequest is the snapshot of the outbound request sent to
// callbacks configured with full_request.
type SandboxCallbackRequest struct {
	Method        string              `json:"method"`
	URL           string              `json:"url"`
	Scheme        string              `json:"scheme"`
	Host          string              `json:"host"`
	Path          string              `json:"path"`
	Query         string              `json:"query,omitempty"`
	Headers       map[string][]string `json:"headers"`
	Body          []byte              `json:"body_base64"`
	BodyTruncated bool                `json:"body_truncated"`
}

// SandboxCallback is a verified proxy callback payload.
type SandboxCallback struct {
	Host     string                  `json:"host"`
	Port     int                     `json:"port"`
	Identity SandboxCallbackIdentity `json:"identity"`
	Request  *SandboxCallbackRequest `json:"request,omitempty"`
}

// SandboxTokenVerifierOptions configures a [SandboxTokenVerifier].
type SandboxTokenVerifierOptions struct {
	// APIURL is the LangSmith API URL whose origin serves the JWKS. Defaults to
	// LANGSMITH_ENDPOINT.
	APIURL string
	// JWKSURL is the full JWKS URL; it overrides APIURL.
	JWKSURL string
	// HTTPClient fetches the JWKS. Defaults to a client with a 10s timeout.
	HTTPClient *http.Client
	// AllowInsecureJWKS allows fetching the JWKS over plain HTTP from a
	// non-loopback host. Anyone who can tamper with that traffic can forge
	// tokens this verifier accepts.
	AllowInsecureJWKS bool
}

// SandboxUserTokenOptions configures [SandboxTokenVerifier.VerifyUserToken].
type SandboxUserTokenOptions struct {
	// Audience is the service URL host the request was sent to, such as the
	// request's Host header. A full service URL is also accepted. Required.
	Audience string
	// Issuer, if set, is the LangSmith app URL the token must be issued by.
	Issuer string
}

// SandboxCallbackOptions configures [SandboxTokenVerifier.VerifyCallback].
type SandboxCallbackOptions struct {
	// Audience, if set, is called with each audience in the signature, which is
	// the callback URL as configured in the proxy config, and must return true
	// for at least one. Use [ExactAudience] to match one URL exactly.
	Audience func(aud string) bool
	// Issuer, if set, is the LangSmith OAuth issuer the signature must be
	// issued by.
	Issuer string
}

// ExactAudience returns an audience matcher that accepts only want.
func ExactAudience(want string) func(string) bool {
	return func(aud string) bool { return aud == want }
}

// SandboxTokenVerifier verifies tokens LangSmith signs for code running in or
// behind a sandbox: the X-Langsmith-User-Token header of service URL requests
// and the X-LangSmith-Signature-JWT header of proxy callbacks. Keys are fetched
// from LangSmith's JWKS endpoint and cached. It is safe for concurrent use.
//
// Code that knows it is running in a sandbox can instead trust the unsigned
// X-Langsmith-User-Id and X-Langsmith-User-Email headers on requests that
// arrive through the service URL, which the sandbox runtime strips from inbound
// requests and sets itself. TCP tunnels and other processes in the sandbox
// calling the port over localhost bypass that stripping.
type SandboxTokenVerifier struct {
	jwksURL string
	client  *http.Client
	// Serializes JWKS fetches so concurrent cache misses share one request.
	fetchSem chan struct{}

	mu        sync.Mutex
	keys      map[string]ed25519.PublicKey
	fetchedAt time.Time
}

// NewSandboxTokenVerifier returns a verifier. It fails if the JWKS URL is not
// HTTPS, unless the host is loopback or opts.AllowInsecureJWKS is set.
func NewSandboxTokenVerifier(opts SandboxTokenVerifierOptions) (*SandboxTokenVerifier, error) {
	jwksURL := opts.JWKSURL
	if jwksURL == "" {
		apiURL := cmp.Or(opts.APIURL, os.Getenv("LANGSMITH_ENDPOINT"), "https://api.smith.langchain.com")
		u, err := url.Parse(apiURL)
		if err != nil {
			return nil, fmt.Errorf("sandbox token verifier: parse API URL: %w", err)
		}
		jwksURL = u.Scheme + "://" + u.Host + sandboxJWKSPath
	}
	if err := checkSandboxJWKSURL(jwksURL, opts.AllowInsecureJWKS); err != nil {
		return nil, err
	}
	client := opts.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &SandboxTokenVerifier{
		jwksURL:  jwksURL,
		client:   client,
		fetchSem: make(chan struct{}, 1),
	}, nil
}

func checkSandboxJWKSURL(raw string, allowInsecure bool) error {
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("sandbox token verifier: parse JWKS URL: %w", err)
	}
	if u.Scheme == "https" || allowInsecure {
		return nil
	}
	if u.Scheme == "http" && isLoopbackHost(u.Hostname()) {
		return nil
	}
	return fmt.Errorf("sandbox token verifier: JWKS URL must use https, got %q; "+
		"set AllowInsecureJWKS only if the network path to LangSmith is trusted", raw)
}

func isLoopbackHost(host string) bool {
	if host == "localhost" {
		return true
	}
	addr, err := netip.ParseAddr(host)
	return err == nil && addr.IsLoopback()
}

// VerifyUserToken verifies the X-Langsmith-User-Token header of a request made
// through a LangSmith-login service URL.
func (v *SandboxTokenVerifier) VerifyUserToken(ctx context.Context, token string, opts SandboxUserTokenOptions) (*SandboxUser, error) {
	if opts.Audience == "" {
		return nil, errors.New("sandbox token verifier: SandboxUserTokenOptions.Audience is required")
	}
	host, err := sandboxServiceHost(opts.Audience)
	if err != nil {
		return nil, err
	}
	claims, err := v.verify(ctx, token, ExactAudience(host), opts.Issuer)
	if err != nil {
		return nil, err
	}
	if claims.Subject == "" || claims.Subject == sandboxCallbackSubject {
		return nil, verificationError("token is not a user token", nil)
	}
	return &SandboxUser{
		Subject:   claims.Subject,
		Email:     claims.Email,
		Name:      claims.Name,
		ExpiresAt: time.Unix(int64(*claims.ExpiresAt), 0).UTC(),
	}, nil
}

// VerifyCallback verifies a proxy callback request and returns its parsed
// payload. body must be the raw request body, exactly as received, and
// signature the X-LangSmith-Signature-JWT header value.
func (v *SandboxTokenVerifier) VerifyCallback(ctx context.Context, body []byte, signature string, opts SandboxCallbackOptions) (*SandboxCallback, error) {
	claims, err := v.verify(ctx, signature, opts.Audience, opts.Issuer)
	if err != nil {
		return nil, err
	}
	if claims.Subject != sandboxCallbackSubject {
		return nil, verificationError("signature is not a callback signature", nil)
	}
	if claims.BodySHA256 == nil {
		return nil, verificationError("signature has no body hash", nil)
	}
	digest := sha256.Sum256(body)
	if subtle.ConstantTimeCompare([]byte(*claims.BodySHA256), []byte(hex.EncodeToString(digest[:]))) != 1 {
		return nil, verificationError("body does not match signature", nil)
	}
	var cb SandboxCallback
	if err := json.Unmarshal(body, &cb); err != nil {
		return nil, verificationError("malformed callback body", err)
	}
	if cb.Host == "" || cb.Identity.TenantID == "" || cb.Identity.SandboxID == "" {
		return nil, verificationError("malformed callback body", nil)
	}
	return &cb, nil
}

func sandboxServiceHost(audience string) (string, error) {
	if !strings.Contains(audience, "://") {
		return audience, nil
	}
	u, err := url.Parse(audience)
	if err != nil || u.Host == "" {
		return "", fmt.Errorf("sandbox token verifier: audience has no host: %q", audience)
	}
	return u.Host, nil
}

type sandboxTokenHeader struct {
	Alg string `json:"alg"`
	Kid string `json:"kid"`
}

type sandboxTokenClaims struct {
	Issuer     string          `json:"iss"`
	Subject    string          `json:"sub"`
	Audience   sandboxAudience `json:"aud"`
	ExpiresAt  *float64        `json:"exp"`
	IssuedAt   *float64        `json:"iat"`
	NotBefore  *float64        `json:"nbf"`
	Email      string          `json:"email"`
	Name       string          `json:"name"`
	BodySHA256 *string         `json:"body_sha256"`
}

// sandboxAudience accepts the JWT aud claim as either a string or an array.
type sandboxAudience []string

func (a *sandboxAudience) UnmarshalJSON(b []byte) error {
	var one string
	if err := json.Unmarshal(b, &one); err == nil {
		*a = sandboxAudience{one}
		return nil
	}
	var many []string
	if err := json.Unmarshal(b, &many); err != nil {
		return err
	}
	*a = many
	return nil
}

func (v *SandboxTokenVerifier) verify(ctx context.Context, token string, audience func(string) bool, issuer string) (*sandboxTokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, verificationError("malformed token", nil)
	}
	var header sandboxTokenHeader
	if err := decodeSandboxTokenSegment(parts[0], &header); err != nil {
		return nil, verificationError("malformed token header", err)
	}
	var claims sandboxTokenClaims
	if err := decodeSandboxTokenSegment(parts[1], &claims); err != nil {
		return nil, verificationError("malformed token claims", err)
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, verificationError("malformed token signature", err)
	}
	if header.Alg != "EdDSA" {
		return nil, verificationError(fmt.Sprintf("unexpected signing algorithm %q", header.Alg), nil)
	}
	if header.Kid == "" {
		return nil, verificationError("token has no kid", nil)
	}

	key, err := v.key(ctx, header.Kid)
	if err != nil {
		return nil, err
	}
	if !ed25519.Verify(key, []byte(parts[0]+"."+parts[1]), sig) {
		return nil, verificationError("invalid signature", nil)
	}

	now := float64(time.Now().Unix())
	leeway := sandboxTokenLeeway.Seconds()
	switch {
	case claims.ExpiresAt == nil:
		return nil, verificationError("token has no exp", nil)
	case now > *claims.ExpiresAt+leeway:
		return nil, verificationError("token is expired", nil)
	case claims.IssuedAt == nil:
		return nil, verificationError("token has no iat", nil)
	case *claims.IssuedAt > now+leeway:
		return nil, verificationError("token is not yet valid", nil)
	case claims.NotBefore != nil && *claims.NotBefore > now+leeway:
		return nil, verificationError("token is not yet valid", nil)
	case claims.Issuer == "":
		return nil, verificationError("token has no iss", nil)
	case issuer != "" && claims.Issuer != strings.TrimRight(issuer, "/"):
		return nil, verificationError("token has the wrong issuer", nil)
	case len(claims.Audience) == 0:
		return nil, verificationError("token has no aud", nil)
	}
	if audience != nil && !anyAudience(claims.Audience, audience) {
		return nil, verificationError("token has the wrong audience", nil)
	}
	return &claims, nil
}

func anyAudience(auds []string, match func(string) bool) bool {
	for _, aud := range auds {
		if match(aud) {
			return true
		}
	}
	return false
}

func decodeSandboxTokenSegment(segment string, into any) error {
	raw, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, into)
}

func (v *SandboxTokenVerifier) cachedKey(kid string, now time.Time) (ed25519.PublicKey, bool) {
	v.mu.Lock()
	defer v.mu.Unlock()
	key, ok := v.keys[kid]
	age := now.Sub(v.fetchedAt)
	fresh := v.keys != nil && age < sandboxJWKSTTL
	if ok && fresh {
		return key, false
	}
	return key, v.keys == nil || !fresh || age >= sandboxJWKSMinRefresh
}

func (v *SandboxTokenVerifier) key(ctx context.Context, kid string) (ed25519.PublicKey, error) {
	key, refresh := v.cachedKey(kid, time.Now())
	if refresh {
		select {
		case v.fetchSem <- struct{}{}:
		case <-ctx.Done():
			return nil, verificationError("waiting for JWKS", ctx.Err())
		}
		defer func() { <-v.fetchSem }()
		now := time.Now()
		if key, refresh = v.cachedKey(kid, now); refresh {
			keys, err := v.fetchKeys(ctx)
			if err != nil {
				return nil, err
			}
			v.mu.Lock()
			v.keys, v.fetchedAt = keys, now
			v.mu.Unlock()
			key = keys[kid]
		}
	}
	if key == nil {
		return nil, verificationError("unknown signing key "+kid, nil)
	}
	return key, nil
}

func (v *SandboxTokenVerifier) fetchKeys(ctx context.Context) (map[string]ed25519.PublicKey, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.jwksURL, nil)
	if err != nil {
		return nil, verificationError("build JWKS request", err)
	}
	resp, err := v.client.Do(req)
	if err != nil {
		return nil, verificationError("fetch JWKS from "+v.jwksURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, verificationError(fmt.Sprintf("fetch JWKS from %s: HTTP %d", v.jwksURL, resp.StatusCode), nil)
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, sandboxJWKSMaxBytes))
	if err != nil {
		return nil, verificationError("read JWKS", err)
	}
	var set struct {
		Keys []json.RawMessage `json:"keys"`
	}
	if err := json.Unmarshal(raw, &set); err != nil || set.Keys == nil {
		return nil, verificationError("invalid JWKS: expected an object with a keys array", err)
	}
	keys := make(map[string]ed25519.PublicKey, len(set.Keys))
	for _, entry := range set.Keys {
		var jwk struct {
			Kty string `json:"kty"`
			Crv string `json:"crv"`
			Kid string `json:"kid"`
			X   string `json:"x"`
		}
		if json.Unmarshal(entry, &jwk) != nil || jwk.Kty != "OKP" || jwk.Crv != "Ed25519" || jwk.Kid == "" {
			continue
		}
		x, err := base64.RawURLEncoding.DecodeString(jwk.X)
		if err != nil || len(x) != ed25519.PublicKeySize {
			continue
		}
		keys[jwk.Kid] = ed25519.PublicKey(x)
	}
	return keys, nil
}
