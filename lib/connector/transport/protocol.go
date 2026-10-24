package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const (
	// ProtocolVersion is the current carrier control protocol version.
	ProtocolVersion byte = 1

	// WebSocketSubprotocol prevents non-carrier WebSocket clients from being
	// mistaken for connector carrier sessions.
	WebSocketSubprotocol = "langsmith-connector.v1"

	HeaderConnectionID        = "X-LangSmith-Connection-Id"
	HeaderConnectorInstanceID = "X-LangSmith-Connector-Instance-Id"
	HeaderSessionID           = "X-LangSmith-Session-Id"
	// HeaderConnectorKind is "tunnel" or "relay": what the connector serves,
	// which must match the connection's kind. A connector that omits it
	// serves a tunnel.
	HeaderConnectorKind = "X-LangSmith-Connector-Kind"
	// HeaderCaller names a first-party caller whose service key the edge has
	// verified on a dedicated control-plane route. Sandbox hosts sign with
	// their own secret.
	HeaderCaller = "X-LangSmith-Caller"
	// CallerSandboxHost is HeaderCaller for a sandbox host.
	CallerSandboxHost = "sandbox-host"

	MaxRegisteredTargets = 256
	MaxTargetNameLength  = 63
)

// TargetKind describes how the connector prepares a named target before it
// reports FLOW_READY. HTTP targets carry the bounded routing metadata the edge
// needs to proxy requests; TCP dial addresses remain connector-local. A proxy
// target has no fixed address: each flow names its destination, and the
// connector's own relay policy decides whether to dial it.
type TargetKind string

const (
	TargetKindTCP   TargetKind = "tcp"
	TargetKindHTTP  TargetKind = "http"
	TargetKindProxy TargetKind = "proxy"
)

// TargetDescriptor is the metadata an edge needs to expose an HTTP target.
// OriginAuthority is deliberately limited to a canonical host:port;
// credentials, paths, queries, and fragments never cross the carrier
// registration boundary.
type TargetDescriptor struct {
	Name            string     `json:"name"`
	Kind            TargetKind `json:"kind"`
	OriginScheme    string     `json:"origin_scheme,omitempty"`
	OriginAuthority string     `json:"origin_authority,omitempty"`
}

// FrameType identifies a control message without exposing application data.
type FrameType byte

const (
	FrameTypeRegister FrameType = iota + 1
	FrameTypeRegistered
	FrameTypeFlowOpen
	FrameTypeFlowResult
	FrameTypeDrain
)

// Register is the first control-stream message sent by a connector.
type Register struct {
	ProtocolVersion     byte               `json:"protocol_version"`
	ConnectionID        string             `json:"connection_id"`
	ConnectorInstanceID string             `json:"connector_instance_id"`
	SessionID           string             `json:"session_id"`
	Targets             []TargetDescriptor `json:"targets"`
	TargetMapHash       string             `json:"target_map_hash"`
}

// Registered confirms the exact carrier identity and assigns its edge runtime.
type Registered struct {
	ProtocolVersion     byte   `json:"protocol_version"`
	ConnectionID        string `json:"connection_id"`
	ConnectorInstanceID string `json:"connector_instance_id"`
	SessionID           string `json:"session_id"`
	RuntimeID           string `json:"runtime_id"`
}

// FlowOpen requests one opaque TCP flow to a previously registered target.
// Destination is the host:port a proxy target dials; it is empty for every
// other kind of target.
type FlowOpen struct {
	FlowID                  string `json:"flow_id"`
	TargetName              string `json:"target_name"`
	DialTimeoutMilliseconds uint32 `json:"dial_timeout_milliseconds"`
	Destination             string `json:"destination,omitempty"`
}

// FlowStatus is deliberately closed so connector failures cannot leak target
// addresses, credentials, or arbitrary error text to the edge.
type FlowStatus string

const (
	FlowStatusReady            FlowStatus = "ready"
	FlowStatusUnknownTarget    FlowStatus = "unknown_target"
	FlowStatusDialTimeout      FlowStatus = "dial_timeout"
	FlowStatusDialFailed       FlowStatus = "dial_failed"
	FlowStatusCapacityExceeded FlowStatus = "capacity_exceeded"
	FlowStatusProtocolError    FlowStatus = "protocol_error"
	// FlowStatusDestinationDenied means the connector's relay policy refused
	// the requested destination.
	FlowStatusDestinationDenied FlowStatus = "destination_denied"
)

var (
	ErrInvalidFlowStatus      = errors.New("invalid flow status")
	ErrInvalidProtocolMessage = errors.New("invalid carrier protocol message")
)

// FlowResult is the connector's bounded response to a flow request.
type FlowResult struct {
	FlowID string     `json:"flow_id"`
	Status FlowStatus `json:"status"`
}

// Drain is sent by a connector on its control stream to ask the edge to stop
// routing new flows to this carrier while its in-flight flows finish, so a
// replacement carrier can take over without dropping traffic.
type Drain struct {
	ProtocolVersion     byte   `json:"protocol_version"`
	ConnectionID        string `json:"connection_id"`
	ConnectorInstanceID string `json:"connector_instance_id"`
	SessionID           string `json:"session_id"`
}

func (s FlowStatus) valid() bool {
	switch s {
	case FlowStatusReady,
		FlowStatusUnknownTarget,
		FlowStatusDialTimeout,
		FlowStatusDialFailed,
		FlowStatusCapacityExceeded,
		FlowStatusProtocolError,
		FlowStatusDestinationDenied:
		return true
	default:
		return false
	}
}

func (s FlowStatus) MarshalJSON() ([]byte, error) {
	if !s.valid() {
		return nil, ErrInvalidFlowStatus
	}
	return json.Marshal(string(s))
}

func (s *FlowStatus) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	status := FlowStatus(value)
	if !status.valid() {
		return ErrInvalidFlowStatus
	}
	*s = status
	return nil
}

func validFrameType(frameType FrameType) bool {
	return frameType >= FrameTypeRegister && frameType <= FrameTypeDrain
}

// ValidateTargetName applies the address-blind grammar shared by connector
// configuration and edge admission.
func ValidateTargetName(name string) error {
	if len(name) == 0 || len(name) > MaxTargetNameLength {
		return invalidProtocolField("target_name")
	}
	if !isLowerLetter(name[0]) {
		return invalidProtocolField("target_name")
	}
	for i := 1; i < len(name); i++ {
		c := name[i]
		if !isLowerLetter(c) && !(c >= '0' && c <= '9') && c != '_' {
			return invalidProtocolField("target_name")
		}
	}
	return nil
}

func ValidateRegister(message Register) error {
	if message.ProtocolVersion != ProtocolVersion {
		return invalidProtocolField("protocol_version")
	}
	if err := validateCanonicalUUID(message.ConnectionID); err != nil {
		return invalidProtocolField("connection_id")
	}
	if err := validateCanonicalUUID(message.ConnectorInstanceID); err != nil {
		return invalidProtocolField("connector_instance_id")
	}
	if err := validateCanonicalUUID(message.SessionID); err != nil {
		return invalidProtocolField("session_id")
	}
	if len(message.Targets) == 0 || len(message.Targets) > MaxRegisteredTargets {
		return invalidProtocolField("targets")
	}
	targets := make(map[string]struct{}, len(message.Targets))
	for _, descriptor := range message.Targets {
		if validateTargetDescriptor(descriptor) != nil {
			return invalidProtocolField("targets")
		}
		if _, duplicate := targets[descriptor.Name]; duplicate {
			return invalidProtocolField("targets")
		}
		targets[descriptor.Name] = struct{}{}
	}
	if !validTargetMapHash(message.TargetMapHash) {
		return invalidProtocolField("target_map_hash")
	}
	return nil
}

func validateTargetDescriptor(descriptor TargetDescriptor) error {
	if err := ValidateTargetName(descriptor.Name); err != nil {
		return err
	}
	switch descriptor.Kind {
	case TargetKindTCP:
		if descriptor.OriginScheme != "" || descriptor.OriginAuthority != "" {
			return invalidProtocolField("targets")
		}
		return nil
	case TargetKindHTTP:
		if descriptor.OriginScheme != "http" && descriptor.OriginScheme != "https" {
			return invalidProtocolField("targets")
		}
		return validateOriginAuthority(descriptor.OriginAuthority)
	case TargetKindProxy:
		if descriptor.OriginScheme != "" || descriptor.OriginAuthority != "" {
			return invalidProtocolField("targets")
		}
		return nil
	default:
		return invalidProtocolField("targets")
	}
}

func validateOriginAuthority(authority string) error {
	if authority == "" || strings.TrimSpace(authority) != authority {
		return invalidProtocolField("origin_authority")
	}
	host, rawPort, err := net.SplitHostPort(authority)
	if err != nil || host == "" || rawPort == "" {
		return invalidProtocolField("origin_authority")
	}
	port, err := strconv.Atoi(rawPort)
	if err != nil || port < 1 || port > 65535 || strconv.Itoa(port) != rawPort {
		return invalidProtocolField("origin_authority")
	}
	if ip := net.ParseIP(host); ip != nil {
		if ip.IsUnspecified() || net.JoinHostPort(ip.String(), rawPort) != authority {
			return invalidProtocolField("origin_authority")
		}
		return nil
	}
	if len(host) > 253 || host != strings.ToLower(host) || strings.HasSuffix(host, ".") {
		return invalidProtocolField("origin_authority")
	}
	for _, label := range strings.Split(host, ".") {
		if len(label) == 0 || len(label) > 63 || !isAlphaNumeric(label[0]) ||
			!isAlphaNumeric(label[len(label)-1]) {
			return invalidProtocolField("origin_authority")
		}
		for index := range len(label) {
			if !isAlphaNumeric(label[index]) && label[index] != '-' {
				return invalidProtocolField("origin_authority")
			}
		}
	}
	return nil
}

func isAlphaNumeric(char byte) bool {
	return char >= 'a' && char <= 'z' || char >= '0' && char <= '9'
}

func ValidateRegistered(message Registered) error {
	if message.ProtocolVersion != ProtocolVersion {
		return invalidProtocolField("protocol_version")
	}
	if err := validateCanonicalUUID(message.ConnectionID); err != nil {
		return invalidProtocolField("connection_id")
	}
	if err := validateCanonicalUUID(message.ConnectorInstanceID); err != nil {
		return invalidProtocolField("connector_instance_id")
	}
	if err := validateCanonicalUUID(message.SessionID); err != nil {
		return invalidProtocolField("session_id")
	}
	if err := validateCanonicalUUID(message.RuntimeID); err != nil {
		return invalidProtocolField("runtime_id")
	}
	return nil
}

func ValidateFlowOpen(message FlowOpen) error {
	if err := validateCanonicalUUID(message.FlowID); err != nil {
		return invalidProtocolField("flow_id")
	}
	if err := ValidateTargetName(message.TargetName); err != nil {
		return invalidProtocolField("target_name")
	}
	if message.DialTimeoutMilliseconds == 0 {
		return invalidProtocolField("dial_timeout_milliseconds")
	}
	if message.Destination != "" {
		if err := ValidateDestination(message.Destination); err != nil {
			return invalidProtocolField("destination")
		}
	}
	return nil
}

// ValidateDestination checks that destination is a canonical host:port: a
// lowercase DNS name or an IP literal, and a port from 1 to 65535.
func ValidateDestination(destination string) error {
	return validateOriginAuthority(destination)
}

func ValidateFlowResult(message FlowResult) error {
	if err := validateCanonicalUUID(message.FlowID); err != nil {
		return invalidProtocolField("flow_id")
	}
	if !message.Status.valid() {
		return invalidProtocolField("status")
	}
	return nil
}

func ValidateDrain(message Drain) error {
	if message.ProtocolVersion != ProtocolVersion {
		return invalidProtocolField("protocol_version")
	}
	if err := validateCanonicalUUID(message.ConnectionID); err != nil {
		return invalidProtocolField("connection_id")
	}
	if err := validateCanonicalUUID(message.ConnectorInstanceID); err != nil {
		return invalidProtocolField("connector_instance_id")
	}
	if err := validateCanonicalUUID(message.SessionID); err != nil {
		return invalidProtocolField("session_id")
	}
	return nil
}

func validateCanonicalUUID(value string) error {
	parsed, err := uuid.Parse(value)
	if err != nil || parsed.String() != value {
		return ErrInvalidProtocolMessage
	}
	return nil
}

func validTargetMapHash(value string) bool {
	if len(value) != 64 {
		return false
	}
	for i := range len(value) {
		c := value[i]
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}

func invalidProtocolField(field string) error {
	return fmt.Errorf("%w: %s", ErrInvalidProtocolMessage, field)
}

func isLowerLetter(character byte) bool {
	return character >= 'a' && character <= 'z'
}
