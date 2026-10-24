package connector

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/langchain-ai/langsmith-go/lib/connector/transport"
	"github.com/stretchr/testify/require"
)

func TestNewTargetMapAcceptsAndCanonicalizesNamesAndAddresses(t *testing.T) {
	t.Parallel()

	targets, err := NewTargetMap(map[string]string{
		"a":                     "LOCALHOST:080",
		"api2":                  "127.0.0.1:443",
		"internal_api_":         "[2001:0db8::1]:5432",
		strings.Repeat("z", 63): "service.internal.:65535",
	})
	require.NoError(t, err)
	require.Equal(t, "localhost:80", mustLookupTarget(t, targets, "a"))
	require.Equal(t, "127.0.0.1:443", mustLookupTarget(t, targets, "api2"))
	require.Equal(t, "[2001:db8::1]:5432", mustLookupTarget(t, targets, "internal_api_"))
	require.Equal(t, "service.internal:65535", mustLookupTarget(t, targets, strings.Repeat("z", 63)))
}

func TestNewTargetMapParsesHTTPOriginsWithoutExposingTCPAddresses(t *testing.T) {
	t.Parallel()

	targets, err := NewTargetMap(map[string]string{
		"payments":  "http://PAYMENTS.default.:080",
		"model_api": "https://model.internal:8443",
		"postgres":  "postgres.database:5432",
	})
	require.NoError(t, err)

	payments, ok := targets.Resolve("payments")
	require.True(t, ok)
	require.Equal(t, "payments.default:80", payments.Address)
	require.Equal(t, transport.TargetDescriptor{
		Name:            "payments",
		Kind:            transport.TargetKindHTTP,
		OriginScheme:    "http",
		OriginAuthority: "payments.default:80",
	}, payments.Descriptor)

	modelAPI, ok := targets.Resolve("model_api")
	require.True(t, ok)
	require.Equal(t, "model.internal:8443", modelAPI.Address)
	require.Equal(t, "https", modelAPI.Descriptor.OriginScheme)

	postgres, ok := targets.Resolve("postgres")
	require.True(t, ok)
	require.Equal(t, transport.TargetKindTCP, postgres.Descriptor.Kind)
	require.Empty(t, postgres.Descriptor.OriginScheme)
	require.Empty(t, postgres.Descriptor.OriginAuthority)

	metadata := targets.RegistrationMetadata()
	require.Equal(t, []transport.TargetDescriptor{
		{
			Name:            "model_api",
			Kind:            transport.TargetKindHTTP,
			OriginScheme:    "https",
			OriginAuthority: "model.internal:8443",
		},
		{
			Name:            "payments",
			Kind:            transport.TargetKindHTTP,
			OriginScheme:    "http",
			OriginAuthority: "payments.default:80",
		},
		{Name: "postgres", Kind: transport.TargetKindTCP},
	}, metadata.Targets)
}

func TestNewTargetMapRejectsInvalidNames(t *testing.T) {
	t.Parallel()

	names := []string{
		"",
		"Database",
		"1database",
		"api-v2",
		"database.primary",
		"database/primary",
		"database primary",
		"-database",
		"_database",
		strings.Repeat("a", 64),
		"d\u00e1tabase",
		"*",
	}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := NewTargetMap(map[string]string{name: "localhost:5432"})
			require.Error(t, err)
		})
	}
}

func TestNewTargetMapRejectsMoreThan256Targets(t *testing.T) {
	t.Parallel()

	raw := make(map[string]string, 257)
	for index := range 257 {
		raw[fmt.Sprintf("target_%03d", index)] = "localhost:5432"
	}

	_, err := NewTargetMap(raw)
	require.Error(t, err)
}

func TestNewTargetMapRejectsInvalidAddresses(t *testing.T) {
	t.Parallel()

	addresses := []string{
		"",
		"postgres://db.internal:5432",
		"http://db.internal",
		"http://db.internal:80/path",
		"http://user@db.internal:80",
		"http://db.internal:80?query=yes",
		"http://db.internal:80#fragment",
		"ftp://db.internal:21",
		"db.internal/path:5432",
		"*:5432",
		"*.internal:5432",
		":5432",
		"db.internal",
		"db.internal:0",
		"db.internal:65536",
		"db.internal:not-a-port",
		"db.internal:+5432",
		"db.internal:-1",
		"2001:db8::1:5432",
		"[not-an-ip]:5432",
		" db.internal:5432",
		"db.internal:5432 ",
		"user@db.internal:5432",
		"/tmp/database.sock:5432",
		"_postgres._tcp:5432",
		"-db.internal:5432",
		"db-.internal:5432",
		"db..internal:5432",
		strings.Repeat("a", 64) + ".internal:5432",
		strings.Repeat("a", 250) + ".com:5432",
		"d\u00e1tabase.internal:5432",
		"999.999.999.999:5432",
		"0.0.0.0:5432",
		"[::]:5432",
	}
	for _, address := range addresses {
		t.Run(address, func(t *testing.T) {
			t.Parallel()

			_, err := NewTargetMap(map[string]string{"database": address})
			require.Error(t, err)
		})
	}
}

func TestTargetMapHashIsDeterministicForCanonicalMap(t *testing.T) {
	t.Parallel()

	first, err := NewTargetMap(map[string]string{
		"database": "DB.INTERNAL:05432",
		"api":      "api.internal:8443",
	})
	require.NoError(t, err)
	second, err := NewTargetMap(map[string]string{
		"api":      "API.INTERNAL:8443",
		"database": "db.internal:5432",
	})
	require.NoError(t, err)

	require.Equal(t, first.Hash(), second.Hash())
	require.Regexp(t, regexp.MustCompile(`^[0-9a-f]{64}$`), first.Hash())
}

func TestTargetMapCopiesInputAndReturnedNames(t *testing.T) {
	t.Parallel()

	raw := map[string]string{"database": "db.internal:5432", "api": "api.internal:8443"}
	targets, err := NewTargetMap(raw)
	require.NoError(t, err)

	raw["database"] = "other.internal:1234"
	firstNames := targets.Names()
	require.Equal(t, []string{"api", "database"}, firstNames)
	firstNames[0] = "modified"

	require.Equal(t, "db.internal:5432", mustLookupTarget(t, targets, "database"))
	require.Equal(t, []string{"api", "database"}, targets.Names())
}

func TestRegistrationMetadataExposesOnlyHTTPRoutingAuthorities(t *testing.T) {
	t.Parallel()

	targets, err := NewTargetMap(map[string]string{
		"database": "db.internal:5432",
		"api":      "https://api.internal:8443",
	})
	require.NoError(t, err)

	metadata := targets.RegistrationMetadata()
	require.Equal(t, targets.Hash(), metadata.TargetMapHash)
	require.Equal(t, []transport.TargetDescriptor{
		{
			Name:            "api",
			Kind:            transport.TargetKindHTTP,
			OriginScheme:    "https",
			OriginAuthority: "api.internal:8443",
		},
		{Name: "database", Kind: transport.TargetKindTCP},
	}, metadata.Targets)

	serialized, err := json.Marshal(metadata)
	require.NoError(t, err)
	require.NotContains(t, string(serialized), "db.internal")
	require.NotContains(t, string(serialized), "5432")
	require.Contains(t, string(serialized), "api.internal:8443")
}

func TestRegistrationMetadataReturnsDefensiveDescriptorCopy(t *testing.T) {
	t.Parallel()

	targets, err := NewTargetMap(map[string]string{"database": "db.internal:5432"})
	require.NoError(t, err)

	metadata := targets.RegistrationMetadata()
	metadata.Targets[0].Name = "modified"

	require.Equal(t, "database", targets.RegistrationMetadata().Targets[0].Name)
}
