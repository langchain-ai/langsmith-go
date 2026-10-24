package connector

import (
	"context"
	"errors"
	"fmt"
	"net"
	"slices"
	"strconv"
	"strings"
)

// ProxyTargetName is the target a connector in relay mode registers. Each flow
// to it names its own destination.
const ProxyTargetName = "proxy"

// defaultProxyPorts are the ports the default proxy rule allows on public
// addresses. Every other port must be allowed explicitly, so a proxy cannot
// relay mail or other abuse-prone protocols unless its owner opts in.
var defaultProxyPorts = []int{80, 443}

var errDestinationDenied = errors.New("destination denied by relay policy")

// ProxyPolicy decides which destinations a connector in relay mode dials. By
// default it allows public addresses on ports 80 and 443; Allow rules widen
// that and Deny rules narrow it. Cloud metadata is always refused.
type ProxyPolicy struct {
	allow []proxyRule
	deny  []proxyRule
	spec  string
}

// proxyRule matches a destination by CIDR or by hostname (exact or "*.suffix"),
// optionally limited to one port.
type proxyRule struct {
	network    *net.IPNet
	host       string
	hostSuffix string
	port       int
}

// ParseProxyPolicy validates allow and deny rules. Each rule is a CIDR
// ("10.0.0.0/8"), an IP address, an exact hostname ("db.corp.internal"), or a
// wildcard ("*.reddit.com"), optionally followed by ":PORT". An IPv6 address or
// range takes a port only in brackets ("[2001:db8::1]:443").
func ParseProxyPolicy(allow, deny []string) (ProxyPolicy, error) {
	allowRules, err := parseProxyRules(allow)
	if err != nil {
		return ProxyPolicy{}, fmt.Errorf("relay allow: %w", err)
	}
	denyRules, err := parseProxyRules(deny)
	if err != nil {
		return ProxyPolicy{}, fmt.Errorf("relay deny: %w", err)
	}
	return ProxyPolicy{
		allow: allowRules,
		deny:  denyRules,
		spec:  "allow=" + canonicalRules(allowRules) + ";deny=" + canonicalRules(denyRules),
	}, nil
}

// canonicalRules writes rules in a form that ignores case, order, and
// duplicates, so equivalent policies describe themselves identically.
func canonicalRules(rules []proxyRule) string {
	texts := make([]string, 0, len(rules))
	for _, rule := range rules {
		texts = append(texts, rule.String())
	}
	slices.Sort(texts)
	return strings.Join(slices.Compact(texts), ",")
}

// String describes the policy in a stable form; it feeds the target map hash
// so replicas of one connector must agree on their relay policy.
func (p ProxyPolicy) String() string {
	return p.spec
}

func parseProxyRules(raw []string) ([]proxyRule, error) {
	rules := make([]proxyRule, 0, len(raw))
	for _, entry := range raw {
		rule, err := parseProxyRule(strings.TrimSpace(entry))
		if err != nil {
			return nil, fmt.Errorf("%q: %w", entry, err)
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func parseProxyRule(entry string) (proxyRule, error) {
	if entry == "" {
		return proxyRule{}, errors.New("rule must not be empty")
	}
	pattern, port, err := splitRulePort(entry)
	if err != nil {
		return proxyRule{}, err
	}
	rule := proxyRule{port: port}
	if strings.Contains(pattern, "/") {
		_, network, err := net.ParseCIDR(pattern)
		if err != nil {
			return proxyRule{}, fmt.Errorf("invalid CIDR: %w", err)
		}
		rule.network = network
		return rule, nil
	}
	if ip := net.ParseIP(pattern); ip != nil {
		bits := 32
		if ip.To4() == nil {
			bits = 128
		}
		rule.network = &net.IPNet{IP: ip, Mask: net.CIDRMask(bits, bits)}
		return rule, nil
	}
	host := strings.ToLower(pattern)
	if suffix, ok := strings.CutPrefix(host, "*."); ok {
		if suffix == "" || strings.Contains(suffix, "*") {
			return proxyRule{}, errors.New("wildcard must be *.domain")
		}
		rule.hostSuffix = "." + suffix
		return rule, nil
	}
	if strings.Contains(host, "*") {
		return proxyRule{}, errors.New("wildcard must be *.domain")
	}
	if strings.ContainsAny(host, ":[]") {
		return proxyRule{}, errors.New("host must be a DNS name; put an IPv6 address with a port in brackets")
	}
	rule.host = host
	return rule, nil
}

// splitRulePort separates a rule's optional ":PORT". Unbracketed, every colon
// in an IPv6 address belongs to the address, so "2001:db8::1:443" is one
// address on every port and "[2001:db8::1]:443" is that address on port 443.
func splitRulePort(entry string) (string, int, error) {
	if strings.HasPrefix(entry, "[") {
		pattern, rawPort, err := net.SplitHostPort(entry)
		if err != nil {
			return "", 0, errors.New("a bracketed rule must be [address]:port")
		}
		port, err := parseRulePort(rawPort)
		return pattern, port, err
	}
	index := strings.LastIndex(entry, ":")
	if index <= 0 || strings.Contains(entry[:index], ":") {
		return entry, 0, nil
	}
	port, err := parseRulePort(entry[index+1:])
	return entry[:index], port, err
}

func parseRulePort(raw string) (int, error) {
	port, err := strconv.Atoi(raw)
	if err != nil || port < 1 || port > 65535 {
		return 0, errors.New("port must be 1-65535")
	}
	return port, nil
}

// String is the rule's canonical text.
func (r proxyRule) String() string {
	var pattern string
	switch {
	case r.network != nil:
		pattern = r.network.String()
	case r.hostSuffix != "":
		pattern = "*" + r.hostSuffix
	default:
		pattern = r.host
	}
	if r.port == 0 {
		return pattern
	}
	return net.JoinHostPort(pattern, strconv.Itoa(r.port))
}

func (r proxyRule) matchesPort(port int) bool {
	return r.port == 0 || r.port == port
}

func (r proxyRule) matchesHost(host string, port int) bool {
	if !r.matchesPort(port) {
		return false
	}
	return (r.host != "" && r.host == host) ||
		(r.hostSuffix != "" && strings.HasSuffix(host, r.hostSuffix))
}

func (r proxyRule) matchesIP(ip net.IP, port int) bool {
	return r.network != nil && r.matchesPort(port) && r.network.Contains(ip)
}

// Resolver looks up a proxy destination's addresses. *net.Resolver
// satisfies it.
type Resolver interface {
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
}

// dial resolves destination, checks every address against the policy, and
// connects to the first allowed one. It returns errDestinationDenied when the
// policy refuses the destination or all of its addresses.
func (p ProxyPolicy) dial(
	ctx context.Context,
	destination string,
	lookup Resolver,
	dialer TargetDialer,
) (net.Conn, error) {
	host, rawPort, err := net.SplitHostPort(destination)
	if err != nil {
		return nil, errDestinationDenied
	}
	port, err := strconv.Atoi(rawPort)
	if err != nil {
		return nil, errDestinationDenied
	}
	host = strings.ToLower(host)
	if isMetadataHost(host) {
		return nil, errDestinationDenied
	}
	if slices.ContainsFunc(p.deny, func(rule proxyRule) bool { return rule.matchesHost(host, port) }) {
		return nil, errDestinationDenied
	}
	hostAllowed := slices.ContainsFunc(p.allow, func(rule proxyRule) bool { return rule.matchesHost(host, port) })

	addresses, err := lookupAddresses(ctx, lookup, host)
	if err != nil {
		return nil, err
	}
	lastErr := errDestinationDenied
	for _, ip := range addresses {
		if !p.allowsAddress(ip, port, hostAllowed) {
			continue
		}
		conn, dialErr := dialer.DialContext(ctx, "tcp", net.JoinHostPort(ip.String(), rawPort))
		if dialErr == nil {
			return conn, nil
		}
		lastErr = dialErr
	}
	return nil, lastErr
}

func (p ProxyPolicy) allowsAddress(ip net.IP, port int, hostAllowed bool) bool {
	if isMetadataAddress(ip) {
		return false
	}
	if slices.ContainsFunc(p.deny, func(rule proxyRule) bool { return rule.matchesIP(ip, port) }) {
		return false
	}
	if hostAllowed {
		return true
	}
	if slices.ContainsFunc(p.allow, func(rule proxyRule) bool { return rule.matchesIP(ip, port) }) {
		return true
	}
	return slices.Contains(defaultProxyPorts, port) && isPublicAddress(ip)
}

func lookupAddresses(ctx context.Context, lookup Resolver, host string) ([]net.IP, error) {
	if ip := net.ParseIP(host); ip != nil {
		return []net.IP{ip}, nil
	}
	answers, err := lookup.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}
	addresses := make([]net.IP, 0, len(answers))
	for _, answer := range answers {
		addresses = append(addresses, answer.IP)
	}
	return addresses, nil
}
