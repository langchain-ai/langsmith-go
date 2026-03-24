package models

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

const maxJSONMergeDepth = 5

// MergePatchToPost merges patch ops into their corresponding post ops for the same run ID.
func MergePatchToPost(ops []*SerializedOp) ([]*SerializedOp, error) {
	// Pass 1: index all posts by run ID.
	postIdxs := make(map[uuid.UUID]int)
	out := make([]*SerializedOp, 0, len(ops))
	var patches []*SerializedOp
	for _, op := range ops {
		switch op.Kind {
		case OpKindPost:
			postIdxs[op.ID] = len(out)
			out = append(out, op)
		case OpKindPatch:
			patches = append(patches, op)
		default:
			return nil, fmt.Errorf("unknown op kind %s", op.Kind)
		}
	}

	// Pass 2: merge each patch into its post, or keep as standalone.
	for _, op := range patches {
		if idx, ok := postIdxs[op.ID]; ok {
			mergeInto(out[idx], op)
		} else {
			out = append(out, op)
		}
	}
	return out, nil
}

func mergeInto(dst, src *SerializedOp) {
	if len(src.RunInfo) > 0 {
		dst.RunInfo = overlayJSON(dst.RunInfo, src.RunInfo)
	}
	if src.Inputs != nil {
		dst.Inputs = src.Inputs
	}
	if src.Outputs != nil {
		dst.Outputs = src.Outputs
	}
	if src.Events != nil {
		dst.Events = src.Events
	}
	if src.Extra != nil {
		dst.Extra = overlayJSON(dst.Extra, src.Extra)
	}
	if src.Error != nil {
		dst.Error = src.Error
	}
	if src.Serialized != nil {
		dst.Serialized = src.Serialized
	}
	if len(src.Attachments) > 0 {
		if dst.Attachments == nil {
			dst.Attachments = make(map[string]Attachment, len(src.Attachments))
		}
		for k, v := range src.Attachments {
			dst.Attachments[k] = v
		}
	}
}

func overlayJSON(base, overlay []byte) []byte {
	if len(overlay) == 0 {
		return base
	}
	if len(base) == 0 {
		return overlay
	}
	var baseMap map[string]any
	if err := json.Unmarshal(base, &baseMap); err != nil {
		return overlay
	}
	var overlayMap map[string]any
	if err := json.Unmarshal(overlay, &overlayMap); err != nil {
		return base
	}
	merged := mergeJSONMaps(baseMap, overlayMap, maxJSONMergeDepth)
	out, err := json.Marshal(merged)
	if err != nil {
		return overlay
	}
	return out
}

func mergeJSONMaps(dst, src map[string]any, maxDepth int) map[string]any {
	if maxDepth == 0 {
		return dst
	}
	if dst == nil {
		dst = make(map[string]any, len(src))
	}
	for k, v := range src {
		if v == nil {
			continue
		}
		if srcMap, ok := v.(map[string]any); ok {
			if dstMap, ok := dst[k].(map[string]any); ok {
				dst[k] = mergeJSONMaps(dstMap, srcMap, maxDepth-1)
				continue
			}
		}
		dst[k] = v
	}
	return dst
}
