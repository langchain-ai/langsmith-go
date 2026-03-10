package models

import "github.com/google/uuid"

// OpKind indicates whether a serialized operation is a create (post) or update (patch).
type OpKind string

const (
	OpKindPost  OpKind = "post"
	OpKindPatch OpKind = "patch"
)

// Attachment is a binary file associated with a run.
// It is sent as a separate multipart part with name "attachment.<run_id>.<name>".
type Attachment struct {
	// ContentType is the MIME type (e.g. "image/png", "application/pdf").
	ContentType string
	// Data is the raw file content.
	Data []byte
}

// SerializedOp is a run create/update event ready for the multipart exporter.
// Each field (except Kind/ID/TraceID) is a pre-serialized JSON byte slice
// that becomes a separate multipart part.
type SerializedOp struct {
	Kind    OpKind
	ID      uuid.UUID
	TraceID uuid.UUID

	RunInfo    []byte // Main run JSON (everything except the split-out fields below).
	Inputs     []byte // JSON-serialized inputs.
	Outputs    []byte // JSON-serialized outputs.
	Events     []byte // JSON-serialized events (e.g. first token timing).
	Extra      []byte // JSON-serialized extra (metadata, runtime).
	Error      []byte // JSON-serialized error.
	Serialized []byte // JSON-serialized model manifest (kept only for "llm"/"prompt" run types).

	Attachments map[string]Attachment // Binary attachments keyed by name.
}

// Size returns the total byte size of all payload fields.
func (o *SerializedOp) Size() int {
	if o == nil {
		return 0
	}
	n := len(o.RunInfo) + len(o.Inputs) + len(o.Outputs) +
		len(o.Events) + len(o.Extra) + len(o.Error) + len(o.Serialized)
	for _, a := range o.Attachments {
		n += len(a.Data)
	}
	return n
}
