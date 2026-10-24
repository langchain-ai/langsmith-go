package connector

// newDefaultTargetDialer dials a connector's targets. Reaching private,
// loopback, and Kubernetes-internal services is the purpose of a connector, so
// only cloud metadata is refused, by hostname and by every resolved address at
// connect time.
func newDefaultTargetDialer() TargetDialer {
	return newMetadataGuardDialer()
}
