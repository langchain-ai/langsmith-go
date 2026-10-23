package traceutil

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"hash/crc32"
	"io"
	"reflect"
	"testing"
)

func stringHeader(name, value string) []byte {
	out := append([]byte{byte(len(name))}, name...)
	out = append(out, 7)
	out = binary.BigEndian.AppendUint16(out, uint16(len(value)))
	return append(out, value...)
}

func encodeEventStreamMessage(headers, payload []byte) []byte {
	total := eventStreamMinLen + len(headers) + len(payload)
	msg := binary.BigEndian.AppendUint32(nil, uint32(total))
	msg = binary.BigEndian.AppendUint32(msg, uint32(len(headers)))
	msg = binary.BigEndian.AppendUint32(msg, crc32.ChecksumIEEE(msg))
	msg = append(msg, headers...)
	msg = append(msg, payload...)
	return binary.BigEndian.AppendUint32(msg, crc32.ChecksumIEEE(msg))
}

func chunkMessage(t *testing.T, event string) []byte {
	t.Helper()
	payload, err := json.Marshal(map[string]any{"bytes": []byte(event), "p": "abc"})
	if err != nil {
		t.Fatal(err)
	}
	headers := append(stringHeader(":message-type", "event"), stringHeader(":event-type", "chunk")...)
	headers = append(headers, append([]byte{4, 's', 'e', 'q', 'n', 4}, 0, 0, 0, 1)...)
	headers = append(headers, stringHeader(":content-type", "application/json")...)
	return encodeEventStreamMessage(headers, payload)
}

func TestEventStreamScanner(t *testing.T) {
	exception := encodeEventStreamMessage(
		append(stringHeader(":message-type", "exception"), stringHeader(":exception-type", "throttlingException")...),
		[]byte(`{"message":"slow down"}`),
	)
	valid := chunkMessage(t, `{"type":"message_start"}`)
	corrupt := bytes.Clone(valid)
	corrupt[len(corrupt)-6] ^= 0xff

	for _, tc := range []struct {
		name       string
		feeds      [][]byte
		wantChunks []map[string]any
		wantErr    error
	}{
		{
			name: "chunks split at every byte",
			feeds: func() [][]byte {
				stream := append(chunkMessage(t, `{"type":"message_start"}`), chunkMessage(t, `{"type":"message_stop"}`)...)
				var feeds [][]byte
				for i := range stream {
					feeds = append(feeds, stream[i:i+1])
				}
				return feeds
			}(),
			wantChunks: []map[string]any{{"type": "message_start"}, {"type": "message_stop"}},
		},
		{
			name:  "exception becomes an error event",
			feeds: [][]byte{exception},
			wantChunks: []map[string]any{{
				"type":  "error",
				"error": map[string]any{"type": "throttlingException", "message": "slow down"},
			}},
		},
		{
			name:    "corrupt message",
			feeds:   [][]byte{corrupt, valid},
			wantErr: errMalformedEventStream,
		},
		{
			name: "invalid chunk JSON",
			feeds: [][]byte{encodeEventStreamMessage(
				append(stringHeader(":message-type", "event"), stringHeader(":event-type", "chunk")...),
				[]byte(`{"bytes":"bm90IGpzb24="}`),
			), valid},
			wantErr: errMalformedEventStream,
		},
		{
			name:    "not an event stream",
			feeds:   [][]byte{[]byte("data: {\"type\":\"message_start\"}\n\n")},
			wantErr: errMalformedEventStream,
		},
		{
			name:       "truncated final message",
			feeds:      [][]byte{valid, valid[:len(valid)-1]},
			wantChunks: []map[string]any{{"type": "message_start"}},
			wantErr:    io.ErrUnexpectedEOF,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var chunks []map[string]any
			var errs []error
			s := NewEventStreamScanner(func(c map[string]any) { chunks = append(chunks, c) })
			s.OnError = func(err error) { errs = append(errs, err) }
			for _, feed := range tc.feeds {
				s.Feed(feed)
			}
			s.Finish()
			if !reflect.DeepEqual(chunks, tc.wantChunks) {
				t.Errorf("chunks = %v, want %v", chunks, tc.wantChunks)
			}
			switch {
			case tc.wantErr == nil && len(errs) > 0:
				t.Errorf("unexpected errors %v", errs)
			case tc.wantErr != nil && (len(errs) != 1 || !errors.Is(errs[0], tc.wantErr)):
				t.Errorf("errors = %v, want exactly %v", errs, tc.wantErr)
			}
		})
	}
}
