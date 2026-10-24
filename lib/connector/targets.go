package connector

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"net"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/langchain-ai/langsmith-go/lib/connector/transport"
)

const maxHostnameLength = 253

// RegistrationMetadata is the target metadata sent to the edge. TCP dial
// addresses stay connector-local; HTTP descriptors include the authority the
// edge needs for Host rewriting while proxying.
type RegistrationMetadata struct {
	Targets       []transport.TargetDescriptor `json:"targets"`
	TargetMapHash string                       `json:"target_map_hash"`
}

// Target is the connector-local target selected by a logical name. Address is
// used only for connector-side dialing; Descriptor contains the HTTP routing
// metadata registered with the edge. A proxy target has no fixed Address; Proxy
// decides which destinations its flows may reach.
type Target struct {
	Address    string
	Descriptor transport.TargetDescriptor
	Proxy      *ProxyPolicy
}

// TargetMap is an immutable map from logical target name to local TCP address.
type TargetMap struct {
	targets map[string]Target
	names   []string
	hash    string
}

// NewTargetMap validates, canonicalizes, and copies the configured targets.
func NewTargetMap(raw map[string]string) (TargetMap, error) {
	return newTargetMap(raw, nil)
}

// newTargetMap builds the target map. A connector serves either named targets
// for a tunnel or the single proxy target for a proxy, never both, matching the
// connection type the edge enforces.
func newTargetMap(raw map[string]string, proxy *ProxyPolicy) (TargetMap, error) {
	if proxy != nil && len(raw) > 0 {
		return TargetMap{}, fmt.Errorf("a connector serves either targets or relay mode, not both")
	}
	count := len(raw)
	if proxy != nil {
		count++
	}
	if count == 0 {
		return TargetMap{}, fmt.Errorf("targets must not be empty")
	}
	if count > transport.MaxRegisteredTargets {
		return TargetMap{}, fmt.Errorf("targets must contain at most %d entries", transport.MaxRegisteredTargets)
	}

	targets := make(map[string]Target, count)
	names := make([]string, 0, count)
	if proxy != nil {
		targets[ProxyTargetName] = Target{
			Descriptor: transport.TargetDescriptor{Name: ProxyTargetName, Kind: transport.TargetKindProxy},
			Proxy:      proxy,
		}
		names = append(names, ProxyTargetName)
	}
	for name, rawAddress := range raw {
		if name == ProxyTargetName {
			return TargetMap{}, fmt.Errorf("target name %q is reserved for relay mode", ProxyTargetName)
		}
		if err := validateTargetName(name); err != nil {
			return TargetMap{}, err
		}
		target, err := parseTarget(name, rawAddress)
		if err != nil {
			return TargetMap{}, fmt.Errorf("target %q: %w", name, err)
		}
		targets[name] = target
		names = append(names, name)
	}
	slices.Sort(names)

	return TargetMap{
		targets: targets,
		names:   names,
		hash:    hashTargetMap(names, targets),
	}, nil
}

// Lookup returns the exact canonical address configured for name.
func (t TargetMap) Lookup(name string) (string, bool) {
	target, ok := t.targets[name]
	return target.Address, ok
}

// Resolve returns the connector-local target selected by name.
func (t TargetMap) Resolve(name string) (Target, bool) {
	target, ok := t.targets[name]
	return target, ok
}

// Names returns the target names in canonical sorted order.
func (t TargetMap) Names() []string {
	return slices.Clone(t.names)
}

// Kind is "relay" for a connector serving the proxy target and "tunnel"
// otherwise.
func (t TargetMap) Kind() string {
	if _, ok := t.targets[ProxyTargetName]; ok && len(t.targets) == 1 {
		return "relay"
	}
	return "tunnel"
}

// Hash returns the lowercase hexadecimal SHA-256 hash of the canonical map.
func (t TargetMap) Hash() string {
	return t.hash
}

// RegistrationMetadata returns the bounded description advertised to the edge.
func (t TargetMap) RegistrationMetadata() RegistrationMetadata {
	descriptors := make([]transport.TargetDescriptor, 0, len(t.names))
	for _, name := range t.names {
		descriptors = append(descriptors, t.targets[name].Descriptor)
	}
	return RegistrationMetadata{
		Targets:       descriptors,
		TargetMapHash: t.hash,
	}
}

func parseTarget(name, raw string) (Target, error) {
	if raw == "" || strings.TrimSpace(raw) != raw {
		return Target{}, fmt.Errorf("address must be an exact host:port or http(s) URL")
	}
	if !strings.Contains(raw, "://") {
		address, err := canonicalTargetAddress(raw)
		if err != nil {
			return Target{}, err
		}
		return Target{
			Address: address,
			Descriptor: transport.TargetDescriptor{
				Name: name,
				Kind: transport.TargetKindTCP,
			},
		}, nil
	}

	parsed, err := url.Parse(raw)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") ||
		parsed.Opaque != "" || parsed.User != nil || parsed.Host == "" ||
		parsed.Path != "" || parsed.RawPath != "" || parsed.RawQuery != "" ||
		parsed.ForceQuery || parsed.Fragment != "" || parsed.RawFragment != "" {
		return Target{}, fmt.Errorf("HTTP target must be an exact http(s)://host:port URL without credentials, path, query, or fragment")
	}
	address, err := canonicalTargetAddress(parsed.Host)
	if err != nil {
		return Target{}, fmt.Errorf("HTTP target: %w", err)
	}
	return Target{
		Address: address,
		Descriptor: transport.TargetDescriptor{
			Name:            name,
			Kind:            transport.TargetKindHTTP,
			OriginScheme:    parsed.Scheme,
			OriginAuthority: address,
		},
	}, nil
}

func validateTargetName(name string) error {
	if err := transport.ValidateTargetName(name); err != nil {
		return fmt.Errorf("target name %q: %w", name, err)
	}
	return nil
}

func canonicalTargetAddress(raw string) (string, error) {
	if raw == "" || strings.TrimSpace(raw) != raw {
		return "", fmt.Errorf("address must be an exact host:port")
	}

	host, rawPort, err := net.SplitHostPort(raw)
	if err != nil || host == "" {
		return "", fmt.Errorf("address must be an exact host:port")
	}
	port, err := canonicalPort(rawPort)
	if err != nil {
		return "", err
	}

	if ip := net.ParseIP(host); ip != nil {
		if ip.IsUnspecified() {
			return "", fmt.Errorf("address must not use a wildcard host")
		}
		return net.JoinHostPort(ip.String(), port), nil
	}
	if strings.HasPrefix(raw, "[") {
		return "", fmt.Errorf("bracketed host must be an IPv6 address")
	}

	host, err = canonicalHostname(host)
	if err != nil {
		return "", err
	}
	return net.JoinHostPort(host, port), nil
}

func canonicalPort(raw string) (string, error) {
	if raw == "" || len(raw) > 5 {
		return "", fmt.Errorf("port must be an integer from 1 to 65535")
	}
	for index := range len(raw) {
		if raw[index] < '0' || raw[index] > '9' {
			return "", fmt.Errorf("port must be an integer from 1 to 65535")
		}
	}
	port, err := strconv.Atoi(raw)
	if err != nil || port < 1 || port > 65535 {
		return "", fmt.Errorf("port must be an integer from 1 to 65535")
	}
	return strconv.Itoa(port), nil
}

func canonicalHostname(raw string) (string, error) {
	hostname := strings.ToLower(strings.TrimSuffix(raw, "."))
	if hostname == "" || len(hostname) > maxHostnameLength || strings.Contains(hostname, "*") {
		return "", fmt.Errorf("host must be a valid DNS name or IP address")
	}

	labels := strings.Split(hostname, ".")
	if looksLikeInvalidIPv4(labels) {
		return "", fmt.Errorf("host must be a valid DNS name or IP address")
	}
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 || !isAlphaNumeric(label[0]) || !isAlphaNumeric(label[len(label)-1]) {
			return "", fmt.Errorf("host must be a valid DNS name or IP address")
		}
		for index := range len(label) {
			if !isAlphaNumeric(label[index]) && label[index] != '-' {
				return "", fmt.Errorf("host must be a valid DNS name or IP address")
			}
		}
	}
	return hostname, nil
}

func looksLikeInvalidIPv4(labels []string) bool {
	if len(labels) != net.IPv4len {
		return false
	}
	for _, label := range labels {
		if label == "" {
			return false
		}
		for index := range len(label) {
			if label[index] < '0' || label[index] > '9' {
				return false
			}
		}
	}
	return true
}

func hashTargetMap(names []string, targets map[string]Target) string {
	hasher := sha256.New()
	_, _ = hasher.Write([]byte("langsmith-target-map-v1\x00"))
	for _, name := range names {
		writeHashString(hasher, name)
		target := targets[name]
		canonical := target.Address
		switch target.Descriptor.Kind {
		case transport.TargetKindHTTP:
			canonical = target.Descriptor.OriginScheme + "://" + target.Address
		case transport.TargetKindProxy:
			canonical = "proxy:" + target.Proxy.String()
		}
		writeHashString(hasher, canonical)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func writeHashString(hasher hash.Hash, value string) {
	var size [4]byte
	binary.BigEndian.PutUint32(size[:], uint32(len(value)))
	_, _ = hasher.Write(size[:])
	_, _ = hasher.Write([]byte(value))
}

func isAlphaNumeric(char byte) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9'
}
