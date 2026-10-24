package transport

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestFrameRoundTrips(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		frameType FrameType
		message   any
		newResult func() any
	}{
		{
			name:      "register",
			frameType: FrameTypeRegister,
			message: Register{
				ProtocolVersion:     ProtocolVersion,
				ConnectionID:        "0190f9d8-4d7a-7c54-8c74-4ee0049f07e0",
				ConnectorInstanceID: "0190f9d8-4d7a-7c54-8c74-4ee0049f07e1",
				SessionID:           "0190f9d8-4d7a-7c54-8c74-4ee0049f07e2",
				Targets: []TargetDescriptor{
					{Name: "primary_db", Kind: TargetKindTCP},
					{
						Name:            "redis",
						Kind:            TargetKindHTTP,
						OriginScheme:    "https",
						OriginAuthority: "redis.internal:8443",
					},
				},
				TargetMapHash: "89d9a681336ba950161134b0de56a43d67b1e2cc6eaa1371c8e3ba34555fe560",
			},
			newResult: func() any { return &Register{} },
		},
		{
			name:      "registered",
			frameType: FrameTypeRegistered,
			message: Registered{
				ProtocolVersion:     ProtocolVersion,
				ConnectionID:        "0190f9d8-4d7a-7c54-8c74-4ee0049f07e0",
				ConnectorInstanceID: "0190f9d8-4d7a-7c54-8c74-4ee0049f07e1",
				SessionID:           "0190f9d8-4d7a-7c54-8c74-4ee0049f07e2",
				RuntimeID:           "0190f9d8-4d7a-7c54-8c74-4ee0049f07e3",
			},
			newResult: func() any { return &Registered{} },
		},
		{
			name:      "flow open",
			frameType: FrameTypeFlowOpen,
			message: FlowOpen{
				FlowID:                  "0190f9d8-4d7a-7c54-8c74-4ee0049f07e4",
				TargetName:              "primary_db",
				DialTimeoutMilliseconds: 5_000,
			},
			newResult: func() any { return &FlowOpen{} },
		},
		{
			name:      "flow result",
			frameType: FrameTypeFlowResult,
			message: FlowResult{
				FlowID: "0190f9d8-4d7a-7c54-8c74-4ee0049f07e4",
				Status: FlowStatusReady,
			},
			newResult: func() any { return &FlowResult{} },
		},
		{
			name:      "drain",
			frameType: FrameTypeDrain,
			message: Drain{
				ProtocolVersion:     ProtocolVersion,
				ConnectionID:        "0190f9d8-4d7a-7c54-8c74-4ee0049f07e0",
				ConnectorInstanceID: "0190f9d8-4d7a-7c54-8c74-4ee0049f07e1",
				SessionID:           "0190f9d8-4d7a-7c54-8c74-4ee0049f07e2",
			},
			newResult: func() any { return &Drain{} },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var wire bytes.Buffer
			written, err := WriteFrame(&wire, tt.frameType, tt.message)
			if err != nil {
				t.Fatalf("WriteFrame() error = %v", err)
			}
			if written != wire.Len() {
				t.Fatalf("WriteFrame() wrote %d bytes, buffer contains %d", written, wire.Len())
			}

			got := tt.newResult()
			if err := ReadFrame(&wire, tt.frameType, got); err != nil {
				t.Fatalf("ReadFrame() error = %v", err)
			}
			if !reflect.DeepEqual(got, pointerTo(tt.message)) {
				t.Fatalf("ReadFrame() = %#v, want %#v", got, tt.message)
			}
		})
	}
}

func TestReadFrameHandlesShortReads(t *testing.T) {
	t.Parallel()

	message := FlowOpen{
		FlowID:                  "0190f9d8-4d7a-7c54-8c74-4ee0049f07e4",
		TargetName:              "primary_db",
		DialTimeoutMilliseconds: 5_000,
	}
	var wire bytes.Buffer
	if _, err := WriteFrame(&wire, FrameTypeFlowOpen, message); err != nil {
		t.Fatalf("WriteFrame() error = %v", err)
	}

	var got FlowOpen
	if err := ReadFrame(&chunkReader{reader: bytes.NewReader(wire.Bytes()), max: 1}, FrameTypeFlowOpen, &got); err != nil {
		t.Fatalf("ReadFrame() error = %v", err)
	}
	if !reflect.DeepEqual(got, message) {
		t.Fatalf("ReadFrame() = %#v, want %#v", got, message)
	}
}

func TestWriteFrameHandlesShortWrites(t *testing.T) {
	t.Parallel()

	message := FlowResult{
		FlowID: "0190f9d8-4d7a-7c54-8c74-4ee0049f07e4",
		Status: FlowStatusReady,
	}
	writer := &chunkWriter{max: 2}
	written, err := WriteFrame(writer, FrameTypeFlowResult, message)
	if err != nil {
		t.Fatalf("WriteFrame() error = %v", err)
	}
	if written != writer.buffer.Len() {
		t.Fatalf("WriteFrame() wrote %d bytes, buffer contains %d", written, writer.buffer.Len())
	}

	var got FlowResult
	if err := ReadFrame(&writer.buffer, FrameTypeFlowResult, &got); err != nil {
		t.Fatalf("ReadFrame() error = %v", err)
	}
	if !reflect.DeepEqual(got, message) {
		t.Fatalf("ReadFrame() = %#v, want %#v", got, message)
	}
}

func TestWriteFrameReportsPartialWrite(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("connection lost")
	writer := &failingWriter{limit: 5, err: wantErr}
	written, err := WriteFrame(writer, FrameTypeFlowOpen, FlowOpen{
		FlowID:                  "0190f9d8-4d7a-7c54-8c74-4ee0049f07e4",
		TargetName:              "primary_db",
		DialTimeoutMilliseconds: 5_000,
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("WriteFrame() error = %v, want %v", err, wantErr)
	}
	if written != writer.limit {
		t.Fatalf("WriteFrame() wrote %d bytes, want %d", written, writer.limit)
	}
}

func TestReadFrameRejectsInvalidHeaders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		version      byte
		frameType    byte
		expectedType FrameType
		wantErr      error
	}{
		{
			name:         "unknown version",
			version:      99,
			frameType:    byte(FrameTypeRegister),
			expectedType: FrameTypeRegister,
			wantErr:      ErrUnsupportedVersion,
		},
		{
			name:         "unknown type",
			version:      ProtocolVersion,
			frameType:    99,
			expectedType: FrameTypeRegister,
			wantErr:      ErrUnknownFrameType,
		},
		{
			name:         "unexpected type",
			version:      ProtocolVersion,
			frameType:    byte(FrameTypeRegistered),
			expectedType: FrameTypeRegister,
			wantErr:      ErrUnexpectedFrameType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			wire := rawFrame(tt.version, tt.frameType, []byte(`{}`))
			var result Register
			err := ReadFrame(bytes.NewReader(wire), tt.expectedType, &result)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ReadFrame() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriteFrameRejectsUnknownType(t *testing.T) {
	t.Parallel()

	written, err := WriteFrame(io.Discard, FrameType(99), Register{})
	if !errors.Is(err, ErrUnknownFrameType) {
		t.Fatalf("WriteFrame() error = %v, want %v", err, ErrUnknownFrameType)
	}
	if written != 0 {
		t.Fatalf("WriteFrame() wrote %d bytes, want 0", written)
	}
}

func TestReadFrameRejectsMalformedOrTruncatedData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wire    []byte
		wantErr error
	}{
		{
			name:    "truncated header",
			wire:    []byte{ProtocolVersion, byte(FrameTypeRegister)},
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "truncated payload",
			wire:    rawFrameWithDeclaredLength(ProtocolVersion, byte(FrameTypeRegister), []byte(`{}`), 3),
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "malformed JSON",
			wire:    rawFrame(ProtocolVersion, byte(FrameTypeRegister), []byte(`{"connection_id":`)),
			wantErr: ErrMalformedFrame,
		},
		{
			name:    "trailing JSON",
			wire:    rawFrame(ProtocolVersion, byte(FrameTypeRegister), []byte(`{} {}`)),
			wantErr: ErrMalformedFrame,
		},
		{
			name:    "unknown JSON field",
			wire:    rawFrame(ProtocolVersion, byte(FrameTypeRegister), []byte(`{"credential-that-must-not-leak":true}`)),
			wantErr: ErrMalformedFrame,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var result Register
			err := ReadFrame(bytes.NewReader(tt.wire), FrameTypeRegister, &result)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ReadFrame() error = %v, want %v", err, tt.wantErr)
			}
			if strings.Contains(err.Error(), "credential-that-must-not-leak") {
				t.Fatalf("ReadFrame() error leaked peer-controlled content: %v", err)
			}
		})
	}
}

func TestControlFrameSizeLimit(t *testing.T) {
	t.Parallel()

	t.Run("read rejects declared oversized payload before allocation", func(t *testing.T) {
		t.Parallel()

		header := make([]byte, frameHeaderSizeBytes)
		header[0] = ProtocolVersion
		header[1] = byte(FrameTypeRegister)
		binary.BigEndian.PutUint32(header[2:], MaxControlFrameSizeBytes+1)
		var result Register
		err := ReadFrame(bytes.NewReader(header), FrameTypeRegister, &result)
		if !errors.Is(err, ErrFrameTooLarge) {
			t.Fatalf("ReadFrame() error = %v, want %v", err, ErrFrameTooLarge)
		}
	})

	t.Run("write rejects oversized payload without writing", func(t *testing.T) {
		t.Parallel()

		var wire bytes.Buffer
		written, err := WriteFrame(&wire, FrameTypeRegister, Register{
			Targets: []TargetDescriptor{{
				Name: string(bytes.Repeat([]byte{'a'}, MaxControlFrameSizeBytes)),
				Kind: TargetKindTCP,
			}},
		})
		if !errors.Is(err, ErrFrameTooLarge) {
			t.Fatalf("WriteFrame() error = %v, want %v", err, ErrFrameTooLarge)
		}
		if written != 0 || wire.Len() != 0 {
			t.Fatalf("WriteFrame() wrote %d bytes and buffered %d, want no output", written, wire.Len())
		}
	})
}

func TestFlowResultRejectsUnknownStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		message FlowResult
		wire    []byte
	}{
		{
			name: "write",
			message: FlowResult{
				FlowID: "0190f9d8-4d7a-7c54-8c74-4ee0049f07e4",
				Status: FlowStatus("surprise"),
			},
		},
		{
			name: "read",
			wire: rawFrame(
				ProtocolVersion,
				byte(FrameTypeFlowResult),
				[]byte(`{"flow_id":"0190f9d8-4d7a-7c54-8c74-4ee0049f07e4","status":"surprise"}`),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.wire == nil {
				written, err := WriteFrame(io.Discard, FrameTypeFlowResult, tt.message)
				if !errors.Is(err, ErrInvalidFlowStatus) {
					t.Fatalf("WriteFrame() error = %v, want %v", err, ErrInvalidFlowStatus)
				}
				if strings.Contains(err.Error(), "surprise") {
					t.Fatalf("WriteFrame() error leaked invalid status: %v", err)
				}
				if written != 0 {
					t.Fatalf("WriteFrame() wrote %d bytes, want 0", written)
				}
				return
			}

			var result FlowResult
			err := ReadFrame(bytes.NewReader(tt.wire), FrameTypeFlowResult, &result)
			if !errors.Is(err, ErrInvalidFlowStatus) {
				t.Fatalf("ReadFrame() error = %v, want %v", err, ErrInvalidFlowStatus)
			}
			if strings.Contains(err.Error(), "surprise") {
				t.Fatalf("ReadFrame() error leaked invalid status: %v", err)
			}
		})
	}
}

func TestValidateTargetName(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"a", "database", "api2", "internal_api", strings.Repeat("z", 63)} {
		t.Run("valid_"+name, func(t *testing.T) {
			t.Parallel()
			if err := ValidateTargetName(name); err != nil {
				t.Fatalf("ValidateTargetName(%q) error = %v", name, err)
			}
		})
	}

	for _, name := range []string{
		"",
		"Database",
		"1database",
		"_database",
		"api-v2",
		"database.primary",
		"database primary",
		strings.Repeat("z", 64),
		"d\u00e1tabase",
	} {
		t.Run("invalid_"+name, func(t *testing.T) {
			t.Parallel()
			if err := ValidateTargetName(name); !errors.Is(err, ErrInvalidProtocolMessage) {
				t.Fatalf("ValidateTargetName(%q) error = %v, want %v", name, err, ErrInvalidProtocolMessage)
			}
		})
	}
}

func TestValidateRegister(t *testing.T) {
	t.Parallel()

	valid := validRegisterMessage()
	if err := ValidateRegister(valid); err != nil {
		t.Fatalf("ValidateRegister() error = %v", err)
	}
	missingDescriptors := valid
	missingDescriptors.Targets = nil
	if err := ValidateRegister(missingDescriptors); !errors.Is(err, ErrInvalidProtocolMessage) {
		t.Fatalf("ValidateRegister() error = %v, want %v", err, ErrInvalidProtocolMessage)
	}

	tests := []struct {
		name   string
		mutate func(*Register)
	}{
		{name: "version", mutate: func(message *Register) { message.ProtocolVersion = 2 }},
		{name: "connection ID", mutate: func(message *Register) { message.ConnectionID = "not-a-uuid" }},
		{name: "connector instance ID", mutate: func(message *Register) { message.ConnectorInstanceID = "not-a-uuid" }},
		{name: "session ID", mutate: func(message *Register) { message.SessionID = "not-a-uuid" }},
		{name: "no targets", mutate: func(message *Register) { message.Targets = nil }},
		{name: "too many targets", mutate: func(message *Register) {
			message.Targets = make([]TargetDescriptor, MaxRegisteredTargets+1)
			for i := range message.Targets {
				message.Targets[i] = TargetDescriptor{
					Name: "target_" + threeDigits(i),
					Kind: TargetKindTCP,
				}
			}
		}},
		{name: "invalid target", mutate: func(message *Register) {
			message.Targets[0].Name = "api-v2"
		}},
		{name: "duplicate target", mutate: func(message *Register) {
			message.Targets[1].Name = message.Targets[0].Name
		}},
		{name: "invalid descriptor kind", mutate: func(message *Register) {
			message.Targets[0].Kind = TargetKind("udp")
		}},
		{name: "tcp descriptor includes authority", mutate: func(message *Register) {
			message.Targets[0].OriginAuthority = "database.internal:5432"
		}},
		{name: "http descriptor missing scheme", mutate: func(message *Register) {
			message.Targets[1].OriginScheme = ""
		}},
		{name: "http descriptor invalid authority", mutate: func(message *Register) {
			message.Targets[1].OriginAuthority = "bad authority"
		}},
		{name: "short hash", mutate: func(message *Register) { message.TargetMapHash = "abcd" }},
		{name: "uppercase hash", mutate: func(message *Register) { message.TargetMapHash = strings.Repeat("A", 64) }},
		{name: "non hexadecimal hash", mutate: func(message *Register) { message.TargetMapHash = strings.Repeat("z", 64) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			message := validRegisterMessage()
			tt.mutate(&message)
			if err := ValidateRegister(message); !errors.Is(err, ErrInvalidProtocolMessage) {
				t.Fatalf("ValidateRegister() error = %v, want %v", err, ErrInvalidProtocolMessage)
			}
		})
	}
}

func TestValidateRegistered(t *testing.T) {
	t.Parallel()

	valid := Registered{
		ProtocolVersion:     ProtocolVersion,
		ConnectionID:        "0190f9d8-4d7a-7c54-8c74-4ee0049f07e0",
		ConnectorInstanceID: "0190f9d8-4d7a-7c54-8c74-4ee0049f07e1",
		SessionID:           "0190f9d8-4d7a-7c54-8c74-4ee0049f07e2",
		RuntimeID:           "0190f9d8-4d7a-7c54-8c74-4ee0049f07e3",
	}
	if err := ValidateRegistered(valid); err != nil {
		t.Fatalf("ValidateRegistered() error = %v", err)
	}

	valid.RuntimeID = "not-a-uuid"
	if err := ValidateRegistered(valid); !errors.Is(err, ErrInvalidProtocolMessage) {
		t.Fatalf("ValidateRegistered() error = %v, want %v", err, ErrInvalidProtocolMessage)
	}
}

func TestValidateFlowMessages(t *testing.T) {
	t.Parallel()

	flowOpen := FlowOpen{
		FlowID:                  "0190f9d8-4d7a-7c54-8c74-4ee0049f07e4",
		TargetName:              "primary_db",
		DialTimeoutMilliseconds: 5_000,
	}
	if err := ValidateFlowOpen(flowOpen); err != nil {
		t.Fatalf("ValidateFlowOpen() error = %v", err)
	}
	flowOpen.DialTimeoutMilliseconds = 0
	if err := ValidateFlowOpen(flowOpen); !errors.Is(err, ErrInvalidProtocolMessage) {
		t.Fatalf("ValidateFlowOpen() error = %v, want %v", err, ErrInvalidProtocolMessage)
	}

	flowResult := FlowResult{
		FlowID: "0190f9d8-4d7a-7c54-8c74-4ee0049f07e4",
		Status: FlowStatusReady,
	}
	if err := ValidateFlowResult(flowResult); err != nil {
		t.Fatalf("ValidateFlowResult() error = %v", err)
	}
	flowResult.Status = FlowStatus("surprise")
	if err := ValidateFlowResult(flowResult); !errors.Is(err, ErrInvalidProtocolMessage) {
		t.Fatalf("ValidateFlowResult() error = %v, want %v", err, ErrInvalidProtocolMessage)
	}
}

func validRegisterMessage() Register {
	return Register{
		ProtocolVersion:     ProtocolVersion,
		ConnectionID:        "0190f9d8-4d7a-7c54-8c74-4ee0049f07e0",
		ConnectorInstanceID: "0190f9d8-4d7a-7c54-8c74-4ee0049f07e1",
		SessionID:           "0190f9d8-4d7a-7c54-8c74-4ee0049f07e2",
		Targets: []TargetDescriptor{
			{Name: "primary_db", Kind: TargetKindTCP},
			{
				Name:            "redis",
				Kind:            TargetKindHTTP,
				OriginScheme:    "https",
				OriginAuthority: "redis.internal:8443",
			},
		},
		TargetMapHash: "89d9a681336ba950161134b0de56a43d67b1e2cc6eaa1371c8e3ba34555fe560",
	}
}

func threeDigits(value int) string {
	return string([]byte{
		byte('0' + value/100%10),
		byte('0' + value/10%10),
		byte('0' + value%10),
	})
}

func pointerTo(value any) any {
	switch value := value.(type) {
	case Register:
		return &value
	case Registered:
		return &value
	case FlowOpen:
		return &value
	case FlowResult:
		return &value
	case Drain:
		return &value
	default:
		panic("unsupported test value")
	}
}

func rawFrame(version, frameType byte, payload []byte) []byte {
	return rawFrameWithDeclaredLength(version, frameType, payload, uint32(len(payload)))
}

func rawFrameWithDeclaredLength(version, frameType byte, payload []byte, length uint32) []byte {
	wire := make([]byte, frameHeaderSizeBytes+len(payload))
	wire[0] = version
	wire[1] = frameType
	binary.BigEndian.PutUint32(wire[2:], length)
	copy(wire[frameHeaderSizeBytes:], payload)
	return wire
}

type chunkReader struct {
	reader io.Reader
	max    int
}

func (r *chunkReader) Read(p []byte) (int, error) {
	if len(p) > r.max {
		p = p[:r.max]
	}
	return r.reader.Read(p)
}

type chunkWriter struct {
	buffer bytes.Buffer
	max    int
}

func (w *chunkWriter) Write(p []byte) (int, error) {
	if len(p) > w.max {
		p = p[:w.max]
	}
	return w.buffer.Write(p)
}

type failingWriter struct {
	written int
	limit   int
	err     error
}

func (w *failingWriter) Write(p []byte) (int, error) {
	remaining := w.limit - w.written
	if remaining <= 0 {
		return 0, w.err
	}
	if len(p) > remaining {
		p = p[:remaining]
	}
	w.written += len(p)
	if w.written == w.limit {
		return len(p), w.err
	}
	return len(p), nil
}

func TestValidateDrain(t *testing.T) {
	t.Parallel()

	valid := Drain{
		ProtocolVersion:     ProtocolVersion,
		ConnectionID:        "0190f9d8-4d7a-7c54-8c74-4ee0049f07e0",
		ConnectorInstanceID: "0190f9d8-4d7a-7c54-8c74-4ee0049f07e1",
		SessionID:           "0190f9d8-4d7a-7c54-8c74-4ee0049f07e2",
	}
	if err := ValidateDrain(valid); err != nil {
		t.Fatalf("ValidateDrain(valid) = %v", err)
	}

	for _, mutate := range []struct {
		name  string
		apply func(*Drain)
	}{
		{name: "bad version", apply: func(d *Drain) { d.ProtocolVersion = 2 }},
		{name: "bad connection id", apply: func(d *Drain) { d.ConnectionID = "not-a-uuid" }},
		{name: "bad session id", apply: func(d *Drain) { d.SessionID = "" }},
	} {
		t.Run(mutate.name, func(t *testing.T) {
			t.Parallel()
			bad := valid
			mutate.apply(&bad)
			if err := ValidateDrain(bad); err == nil {
				t.Fatal("ValidateDrain() = nil, want error")
			}
		})
	}
}
