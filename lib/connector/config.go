package connector

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"
)

const (
	// DefaultEndpoint is the LangSmith API URL used when none is configured.
	DefaultEndpoint = "https://api.smith.langchain.com"

	// carrierPath is where the edge serves connector carriers. It is derived
	// from the API endpoint rather than configured: which transport and path
	// the tunnel uses is an internal protocol detail.
	carrierPath = "/v1/connections/connect"
)

// Config contains the validated settings for a connector process.
type Config struct {
	// Endpoint is the carrier URL the connector dials, derived from the
	// configured LangSmith API URL.
	Endpoint     string
	ConnectionID uuid.UUID
	Targets      TargetMap

	apiKey     string
	apiKeyFile string
}

// Settings are a connector's settings before validation.
type Settings struct {
	// Endpoint is the LangSmith API URL, DefaultEndpoint when empty. The
	// carrier URL is derived from it.
	Endpoint string
	// ConnectionID is the connection this connector serves.
	ConnectionID string
	// APIKey is the API key itself. APIKeyFile is instead a file holding it,
	// read again on every carrier attempt so a rotated key needs no restart.
	// Set exactly one.
	APIKey     string
	APIKeyFile string
	// Targets maps target names to host:port for raw TCP, or to an
	// http(s)://host:port URL for an HTTP service, for a tunnel.
	Targets map[string]string
	// Relay serves a relay instead of a tunnel: LangSmith names each
	// destination, and traffic to it leaves from the connector's network.
	// RelayAllow and RelayDeny take the rules ParseProxyPolicy accepts.
	Relay      bool
	RelayAllow []string
	RelayDeny  []string
}

var (
	errAPIKeyRequired       = errors.New("an API key is required")
	errAPIKeySourceConflict = errors.New("set an API key or an API key file, not both")
	errRelayRulesNeedRelay  = errors.New("relay allow and deny rules require relay mode")
)

// NewConfig validates settings.
func NewConfig(settings Settings) (Config, error) {
	endpoint, err := carrierEndpoint(firstNonEmpty(settings.Endpoint, DefaultEndpoint))
	if err != nil {
		return Config{}, err
	}
	connectionID, err := uuid.Parse(strings.TrimSpace(settings.ConnectionID))
	if err != nil || connectionID == uuid.Nil {
		return Config{}, fmt.Errorf("connection ID must be a non-zero UUID")
	}
	proxy, err := relayPolicy(settings)
	if err != nil {
		return Config{}, err
	}
	targets, err := newTargetMap(settings.Targets, proxy)
	if err != nil {
		return Config{}, err
	}

	config := Config{
		Endpoint:     endpoint,
		ConnectionID: connectionID,
		Targets:      targets,
		apiKey:       strings.TrimSpace(settings.APIKey),
		apiKeyFile:   settings.APIKeyFile,
	}
	switch {
	case config.apiKey != "" && config.apiKeyFile != "":
		return Config{}, errAPIKeySourceConflict
	case config.apiKey == "" && config.apiKeyFile == "":
		return Config{}, errAPIKeyRequired
	}
	if _, err := config.ReadAPIKey(); err != nil {
		return Config{}, err
	}
	return config, nil
}

// relayPolicy returns the relay policy, or nil when relay mode is off.
func relayPolicy(settings Settings) (*ProxyPolicy, error) {
	if !settings.Relay {
		if len(settings.RelayAllow) > 0 || len(settings.RelayDeny) > 0 {
			return nil, errRelayRulesNeedRelay
		}
		return nil, nil
	}
	policy, err := ParseProxyPolicy(settings.RelayAllow, settings.RelayDeny)
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

// ReadAPIKey returns the resolved credential. When a credential file is
// configured, it is re-read on every call so a rotated secret is picked up on
// the next carrier attempt.
func (c Config) ReadAPIKey() (string, error) {
	if c.apiKey != "" {
		return c.apiKey, nil
	}
	if c.apiKeyFile == "" {
		return "", fmt.Errorf("API key is required")
	}
	contents, err := os.ReadFile(c.apiKeyFile)
	if err != nil {
		return "", fmt.Errorf("read API key file: %w", err)
	}
	apiKey := strings.TrimSpace(string(contents))
	if apiKey == "" {
		return "", fmt.Errorf("API key file is empty")
	}
	return apiKey, nil
}

// carrierEndpoint turns a LangSmith API URL into the carrier URL to dial,
// preserving any base path so a self-hosted install mounted under a prefix
// keeps working. Callers configure the same API URL they use everywhere else;
// the WebSocket scheme and carrier path are supplied here.
func carrierEndpoint(rawEndpoint string) (string, error) {
	if strings.TrimSpace(rawEndpoint) != rawEndpoint {
		return "", fmt.Errorf("endpoint must not have surrounding whitespace")
	}
	parsed, err := url.Parse(rawEndpoint)
	if err != nil {
		return "", fmt.Errorf("endpoint must be a valid LangSmith API URL")
	}

	var scheme string
	switch parsed.Scheme {
	case "https":
		scheme = "wss"
	case "http":
		if !isLoopbackHostname(parsed.Hostname()) {
			return "", fmt.Errorf("endpoint must use https unless the host is loopback")
		}
		scheme = "ws"
	case "ws", "wss":
		return "", fmt.Errorf(
			"endpoint must be the LangSmith API URL (for example %s), not a WebSocket URL; "+
				"the tunnel URL is derived automatically", DefaultEndpoint)
	default:
		return "", fmt.Errorf("endpoint must be an http(s) LangSmith API URL")
	}
	if parsed.Host == "" || parsed.Hostname() == "" || parsed.Opaque != "" {
		return "", fmt.Errorf("endpoint must include a host")
	}
	if parsed.User != nil {
		return "", fmt.Errorf("endpoint must not contain credentials")
	}
	if parsed.RawQuery != "" || parsed.ForceQuery {
		return "", fmt.Errorf("endpoint must not contain a query")
	}
	if parsed.Fragment != "" || parsed.RawFragment != "" {
		return "", fmt.Errorf("endpoint must not contain a fragment")
	}

	carrier := url.URL{
		Scheme: scheme,
		Host:   parsed.Host,
		Path:   strings.TrimSuffix(parsed.Path, "/") + carrierPath,
	}
	return carrier.String(), nil
}

// validateCarrierURL checks a value already in carrier form — what
// carrierEndpoint produces and what Client dials.
func validateCarrierURL(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("connector endpoint is not a valid URL")
	}
	if (parsed.Scheme != "wss" && parsed.Scheme != "ws") || parsed.Host == "" || parsed.Hostname() == "" {
		return fmt.Errorf("connector endpoint must be a WebSocket URL with a host")
	}
	if parsed.Scheme == "ws" && !isLoopbackHostname(parsed.Hostname()) {
		return fmt.Errorf("connector endpoint must use wss unless the host is loopback")
	}
	if parsed.User != nil {
		return fmt.Errorf("connector endpoint must not contain credentials")
	}
	return nil
}

func isLoopbackHostname(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
