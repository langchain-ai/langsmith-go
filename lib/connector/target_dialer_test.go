package connector

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDefaultTargetDialerReachesPrivateAddresses(t *testing.T) {
	t.Parallel()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	conn, err := newDefaultTargetDialer().DialContext(t.Context(), "tcp", listener.Addr().String())
	require.NoError(t, err)
	require.NoError(t, conn.Close())
}

func TestDefaultTargetDialerRefusesMetadata(t *testing.T) {
	t.Parallel()

	for _, address := range []string{
		"169.254.169.254:80",
		"[::ffff:169.254.169.254]:80",
		"[fd00:ec2::254]:80",
		"metadata.google.internal:80",
		"METADATA.google.internal.:80",
	} {
		t.Run(address, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), time.Second)
			defer cancel()
			_, err := newDefaultTargetDialer().DialContext(ctx, "tcp", address)
			require.ErrorIs(t, err, errMetadataDestination)
			require.False(t, errors.Is(err, context.DeadlineExceeded))
		})
	}
}
