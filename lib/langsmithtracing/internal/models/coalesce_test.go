package models

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestCoalesce_MergesAttachments(t *testing.T) {
	id := uuid.New()
	traceID := uuid.New()

	ops := []*SerializedOp{
		{
			Kind:    OpKindPost,
			ID:      id,
			TraceID: traceID,
			RunInfo: []byte(`{"name":"test"}`),
			Attachments: map[string]Attachment{
				"image": {ContentType: "image/png", Data: []byte("png-data")},
			},
		},
		{
			Kind:    OpKindPatch,
			ID:      id,
			TraceID: traceID,
			RunInfo: []byte(`{"end_time":"2024-01-01T00:00:00Z"}`),
			Attachments: map[string]Attachment{
				"report": {ContentType: "application/pdf", Data: []byte("pdf-data")},
			},
		},
	}

	result, err := Coalesce(ops)
	if err != nil {
		t.Fatalf("Coalesce: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 op, got %d", len(result))
	}

	merged := result[0]
	if merged.Kind != OpKindPost {
		t.Errorf("expected post, got %s", merged.Kind)
	}
	if len(merged.Attachments) != 2 {
		t.Fatalf("expected 2 attachments, got %d", len(merged.Attachments))
	}
	if string(merged.Attachments["image"].Data) != "png-data" {
		t.Error("image attachment data mismatch")
	}
	if string(merged.Attachments["report"].Data) != "pdf-data" {
		t.Error("report attachment data mismatch")
	}
}

func TestCoalesce_PatchOverridesAttachment(t *testing.T) {
	id := uuid.New()

	ops := []*SerializedOp{
		{
			Kind:    OpKindPost,
			ID:      id,
			TraceID: id,
			RunInfo: []byte(`{}`),
			Attachments: map[string]Attachment{
				"file": {ContentType: "text/plain", Data: []byte("v1")},
			},
		},
		{
			Kind:    OpKindPatch,
			ID:      id,
			TraceID: id,
			RunInfo: []byte(`{}`),
			Attachments: map[string]Attachment{
				"file": {ContentType: "text/plain", Data: []byte("v2")},
			},
		},
	}

	result, err := Coalesce(ops)
	if err != nil {
		t.Fatalf("Coalesce: %v", err)
	}
	if string(result[0].Attachments["file"].Data) != "v2" {
		t.Error("patch should override post attachment with same key")
	}
}

func TestCoalesce_ExtraDeepMerged(t *testing.T) {
	id := uuid.New()

	ops := []*SerializedOp{
		{
			Kind:    OpKindPost,
			ID:      id,
			TraceID: id,
			RunInfo: []byte(`{}`),
			Extra:   []byte(`{"runtime":{"sdk":"langsmith-go","platform":"linux/amd64"},"metadata":{"user_key":"keep"}}`),
		},
		{
			Kind:    OpKindPatch,
			ID:      id,
			TraceID: id,
			RunInfo: []byte(`{}`),
			Extra:   []byte(`{"metadata":{"v":2}}`),
		},
	}

	result, err := Coalesce(ops)
	if err != nil {
		t.Fatalf("Coalesce: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 op, got %d", len(result))
	}

	var extra map[string]any
	if err := json.Unmarshal(result[0].Extra, &extra); err != nil {
		t.Fatalf("unmarshal extra: %v", err)
	}

	runtime, ok := extra["runtime"].(map[string]any)
	if !ok {
		t.Fatal("runtime key missing after coalesce — patch clobbered create's runtime env")
	}
	if runtime["sdk"] != "langsmith-go" {
		t.Errorf("runtime.sdk = %v, want langsmith-go", runtime["sdk"])
	}
	if runtime["platform"] != "linux/amd64" {
		t.Errorf("runtime.platform = %v, want linux/amd64", runtime["platform"])
	}

	metadata, ok := extra["metadata"].(map[string]any)
	if !ok {
		t.Fatal("metadata key missing after coalesce")
	}
	if metadata["v"] != float64(2) {
		t.Errorf("metadata.v = %v, want 2 (patch value)", metadata["v"])
	}
	if metadata["user_key"] != "keep" {
		t.Errorf("metadata.user_key = %v, want keep (create value)", metadata["user_key"])
	}
}

func TestCoalesce_StandalonePatch(t *testing.T) {
	id := uuid.New()

	ops := []*SerializedOp{
		{
			Kind:    OpKindPatch,
			ID:      id,
			TraceID: id,
			RunInfo: []byte(`{}`),
			Attachments: map[string]Attachment{
				"doc": {ContentType: "application/pdf", Data: []byte("data")},
			},
		},
	}

	result, err := Coalesce(ops)
	if err != nil {
		t.Fatalf("Coalesce: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 op, got %d", len(result))
	}
	if len(result[0].Attachments) != 1 {
		t.Error("standalone patch should preserve its attachments")
	}
}

func TestCoalesce_UnknownOpKind(t *testing.T) {
	ops := []*SerializedOp{
		{
			Kind:    OpKind("unknown"),
			ID:      uuid.New(),
			TraceID: uuid.New(),
			RunInfo: []byte(`{}`),
		},
	}
	_, err := Coalesce(ops)
	if err == nil {
		t.Fatal("expected error for unknown op kind, got nil")
	}
	if !strings.Contains(err.Error(), "unknown op kind") {
		t.Errorf("expected error containing %q, got %q", "unknown op kind", err.Error())
	}
}

func TestOverlayJSON_MalformedBase(t *testing.T) {
	base := []byte("not json")
	overlay := []byte(`{"a":1}`)
	got := overlayJSON(base, overlay)
	if !bytes.Equal(got, overlay) {
		t.Errorf("expected overlay returned unchanged, got %s", got)
	}
}

func TestOverlayJSON_MalformedOverlay(t *testing.T) {
	base := []byte(`{"a":1}`)
	overlay := []byte("not json")
	got := overlayJSON(base, overlay)
	if !bytes.Equal(got, base) {
		t.Errorf("expected base returned unchanged, got %s", got)
	}
}

func TestOverlayJSON_EmptyBase(t *testing.T) {
	overlay := []byte(`{"a":1}`)
	for _, base := range [][]byte{nil, {}} {
		got := overlayJSON(base, overlay)
		if !bytes.Equal(got, overlay) {
			t.Errorf("base=%v: expected overlay returned, got %s", base, got)
		}
	}
}

func TestOverlayJSON_EmptyOverlay(t *testing.T) {
	base := []byte(`{"a":1}`)
	for _, overlay := range [][]byte{nil, {}} {
		got := overlayJSON(base, overlay)
		if !bytes.Equal(got, base) {
			t.Errorf("overlay=%v: expected base returned, got %s", overlay, got)
		}
	}
}

func TestMergeJSONMaps_NilValuesSkipped(t *testing.T) {
	dst := map[string]any{"y": 2}
	src := map[string]any{"x": nil, "z": 3}
	result := mergeJSONMaps(dst, src, 5)
	if _, exists := result["x"]; exists {
		t.Error("nil value from src should not be set in dst")
	}
	if result["z"] != 3 {
		t.Error("non-nil value from src should be merged")
	}
	if result["y"] != 2 {
		t.Error("existing dst value should be preserved")
	}
}

func TestMergeJSONMaps_DepthLimit(t *testing.T) {
	dst := map[string]any{"a": map[string]any{"b": 1}}
	src := map[string]any{"a": map[string]any{"c": 2}}
	result := mergeJSONMaps(dst, src, 0)
	inner, ok := result["a"].(map[string]any)
	if !ok {
		t.Fatal("expected nested map for key 'a'")
	}
	if _, exists := inner["c"]; exists {
		t.Error("maxDepth=0 should prevent merging nested maps from src")
	}
	if inner["b"] != 1 {
		t.Error("original dst nested value should be preserved")
	}
}

func TestMergeJSONMaps_NestedMerge(t *testing.T) {
	dst := map[string]any{"a": map[string]any{"b": float64(1)}}
	src := map[string]any{"a": map[string]any{"c": float64(2)}}
	result := mergeJSONMaps(dst, src, 2)

	out, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}

	inner, ok := got["a"].(map[string]any)
	if !ok {
		t.Fatal("expected nested map for key 'a'")
	}
	if inner["b"] != float64(1) {
		t.Errorf("expected b=1, got %v", inner["b"])
	}
	if inner["c"] != float64(2) {
		t.Errorf("expected c=2, got %v", inner["c"])
	}
}

func TestSerializedOp_SizeNilReceiver(t *testing.T) {
	var op *SerializedOp
	if op.Size() != 0 {
		t.Errorf("nil SerializedOp.Size() = %d, want 0", op.Size())
	}
}

func TestSerializedOp_SizeIncludesAllFields(t *testing.T) {
	op := &SerializedOp{
		RunInfo:  []byte(`{"name":"x"}`),
		Inputs:   []byte(`{"a":1}`),
		Outputs:  []byte(`{"b":2}`),
		Events:   []byte(`[1]`),
		Extra:    []byte(`{}`),
		Error:    []byte(`"err"`),
		Serialized: []byte(`{}`),
		Attachments: map[string]Attachment{
			"f": {Data: []byte("12345")},
		},
	}
	expected := len(op.RunInfo) + len(op.Inputs) + len(op.Outputs) +
		len(op.Events) + len(op.Extra) + len(op.Error) + len(op.Serialized) + 5
	if got := op.Size(); got != expected {
		t.Errorf("Size() = %d, want %d", got, expected)
	}
}
