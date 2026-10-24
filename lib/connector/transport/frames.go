package transport

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

const (
	frameHeaderSizeBytes     = 6
	MaxControlFrameSizeBytes = 64 * 1024
)

var (
	ErrUnsupportedVersion  = errors.New("unsupported carrier protocol version")
	ErrUnknownFrameType    = errors.New("unknown carrier frame type")
	ErrUnexpectedFrameType = errors.New("unexpected carrier frame type")
	ErrFrameTooLarge       = errors.New("carrier control frame too large")
	ErrMalformedFrame      = errors.New("malformed carrier control frame")
)

// WriteFrame writes one length-prefixed JSON control frame. The byte count is
// meaningful on error: once any FLOW_OPEN byte is written, retry is ambiguous.
func WriteFrame(w io.Writer, frameType FrameType, value any) (int, error) {
	if !validFrameType(frameType) {
		return 0, fmt.Errorf("%w: %d", ErrUnknownFrameType, frameType)
	}

	payload, err := json.Marshal(value)
	if err != nil {
		return 0, fmt.Errorf("encode carrier frame: %w", err)
	}
	if len(payload) > MaxControlFrameSizeBytes {
		return 0, fmt.Errorf("%w: %d bytes", ErrFrameTooLarge, len(payload))
	}

	wire := make([]byte, frameHeaderSizeBytes+len(payload))
	wire[0] = ProtocolVersion
	wire[1] = byte(frameType)
	binary.BigEndian.PutUint32(wire[2:frameHeaderSizeBytes], uint32(len(payload)))
	copy(wire[frameHeaderSizeBytes:], payload)

	return writeAll(w, wire)
}

// ReadFrame reads one frame and requires the caller's expected message type.
func ReadFrame(r io.Reader, expectedType FrameType, out any) error {
	if !validFrameType(expectedType) {
		return fmt.Errorf("%w: %d", ErrUnknownFrameType, expectedType)
	}

	var header [frameHeaderSizeBytes]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return fmt.Errorf("read carrier frame header: %w", err)
	}
	if header[0] != ProtocolVersion {
		return fmt.Errorf("%w: %d", ErrUnsupportedVersion, header[0])
	}

	frameType := FrameType(header[1])
	if !validFrameType(frameType) {
		return fmt.Errorf("%w: %d", ErrUnknownFrameType, frameType)
	}
	if frameType != expectedType {
		return fmt.Errorf("%w: got %d, want %d", ErrUnexpectedFrameType, frameType, expectedType)
	}

	payloadSize := binary.BigEndian.Uint32(header[2:])
	if payloadSize > MaxControlFrameSizeBytes {
		return fmt.Errorf("%w: %d bytes", ErrFrameTooLarge, payloadSize)
	}
	payload := make([]byte, int(payloadSize))
	if _, err := io.ReadFull(r, payload); err != nil {
		return fmt.Errorf("read carrier frame payload: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(out); err != nil {
		if errors.Is(err, ErrInvalidFlowStatus) {
			return fmt.Errorf("decode carrier frame: %w", err)
		}
		return ErrMalformedFrame
	}
	if err := requireJSONEOF(decoder); err != nil {
		return err
	}
	return nil
}

func writeAll(w io.Writer, data []byte) (int, error) {
	written := 0
	for written < len(data) {
		n, err := w.Write(data[written:])
		if n < 0 || n > len(data)-written {
			return written, io.ErrShortWrite
		}
		written += n
		if err != nil {
			return written, err
		}
		if n == 0 {
			return written, io.ErrShortWrite
		}
	}
	return written, nil
}

func requireJSONEOF(decoder *json.Decoder) error {
	var trailing any
	err := decoder.Decode(&trailing)
	if errors.Is(err, io.EOF) {
		return nil
	}
	if err != nil {
		return ErrMalformedFrame
	}
	return fmt.Errorf("%w: multiple JSON values", ErrMalformedFrame)
}
