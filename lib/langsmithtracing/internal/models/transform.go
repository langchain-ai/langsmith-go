package models

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// RunOp is a decoded run operation exposed to transform hooks.
// It combines all the split-out fields of a SerializedOp into a single
// map for easy inspection and modification.
type RunOp struct {
	Kind    string            // "post" or "patch"
	ID      uuid.UUID         //
	TraceID uuid.UUID         //
	Data    map[string]any    // merged run dict (run info + inputs + outputs + extra + events + error + serialized)
	Attachments map[string]Attachment // preserved as-is; not part of Data
}

// DeserializeOp converts a SerializedOp into a RunOp by merging all
// JSON fields into a single map.
func DeserializeOp(op *SerializedOp) (RunOp, error) {
	var data map[string]any
	if len(op.RunInfo) > 0 {
		if err := json.Unmarshal(op.RunInfo, &data); err != nil {
			return RunOp{}, fmt.Errorf("unmarshal run info for %s: %w", op.ID, err)
		}
	}
	if data == nil {
		data = make(map[string]any)
	}

	merge := func(field string, raw []byte) error {
		if len(raw) == 0 {
			return nil
		}
		var v any
		if err := json.Unmarshal(raw, &v); err != nil {
			return fmt.Errorf("unmarshal %s for %s: %w", field, op.ID, err)
		}
		data[field] = v
		return nil
	}

	for _, f := range []struct {
		name string
		data []byte
	}{
		{"inputs", op.Inputs},
		{"outputs", op.Outputs},
		{"events", op.Events},
		{"extra", op.Extra},
		{"error", op.Error},
		{"serialized", op.Serialized},
	} {
		if err := merge(f.name, f.data); err != nil {
			return RunOp{}, err
		}
	}

	return RunOp{
		Kind:        string(op.Kind),
		ID:          op.ID,
		TraceID:     op.TraceID,
		Data:        data,
		Attachments: op.Attachments,
	}, nil
}

var splitFields = map[string]bool{
	"inputs": true, "outputs": true, "events": true,
	"extra": true, "error": true, "serialized": true,
}

// SerializeOp converts a RunOp back into a SerializedOp by splitting
// known fields out of Data into separate byte slices.
func SerializeOp(r RunOp) (*SerializedOp, error) {
	op := &SerializedOp{
		Kind:        OpKind(r.Kind),
		ID:          r.ID,
		TraceID:     r.TraceID,
		Attachments: r.Attachments,
	}

	marshalField := func(key string) ([]byte, error) {
		v, ok := r.Data[key]
		if !ok || v == nil {
			return nil, nil
		}
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("marshal %s for %s: %w", key, r.ID, err)
		}
		return b, nil
	}

	var err error
	if op.Inputs, err = marshalField("inputs"); err != nil {
		return nil, err
	}
	if op.Outputs, err = marshalField("outputs"); err != nil {
		return nil, err
	}
	if op.Events, err = marshalField("events"); err != nil {
		return nil, err
	}
	if op.Extra, err = marshalField("extra"); err != nil {
		return nil, err
	}
	if op.Error, err = marshalField("error"); err != nil {
		return nil, err
	}
	if op.Serialized, err = marshalField("serialized"); err != nil {
		return nil, err
	}

	// RunInfo is everything in Data except the split-out fields.
	runInfo := make(map[string]any, len(r.Data))
	for k, v := range r.Data {
		if !splitFields[k] {
			runInfo[k] = v
		}
	}
	op.RunInfo, err = json.Marshal(runInfo)
	if err != nil {
		return nil, fmt.Errorf("marshal run info for %s: %w", r.ID, err)
	}

	return op, nil
}
