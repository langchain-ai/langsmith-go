package connector

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/langchain-ai/langsmith-go/lib/connector/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type staticResolver map[string][]string

func (r staticResolver) LookupIPAddr(_ context.Context, host string) ([]net.IPAddr, error) {
	answers, ok := r[host]
	if !ok {
		return nil, &net.DNSError{Err: "no such host", Name: host, IsNotFound: true}
	}
	addresses := make([]net.IPAddr, 0, len(answers))
	for _, answer := range answers {
		addresses = append(addresses, net.IPAddr{IP: net.ParseIP(answer)})
	}
	return addresses, nil
}

// recordingDialer accepts every dial and records the addresses it was asked
// to reach.
type recordingDialer struct {
	mu        sync.Mutex
	addresses []string
}

func (d *recordingDialer) DialContext(_ context.Context, _, address string) (net.Conn, error) {
	d.mu.Lock()
	d.addresses = append(d.addresses, address)
	d.mu.Unlock()
	local, remote := net.Pipe()
	_ = remote.Close()
	return local, nil
}

func TestProxyPolicyDial(t *testing.T) {
	t.Parallel()

	resolver := staticResolver{
		"www.reddit.com":      {"151.101.1.140"},
		"db.corp.internal":    {"10.0.0.5"},
		"mixed.example.com":   {"10.0.0.9", "93.184.216.34"},
		"rebind.example.com":  {"169.254.169.254"},
		"loopback.example.io": {"127.0.0.1"},
	}
	tests := []struct {
		name        string
		allow, deny []string
		destination string
		wantDialed  string
		wantDenied  bool
	}{
		{name: "public host on 443 by default", destination: "www.reddit.com:443", wantDialed: "151.101.1.140:443"},
		{name: "public IP literal on 80 by default", destination: "93.184.216.34:80", wantDialed: "93.184.216.34:80"},
		{name: "public host on other port needs allow", destination: "www.reddit.com:25", wantDenied: true},
		{name: "private host denied by default", destination: "db.corp.internal:443", wantDenied: true},
		{name: "loopback denied by default", destination: "loopback.example.io:443", wantDenied: true},
		{
			name:        "private host allowed by exact host rule",
			allow:       []string{"db.corp.internal:5432"},
			destination: "db.corp.internal:5432",
			wantDialed:  "10.0.0.5:5432",
		},
		{
			name:        "host rule port must match",
			allow:       []string{"db.corp.internal:5432"},
			destination: "db.corp.internal:5433",
			wantDenied:  true,
		},
		{
			name:        "private address allowed by CIDR rule",
			allow:       []string{"10.0.0.0/8"},
			destination: "db.corp.internal:6379",
			wantDialed:  "10.0.0.5:6379",
		},
		{
			name:        "wildcard deny wins over the default",
			deny:        []string{"*.reddit.com"},
			destination: "www.reddit.com:443",
			wantDenied:  true,
		},
		{
			name:        "CIDR deny wins over a host allow",
			allow:       []string{"db.corp.internal"},
			deny:        []string{"10.0.0.0/24"},
			destination: "db.corp.internal:443",
			wantDenied:  true,
		},
		{
			name:        "skips denied addresses and dials an allowed one",
			destination: "mixed.example.com:443",
			wantDialed:  "93.184.216.34:443",
		},
		{
			name:        "metadata address denied even when the host is allowed",
			allow:       []string{"rebind.example.com"},
			destination: "rebind.example.com:80",
			wantDenied:  true,
		},
		{
			name:        "metadata hostname denied even when allowed",
			allow:       []string{"metadata.google.internal"},
			destination: "metadata.google.internal:80",
			wantDenied:  true,
		},
		{
			name:        "bracketed IPv6 rule allows only its port",
			allow:       []string{"[2001:db8::1]:5432"},
			destination: "[2001:db8::1]:5432",
			wantDialed:  "[2001:db8::1]:5432",
		},
		{
			name:        "bracketed IPv6 rule does not widen to other ports",
			allow:       []string{"[2001:db8::1]:5432"},
			destination: "[2001:db8::1]:6379",
			wantDenied:  true,
		},
		{
			name:        "metadata IP denied even when its CIDR is allowed",
			allow:       []string{"169.254.0.0/16"},
			destination: "169.254.169.254:80",
			wantDenied:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			policy, err := ParseProxyPolicy(test.allow, test.deny)
			require.NoError(t, err)
			dialer := &recordingDialer{}

			conn, err := policy.dial(context.Background(), test.destination, resolver, dialer)

			if test.wantDenied {
				require.ErrorIs(t, err, errDestinationDenied)
				assert.Empty(t, dialer.addresses)
				return
			}
			require.NoError(t, err)
			_ = conn.Close()
			assert.Equal(t, []string{test.wantDialed}, dialer.addresses)
		})
	}
}

func TestParseProxyPolicyRejectsMalformedRules(t *testing.T) {
	t.Parallel()

	for _, rule := range []string{
		"", "*", "foo.*.com", "10.0.0.0/33", "host:0", "host:70000", "host:http",
		"[2001:db8::1]", "[2001:db8::1]:0", "db[1].corp",
	} {
		_, err := ParseProxyPolicy([]string{rule}, nil)
		assert.Error(t, err, "rule %q", rule)
	}
}

func TestNewConfigRelayMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		settings  Settings
		wantNames []string
		wantErr   bool
	}{
		{
			name:      "relay connector with allow rules",
			settings:  Settings{Relay: true, RelayAllow: []string{"10.0.0.0/8", "*.corp.internal:5432"}},
			wantNames: []string{ProxyTargetName},
		},
		{
			name:      "relay connector with deny rules",
			settings:  Settings{Relay: true, RelayDeny: []string{"*.example.com"}},
			wantNames: []string{ProxyTargetName},
		},
		{
			name:     "relay mode alongside named targets",
			settings: Settings{Relay: true, Targets: map[string]string{"database": "db:5432"}},
			wantErr:  true,
		},
		{
			name:     "rules without relay mode",
			settings: Settings{RelayDeny: []string{"*.example.com"}, Targets: map[string]string{"database": "db:5432"}},
			wantErr:  true,
		},
		{
			name:     "target named proxy is reserved",
			settings: Settings{Relay: true, Targets: map[string]string{ProxyTargetName: "db:5432"}},
			wantErr:  true,
		},
		{
			name:    "no targets and no relay",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			settings := test.settings
			settings.ConnectionID = uuid.NewString()
			settings.APIKey = "test-key"
			config, err := NewConfig(settings)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.wantNames, config.Targets.Names())
		})
	}
}

func TestTargetMapHashChangesWithProxyPolicy(t *testing.T) {
	t.Parallel()

	open, err := ParseProxyPolicy(nil, nil)
	require.NoError(t, err)
	restricted, err := ParseProxyPolicy(nil, []string{"*.example.com"})
	require.NoError(t, err)
	openMap, err := newTargetMap(nil, &open)
	require.NoError(t, err)
	restrictedMap, err := newTargetMap(nil, &restricted)
	require.NoError(t, err)

	assert.NotEqual(t, openMap.Hash(), restrictedMap.Hash())
	assert.Equal(t, []transport.TargetDescriptor{{Name: ProxyTargetName, Kind: transport.TargetKindProxy}},
		openMap.RegistrationMetadata().Targets)
}

func TestFlowHandlerProxyFlows(t *testing.T) {
	t.Parallel()

	policy, err := ParseProxyPolicy(nil, nil)
	require.NoError(t, err)
	proxyTargets, err := newTargetMap(nil, &policy)
	require.NoError(t, err)
	tunnelTargets, err := NewTargetMap(map[string]string{"database": "db.internal:5432"})
	require.NoError(t, err)

	tests := []struct {
		name        string
		targets     TargetMap
		target      string
		destination string
		wantStatus  transport.FlowStatus
		wantDialed  []string
	}{
		{name: "allowed destination", targets: proxyTargets, target: ProxyTargetName, destination: "www.reddit.com:443",
			wantStatus: transport.FlowStatusReady, wantDialed: []string{"151.101.1.140:443"}},
		{name: "denied destination", targets: proxyTargets, target: ProxyTargetName, destination: "db.internal:443",
			wantStatus: transport.FlowStatusDestinationDenied},
		{name: "proxy without destination", targets: proxyTargets, target: ProxyTargetName,
			wantStatus: transport.FlowStatusProtocolError},
		{name: "tunnel target with destination", targets: tunnelTargets, target: "database", destination: "www.reddit.com:443",
			wantStatus: transport.FlowStatusProtocolError},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			dialer := &recordingDialer{}
			handler, err := newFlowHandler(test.targets, dialer, 1)
			require.NoError(t, err)
			handler.resolver = staticResolver{
				"www.reddit.com": {"151.101.1.140"},
				"db.internal":    {"10.1.2.3"},
			}

			edgeSide, connectorSide := net.Pipe()
			t.Cleanup(func() { _ = edgeSide.Close() })
			done := make(chan error, 1)
			go func() { done <- handler.handle(context.Background(), connectorSide) }()

			flowID := uuid.Must(uuid.NewV7()).String()
			_, err = transport.WriteFrame(edgeSide, transport.FrameTypeFlowOpen, transport.FlowOpen{
				FlowID:                  flowID,
				TargetName:              test.target,
				Destination:             test.destination,
				DialTimeoutMilliseconds: 1_000,
			})
			require.NoError(t, err)
			var result transport.FlowResult
			require.NoError(t, transport.ReadFrame(edgeSide, transport.FrameTypeFlowResult, &result))
			assert.Equal(t, test.wantStatus, result.Status)
			_ = edgeSide.Close()
			<-done
			assert.Equal(t, test.wantDialed, dialer.addresses)
		})
	}
}

func TestProxyPolicyStringIgnoresCaseOrderAndDuplicates(t *testing.T) {
	t.Parallel()

	first, err := ParseProxyPolicy([]string{"*.Reddit.com", "10.0.0.0/8", "DB.corp.internal:5432"}, []string{"192.168.1.1"})
	require.NoError(t, err)
	second, err := ParseProxyPolicy(
		[]string{"db.corp.internal:5432", "10.0.0.0/8", "*.reddit.com", "*.reddit.com"},
		[]string{"192.168.1.1/32"},
	)
	require.NoError(t, err)
	different, err := ParseProxyPolicy([]string{"*.reddit.com"}, nil)
	require.NoError(t, err)

	assert.Equal(t, first.String(), second.String())
	assert.NotEqual(t, first.String(), different.String())
}
