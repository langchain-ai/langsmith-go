package connector

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"slices"
	"strings"
	"syscall"
	"time"
)

// A connector exists to reach the network it runs in, so it does not refuse
// private addresses the way a LangSmith service refuses them for outbound
// requests. It needs only two checks: cloud metadata is refused on every flow,
// and a relay's default rule admits only public addresses.

var errMetadataDestination = errors.New("cloud metadata endpoints are not reachable through a connector")

var (
	// nonPublicIPv4 are the private, loopback, link-local, shared, reserved,
	// documentation, benchmarking, multicast, and broadcast IPv4 ranges.
	nonPublicIPv4 = mustParsePrefixes(
		"0.0.0.0/8",
		"10.0.0.0/8",
		"100.64.0.0/10",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"172.16.0.0/12",
		"192.0.0.0/24",
		"192.0.2.0/24",
		"192.168.0.0/16",
		"198.18.0.0/15",
		"198.51.100.0/24",
		"203.0.113.0/24",
		"224.0.0.0/4",
		"240.0.0.0/4",
		"255.255.255.255/32",
	)
	// nonPublicIPv6 are the loopback, unique local, link-local, and multicast
	// IPv6 ranges.
	nonPublicIPv6 = mustParsePrefixes("::1/128", "fc00::/7", "fe80::/10", "ff00::/8")
	// nat64Prefixes embed an IPv4 address in their last four bytes (RFC 6052,
	// RFC 8215).
	nat64Prefixes = mustParsePrefixes("64:ff9b::/96", "64:ff9b:1::/48")

	metadataAddresses = []netip.Addr{
		netip.MustParseAddr("169.254.169.254"),
		netip.MustParseAddr("169.254.170.2"),
		netip.MustParseAddr("100.100.100.200"),
		netip.MustParseAddr("fd00:ec2::254"),
	}
	metadataHosts = []string{"metadata.google.internal", "metadata.amazonaws.com", "metadata", "instance-data"}
)

func mustParsePrefixes(cidrs ...string) []netip.Prefix {
	prefixes := make([]netip.Prefix, len(cidrs))
	for i, cidr := range cidrs {
		prefixes[i] = netip.MustParsePrefix(cidr)
	}
	return prefixes
}

// addressForms returns ip and any IPv4 address it carries as an IPv4-mapped
// or NAT64 address, so a check cannot be bypassed by rewriting the address.
func addressForms(ip net.IP) []netip.Addr {
	addr, ok := netip.AddrFromSlice(ip)
	if !ok {
		return nil
	}
	forms := []netip.Addr{addr}
	if addr.Is4In6() {
		forms = append(forms, addr.Unmap())
	}
	for _, prefix := range nat64Prefixes {
		if prefix.Contains(addr) {
			embedded := addr.As16()
			forms = append(forms, netip.AddrFrom4([4]byte(embedded[12:])))
			break
		}
	}
	return forms
}

// isMetadataHost reports whether host names a cloud metadata endpoint.
func isMetadataHost(host string) bool {
	host = strings.TrimSuffix(host, ".")
	return slices.ContainsFunc(metadataHosts, func(name string) bool { return strings.EqualFold(name, host) })
}

// isMetadataAddress reports whether ip, in any form, is a cloud metadata
// endpoint. An unparseable address counts as one, so it is refused.
func isMetadataAddress(ip net.IP) bool {
	forms := addressForms(ip)
	if len(forms) == 0 {
		return true
	}
	return slices.ContainsFunc(forms, func(form netip.Addr) bool {
		return slices.Contains(metadataAddresses, form.Unmap())
	})
}

// isPublicAddress reports whether ip, in every form, is a public address and
// not a cloud metadata endpoint.
func isPublicAddress(ip net.IP) bool {
	forms := addressForms(ip)
	if len(forms) == 0 || isMetadataAddress(ip) {
		return false
	}
	return !slices.ContainsFunc(forms, func(form netip.Addr) bool {
		ranges := nonPublicIPv6
		if form.Is4() {
			ranges = nonPublicIPv4
		}
		return slices.ContainsFunc(ranges, func(prefix netip.Prefix) bool { return prefix.Contains(form) })
	})
}

// metadataGuardDialer dials any destination except cloud metadata. It checks
// the hostname, then every address the dialer resolves immediately before it
// connects, so a name that later resolves to a metadata address is refused.
type metadataGuardDialer struct {
	dialer net.Dialer
}

func newMetadataGuardDialer() *metadataGuardDialer {
	return &metadataGuardDialer{dialer: net.Dialer{
		Timeout: 10 * time.Second,
		ControlContext: func(_ context.Context, _, address string, _ syscall.RawConn) error {
			host, _, err := net.SplitHostPort(address)
			if err != nil {
				return fmt.Errorf("invalid resolved address %q: %w", address, err)
			}
			if isMetadataAddress(net.ParseIP(host)) {
				return errMetadataDestination
			}
			return nil
		},
	}}
}

func (d *metadataGuardDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return nil, fmt.Errorf("invalid address %q: %w", address, err)
	}
	if isMetadataHost(host) {
		return nil, errMetadataDestination
	}
	return d.dialer.DialContext(ctx, network, address)
}
