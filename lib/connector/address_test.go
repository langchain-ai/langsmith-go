package connector

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressClassification(t *testing.T) {
	t.Parallel()

	tests := []struct {
		ip       string
		public   bool
		metadata bool
	}{
		{ip: "93.184.215.14", public: true},
		{ip: "2606:4700:4700::1111", public: true},
		{ip: "10.0.0.5"},
		{ip: "172.16.4.1"},
		{ip: "192.168.1.1"},
		{ip: "100.64.0.1"},
		{ip: "127.0.0.1"},
		{ip: "0.0.0.0"},
		{ip: "169.254.1.1"},
		{ip: "198.18.7.7"},
		{ip: "224.0.0.1"},
		{ip: "255.255.255.255"},
		{ip: "::1"},
		{ip: "fd12:3456::1"},
		{ip: "fe80::1"},
		{ip: "ff02::1"},
		{ip: "::ffff:10.0.0.5"},
		{ip: "::ffff:93.184.215.14", public: true},
		{ip: "64:ff9b::a00:5"},
		{ip: "64:ff9b::5db8:d70e", public: true},
		{ip: "169.254.169.254", metadata: true},
		{ip: "169.254.170.2", metadata: true},
		{ip: "100.100.100.200", metadata: true},
		{ip: "fd00:ec2::254", metadata: true},
		{ip: "::ffff:169.254.169.254", metadata: true},
		{ip: "64:ff9b::a9fe:a9fe", metadata: true},
	}
	for _, tt := range tests {
		t.Run(tt.ip, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			assert.Equal(t, tt.public, isPublicAddress(ip), "public")
			assert.Equal(t, tt.metadata, isMetadataAddress(ip), "metadata")
		})
	}
}

func TestUnparseableAddressesAreRefused(t *testing.T) {
	t.Parallel()

	assert.False(t, isPublicAddress(nil))
	assert.True(t, isMetadataAddress(nil))
}

func TestMetadataHosts(t *testing.T) {
	t.Parallel()

	for host, want := range map[string]bool{
		"metadata.google.internal":  true,
		"Metadata.Google.Internal.": true,
		"metadata.amazonaws.com":    true,
		"metadata":                  true,
		"instance-data":             true,
		"metadata.example.com":      false,
		"db.corp.internal":          false,
	} {
		assert.Equal(t, want, isMetadataHost(host), host)
	}
}
