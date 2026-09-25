package langsmith

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	verifyTestKID         = "test-kid"
	verifyTestIssuer      = "https://app.example.com"
	verifyTestServiceHost = "0190aaaa-0000-7000-8000-000000000001--8080.svc.example.com"
	verifyTestCallbackURL = "https://integrator.example.com/sandbox-callback"
)

type verifyTestKey struct {
	kid  string
	priv ed25519.PrivateKey
}

func newVerifyTestKey(t *testing.T, kid string) verifyTestKey {
	t.Helper()
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	return verifyTestKey{kid: kid, priv: priv}
}

func (k verifyTestKey) jwk() map[string]any {
	return map[string]any{
		"kty": "OKP", "crv": "Ed25519", "alg": "EdDSA", "use": "sig",
		"kid": k.kid, "key_ops": []string{"verify"}, "ext": true,
		"x": base64.RawURLEncoding.EncodeToString(k.priv.Public().(ed25519.PublicKey)),
	}
}

func (k verifyTestKey) sign(t *testing.T, claims map[string]any) string {
	return signVerifyTestToken(t, k.priv, map[string]any{"alg": "EdDSA", "kid": k.kid, "typ": "JWT"}, claims)
}

func signVerifyTestToken(t *testing.T, priv ed25519.PrivateKey, header, claims map[string]any) string {
	t.Helper()
	h, err := json.Marshal(header)
	require.NoError(t, err)
	c, err := json.Marshal(claims)
	require.NoError(t, err)
	input := base64.RawURLEncoding.EncodeToString(h) + "." + base64.RawURLEncoding.EncodeToString(c)
	return input + "." + base64.RawURLEncoding.EncodeToString(ed25519.Sign(priv, []byte(input)))
}

func userTestClaims(overrides map[string]any) map[string]any {
	now := time.Now().Unix()
	claims := map[string]any{
		"iss": verifyTestIssuer, "sub": "user-123", "aud": []string{verifyTestServiceHost},
		"exp": now + 600, "iat": now, "email": "ada@example.com", "name": "Ada",
	}
	for k, v := range overrides {
		if v == nil {
			delete(claims, k)
		} else {
			claims[k] = v
		}
	}
	return claims
}

func callbackTestBody(t *testing.T, overrides map[string]any) []byte {
	t.Helper()
	payload := map[string]any{
		"host": "api.github.com",
		"port": 443,
		"identity": map[string]any{
			"tenant_id": "c82ac708-2532-4f48-91b5-2756d23af07e", "sandbox_id": "b2b86f7c-53cf-44a4-a972-76a7d1a4b8a6",
			"organization_id": "1bb652ad-e88f-48e8-87d3-60b551a777a3", "ls_user_id": "6f212209-4157-4e51-bf40-187513192ee6",
		},
	}
	for k, v := range overrides {
		payload[k] = v
	}
	body, err := json.Marshal(payload)
	require.NoError(t, err)
	return body
}

func callbackTestClaims(body []byte, overrides map[string]any) map[string]any {
	now := time.Now().Unix()
	digest := sha256.Sum256(body)
	claims := map[string]any{
		"iss": verifyTestIssuer, "sub": "langsmith-sandbox-callback", "aud": []string{verifyTestCallbackURL},
		"iat": now, "nbf": now, "exp": now + 300, "jti": "b904e2e7-347e-404d-9019-a036c2dfc0b0",
		"body_sha256": hex.EncodeToString(digest[:]),
	}
	for k, v := range overrides {
		if v == nil {
			delete(claims, k)
		} else {
			claims[k] = v
		}
	}
	return claims
}

type jwksTestServer struct {
	*httptest.Server
	mu       sync.Mutex
	body     any
	delay    time.Duration
	requests atomic.Int32
}

func newJWKSTestServer(t *testing.T, keys ...verifyTestKey) *jwksTestServer {
	t.Helper()
	s := &jwksTestServer{}
	s.serveKeys(keys...)
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.requests.Add(1)
		s.mu.Lock()
		body, delay := s.body, s.delay
		s.mu.Unlock()
		time.Sleep(delay)
		w.Header().Set("Content-Type", "application/json")
		assert.NoError(t, json.NewEncoder(w).Encode(body))
	}))
	t.Cleanup(s.Close)
	return s
}

func (s *jwksTestServer) serveKeys(keys ...verifyTestKey) {
	jwks := make([]map[string]any, len(keys))
	for i, k := range keys {
		jwks[i] = k.jwk()
	}
	s.serve(map[string]any{"keys": jwks})
}

func (s *jwksTestServer) serve(body any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.body = body
}

func (s *jwksTestServer) verifier(t *testing.T) *SandboxTokenVerifier {
	t.Helper()
	v, err := NewSandboxTokenVerifier(SandboxTokenVerifierOptions{JWKSURL: s.URL + "/.well-known/jwks.json"})
	require.NoError(t, err)
	return v
}

func requireVerificationError(t *testing.T, err error, contains string) {
	t.Helper()
	var verr *SandboxTokenVerificationError
	require.ErrorAs(t, err, &verr)
	assert.Contains(t, err.Error(), contains)
}

func TestSandboxVerifyUserToken(t *testing.T) {
	key := newVerifyTestKey(t, verifyTestKID)
	srv := newJWKSTestServer(t, key)
	v := srv.verifier(t)

	user, err := v.VerifyUserToken(t.Context(), key.sign(t, userTestClaims(nil)), SandboxUserTokenOptions{
		Audience: verifyTestServiceHost, Issuer: verifyTestIssuer,
	})
	require.NoError(t, err)
	assert.Equal(t, "user-123", user.Subject)
	assert.Equal(t, "ada@example.com", user.Email)
	assert.Equal(t, "Ada", user.Name)
	assert.WithinDuration(t, time.Now().Add(10*time.Minute), user.ExpiresAt, time.Minute)

	user, err = v.VerifyUserToken(t.Context(), key.sign(t, userTestClaims(map[string]any{"aud": verifyTestServiceHost})),
		SandboxUserTokenOptions{Audience: "https://" + verifyTestServiceHost + "/path"})
	require.NoError(t, err)
	assert.Equal(t, "user-123", user.Subject)
}

func TestSandboxVerifyUserTokenRejects(t *testing.T) {
	key := newVerifyTestKey(t, verifyTestKID)
	other := newVerifyTestKey(t, verifyTestKID)
	srv := newJWKSTestServer(t, key)
	v := srv.verifier(t)
	now := time.Now().Unix()
	valid := key.sign(t, userTestClaims(nil))
	parts := strings.Split(valid, ".")
	forgedClaims, err := json.Marshal(userTestClaims(map[string]any{"sub": "admin"}))
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		issuer  string
		wantErr string
	}{
		{"wrong audience", key.sign(t, userTestClaims(map[string]any{"aud": []string{"other--8080.svc.example.com"}})), "", "wrong audience"},
		{"expired", key.sign(t, userTestClaims(map[string]any{"exp": now - 120})), "", "expired"},
		{"future iat", key.sign(t, userTestClaims(map[string]any{"iat": now + 600})), "", "not yet valid"},
		{"missing exp", key.sign(t, userTestClaims(map[string]any{"exp": nil})), "", "no exp"},
		{"wrong issuer", valid, "https://evil.example.com", "wrong issuer"},
		{"empty subject", key.sign(t, userTestClaims(map[string]any{"sub": ""})), "", "not a user token"},
		{"callback subject", key.sign(t, userTestClaims(map[string]any{"sub": "langsmith-sandbox-callback"})), "", "not a user token"},
		{"other key", other.sign(t, userTestClaims(nil)), "", "invalid signature"},
		{"tampered payload", parts[0] + "." + base64.RawURLEncoding.EncodeToString(forgedClaims) + "." + parts[2], "", "invalid signature"},
		{"hs256", signVerifyTestToken(t, key.priv, map[string]any{"alg": "HS256", "kid": verifyTestKID}, userTestClaims(nil)), "", "algorithm"},
		{"no kid", signVerifyTestToken(t, key.priv, map[string]any{"alg": "EdDSA"}, userTestClaims(nil)), "", "no kid"},
		{"malformed", "not-a-jwt", "", "malformed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := v.VerifyUserToken(t.Context(), tt.token, SandboxUserTokenOptions{Audience: verifyTestServiceHost, Issuer: tt.issuer})
			requireVerificationError(t, err, tt.wantErr)
		})
	}
}

func TestSandboxVerifyUserTokenRequiresAudience(t *testing.T) {
	key := newVerifyTestKey(t, verifyTestKID)
	v := newJWKSTestServer(t, key).verifier(t)
	_, err := v.VerifyUserToken(t.Context(), key.sign(t, userTestClaims(nil)), SandboxUserTokenOptions{})
	require.ErrorContains(t, err, "Audience is required")
}

func TestSandboxVerifyCallback(t *testing.T) {
	key := newVerifyTestKey(t, verifyTestKID)
	v := newJWKSTestServer(t, key).verifier(t)
	body := callbackTestBody(t, nil)
	sig := key.sign(t, callbackTestClaims(body, nil))

	cb, err := v.VerifyCallback(t.Context(), body, sig, SandboxCallbackOptions{
		Audience: ExactAudience(verifyTestCallbackURL), Issuer: verifyTestIssuer,
	})
	require.NoError(t, err)
	assert.Equal(t, "api.github.com", cb.Host)
	assert.Equal(t, 443, cb.Port)
	assert.Equal(t, "b2b86f7c-53cf-44a4-a972-76a7d1a4b8a6", cb.Identity.SandboxID)
	assert.Nil(t, cb.Request)

	cb, err = v.VerifyCallback(t.Context(), body, sig, SandboxCallbackOptions{
		Audience: func(aud string) bool { return strings.HasPrefix(aud, "https://integrator.example.com/") },
	})
	require.NoError(t, err)
	assert.Equal(t, 443, cb.Port)

	otherAud := key.sign(t, callbackTestClaims(body, map[string]any{"aud": []string{"https://other.example.com/cb"}}))
	_, err = v.VerifyCallback(t.Context(), body, otherAud, SandboxCallbackOptions{})
	require.NoError(t, err)
}

func TestSandboxVerifyCallbackFullRequest(t *testing.T) {
	key := newVerifyTestKey(t, verifyTestKID)
	v := newJWKSTestServer(t, key).verifier(t)
	body := callbackTestBody(t, map[string]any{"request": map[string]any{
		"method": "POST", "url": "https://api.github.com/repos?x=1", "scheme": "https",
		"host": "api.github.com", "path": "/repos", "query": "x=1",
		"headers":     map[string][]string{"Accept": {"application/json"}},
		"body_base64": base64.StdEncoding.EncodeToString([]byte("hello")), "body_truncated": false,
	}})

	cb, err := v.VerifyCallback(t.Context(), body, key.sign(t, callbackTestClaims(body, nil)), SandboxCallbackOptions{})
	require.NoError(t, err)
	require.NotNil(t, cb.Request)
	assert.Equal(t, []byte("hello"), cb.Request.Body)
	assert.Equal(t, map[string][]string{"Accept": {"application/json"}}, cb.Request.Headers)
}

func TestSandboxVerifyCallbackRejects(t *testing.T) {
	key := newVerifyTestKey(t, verifyTestKID)
	v := newJWKSTestServer(t, key).verifier(t)
	body := callbackTestBody(t, nil)
	now := time.Now().Unix()

	tests := []struct {
		name     string
		body     []byte
		sig      string
		audience func(string) bool
		wantErr  string
	}{
		{"tampered body", callbackTestBody(t, map[string]any{"host": "evil.example.com"}), key.sign(t, callbackTestClaims(body, nil)), nil, "body does not match"},
		{"wrong aud", body, key.sign(t, callbackTestClaims(body, nil)), ExactAudience("https://other.example.com/cb"), "wrong audience"},
		{"predicate rejects", body, key.sign(t, callbackTestClaims(body, nil)), func(string) bool { return false }, "wrong audience"},
		{"wrong subject", body, key.sign(t, callbackTestClaims(body, map[string]any{"sub": "user-123"})), nil, "not a callback signature"},
		{"expired", body, key.sign(t, callbackTestClaims(body, map[string]any{"exp": now - 120})), nil, "expired"},
		{"future nbf", body, key.sign(t, callbackTestClaims(body, map[string]any{"nbf": now + 600})), nil, "not yet valid"},
		{"missing body hash", body, key.sign(t, callbackTestClaims(body, map[string]any{"body_sha256": nil})), nil, "no body hash"},
		{"user token", body, key.sign(t, userTestClaims(nil)), nil, "not a callback signature"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := v.VerifyCallback(t.Context(), tt.body, tt.sig, SandboxCallbackOptions{Audience: tt.audience})
			requireVerificationError(t, err, tt.wantErr)
		})
	}
}

func TestSandboxVerifierCachesJWKS(t *testing.T) {
	key := newVerifyTestKey(t, verifyTestKID)
	srv := newJWKSTestServer(t, key)
	v := srv.verifier(t)
	for range 3 {
		_, err := v.VerifyUserToken(t.Context(), key.sign(t, userTestClaims(nil)), SandboxUserTokenOptions{Audience: verifyTestServiceHost})
		require.NoError(t, err)
	}
	assert.EqualValues(t, 1, srv.requests.Load())
}

func TestSandboxVerifierRefetchesForRotatedKey(t *testing.T) {
	key := newVerifyTestKey(t, verifyTestKID)
	srv := newJWKSTestServer(t, key)
	v := srv.verifier(t)
	_, err := v.VerifyUserToken(t.Context(), key.sign(t, userTestClaims(nil)), SandboxUserTokenOptions{Audience: verifyTestServiceHost})
	require.NoError(t, err)

	rotated := newVerifyTestKey(t, "new")
	srv.serveKeys(key, rotated)
	unknown := rotated.sign(t, userTestClaims(nil))
	_, err = v.VerifyUserToken(t.Context(), unknown, SandboxUserTokenOptions{Audience: verifyTestServiceHost})
	requireVerificationError(t, err, "unknown signing key")
	assert.EqualValues(t, 1, srv.requests.Load(), "unknown kid inside the throttle window must not refetch")

	v.mu.Lock()
	v.fetchedAt = v.fetchedAt.Add(-time.Minute)
	v.mu.Unlock()
	user, err := v.VerifyUserToken(t.Context(), unknown, SandboxUserTokenOptions{Audience: verifyTestServiceHost})
	require.NoError(t, err)
	assert.Equal(t, "user-123", user.Subject)
	assert.EqualValues(t, 2, srv.requests.Load())
}

func TestSandboxVerifierConcurrentMissesShareOneFetch(t *testing.T) {
	key := newVerifyTestKey(t, verifyTestKID)
	srv := newJWKSTestServer(t, key)
	srv.mu.Lock()
	srv.delay = 200 * time.Millisecond
	srv.mu.Unlock()
	v := srv.verifier(t)
	token := key.sign(t, userTestClaims(nil))

	var wg sync.WaitGroup
	errs := make([]error, 8)
	for i := range errs {
		wg.Go(func() {
			_, errs[i] = v.VerifyUserToken(t.Context(), token, SandboxUserTokenOptions{Audience: verifyTestServiceHost})
		})
	}
	wg.Wait()
	require.NoError(t, errors.Join(errs...))
	assert.EqualValues(t, 1, srv.requests.Load())
}

func TestSandboxVerifierRejectsMalformedJWKS(t *testing.T) {
	key := newVerifyTestKey(t, verifyTestKID)
	for _, body := range []any{nil, map[string]any{"keys": map[string]any{}}, []any{}, "keys"} {
		srv := newJWKSTestServer(t)
		srv.serve(body)
		_, err := srv.verifier(t).VerifyUserToken(t.Context(), key.sign(t, userTestClaims(nil)), SandboxUserTokenOptions{Audience: verifyTestServiceHost})
		requireVerificationError(t, err, "invalid JWKS")
	}
}

func TestSandboxVerifierSkipsUnusableJWKSEntries(t *testing.T) {
	key := newVerifyTestKey(t, verifyTestKID)
	srv := newJWKSTestServer(t)
	srv.serve(map[string]any{"keys": []any{nil, 1, map[string]any{"kty": "RSA", "kid": "rsa"}, key.jwk()}})
	user, err := srv.verifier(t).VerifyUserToken(t.Context(), key.sign(t, userTestClaims(nil)), SandboxUserTokenOptions{Audience: verifyTestServiceHost})
	require.NoError(t, err)
	assert.Equal(t, "user-123", user.Subject)
}

func TestNewSandboxTokenVerifierJWKSURL(t *testing.T) {
	t.Run("derives from LANGSMITH_ENDPOINT origin", func(t *testing.T) {
		key := newVerifyTestKey(t, verifyTestKID)
		srv := newJWKSTestServer(t, key)
		t.Setenv("LANGSMITH_ENDPOINT", srv.URL+"/api/v1")
		v, err := NewSandboxTokenVerifier(SandboxTokenVerifierOptions{})
		require.NoError(t, err)
		_, err = v.VerifyUserToken(t.Context(), key.sign(t, userTestClaims(nil)), SandboxUserTokenOptions{Audience: verifyTestServiceHost})
		require.NoError(t, err)
	})

	for _, url := range []string{
		"http://localhost:1984/.well-known/jwks.json",
		"http://127.0.0.1:1984/.well-known/jwks.json",
		"http://[::1]:1984/.well-known/jwks.json",
	} {
		t.Run("allows "+url, func(t *testing.T) {
			_, err := NewSandboxTokenVerifier(SandboxTokenVerifierOptions{JWKSURL: url})
			require.NoError(t, err)
		})
	}

	for name, opts := range map[string]SandboxTokenVerifierOptions{
		"http jwks url":         {JWKSURL: "http://langsmith.internal/.well-known/jwks.json"},
		"http api url":          {APIURL: "http://langsmith.internal/api/v1"},
		"loopback-looking name": {JWKSURL: "http://127.example.com/.well-known/jwks.json"},
		"ftp":                   {JWKSURL: "ftp://langsmith.internal/jwks.json"},
	} {
		t.Run("rejects "+name, func(t *testing.T) {
			_, err := NewSandboxTokenVerifier(opts)
			require.ErrorContains(t, err, "must use https")
		})
	}

	t.Run("rejects insecure LANGSMITH_ENDPOINT", func(t *testing.T) {
		t.Setenv("LANGSMITH_ENDPOINT", "http://langsmith.internal/api/v1")
		_, err := NewSandboxTokenVerifier(SandboxTokenVerifierOptions{})
		require.ErrorContains(t, err, "must use https")
	})

	t.Run("AllowInsecureJWKS opts out", func(t *testing.T) {
		_, err := NewSandboxTokenVerifier(SandboxTokenVerifierOptions{
			JWKSURL: "http://langsmith.internal/.well-known/jwks.json", AllowInsecureJWKS: true,
		})
		require.NoError(t, err)
	})
}
