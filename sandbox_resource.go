package langsmith

import (
	"context"
	"encoding/json"

	"github.com/langchain-ai/langsmith-go/option"
)

// Sandbox is a convenience wrapper around generated sandbox box responses.
type Sandbox struct {
	ID              string
	Name            string
	DataplaneURL    string
	Status          string
	StatusMessage   string
	CreatedAt       string
	UpdatedAt       string
	TTLSeconds      int64
	IdleTTLSeconds  int64
	ExpiresAt       string
	SnapshotID      string
	Vcpus           int64
	MemBytes        int64
	FsCapacityBytes int64

	boxes *SandboxBoxService
}

// NewSandbox creates a sandbox and returns the convenience wrapper.
func (r *SandboxBoxService) NewSandbox(ctx context.Context, body SandboxBoxNewParams, opts ...option.RequestOption) (*Sandbox, error) {
	res, err := r.New(ctx, body, opts...)
	if err != nil {
		return nil, err
	}
	return sandboxFromNewResponse(res, r), nil
}

// GetSandbox retrieves a sandbox and returns the convenience wrapper.
func (r *SandboxBoxService) GetSandbox(ctx context.Context, name string, opts ...option.RequestOption) (*Sandbox, error) {
	res, err := r.Get(ctx, name, opts...)
	if err != nil {
		return nil, err
	}
	return sandboxFromGetResponse(res, r), nil
}

// ListSandboxes lists sandboxes and returns convenience wrappers.
func (r *SandboxBoxService) ListSandboxes(ctx context.Context, query SandboxBoxListParams, opts ...option.RequestOption) ([]*Sandbox, error) {
	res, err := r.List(ctx, query, opts...)
	if err != nil {
		return nil, err
	}
	out := make([]*Sandbox, 0, len(res.Sandboxes))
	for i := range res.Sandboxes {
		out = append(out, sandboxFromListResponse(&res.Sandboxes[i], r))
	}
	return out, nil
}

// Refresh fetches latest sandbox state and updates this object.
func (s *Sandbox) Refresh(ctx context.Context, opts ...option.RequestOption) error {
	box, err := s.boxes.Get(ctx, s.Name, opts...)
	if err != nil {
		return err
	}
	s.applyGetResponse(box)
	return nil
}

// Update updates this sandbox and refreshes this object's fields.
func (s *Sandbox) Update(ctx context.Context, body SandboxBoxUpdateParams, opts ...option.RequestOption) error {
	box, err := s.boxes.Update(ctx, s.Name, body, opts...)
	if err != nil {
		return err
	}
	s.applyUpdateResponse(box)
	return nil
}

// Start starts this sandbox and waits until it is ready.
func (s *Sandbox) Start(ctx context.Context, params SandboxWaitParams, opts ...option.RequestOption) error {
	box, err := s.boxes.StartAndWait(ctx, s.Name, params, opts...)
	if err != nil {
		return err
	}
	s.applyGetResponse(box)
	return nil
}

// Stop stops this sandbox.
func (s *Sandbox) Stop(ctx context.Context, opts ...option.RequestOption) error {
	if err := s.boxes.Stop(ctx, s.Name, opts...); err != nil {
		return err
	}
	s.Status = "stopped"
	s.DataplaneURL = ""
	return nil
}

// Delete deletes this sandbox.
func (s *Sandbox) Delete(ctx context.Context, opts ...option.RequestOption) error {
	return s.boxes.Delete(ctx, s.Name, opts...)
}

func sandboxFromNewResponse(res *SandboxBoxNewResponse, boxes *SandboxBoxService) *Sandbox {
	if res == nil {
		return nil
	}
	return &Sandbox{
		ID:              res.ID,
		Name:            res.Name,
		DataplaneURL:    res.DataplaneURL,
		Status:          res.Status,
		StatusMessage:   res.StatusMessage,
		CreatedAt:       res.CreatedAt,
		UpdatedAt:       res.UpdatedAt,
		TTLSeconds:      sandboxTTLSeconds(res.DeleteAfterStopSeconds, res.JSON.RawJSON()),
		IdleTTLSeconds:  res.IdleTtlSeconds,
		ExpiresAt:       sandboxExpiresAt(res.JSON.RawJSON()),
		SnapshotID:      res.SnapshotID,
		Vcpus:           res.Vcpus,
		MemBytes:        res.MemBytes,
		FsCapacityBytes: res.FsCapacityBytes,
		boxes:           boxes,
	}
}

func sandboxFromGetResponse(res *SandboxBoxGetResponse, boxes *SandboxBoxService) *Sandbox {
	if res == nil {
		return nil
	}
	sb := &Sandbox{boxes: boxes}
	sb.applyGetResponse(res)
	return sb
}

func sandboxFromListResponse(res *SandboxBoxListResponseSandbox, boxes *SandboxBoxService) *Sandbox {
	if res == nil {
		return nil
	}
	return &Sandbox{
		ID:              res.ID,
		Name:            res.Name,
		DataplaneURL:    res.DataplaneURL,
		Status:          res.Status,
		StatusMessage:   res.StatusMessage,
		CreatedAt:       res.CreatedAt,
		UpdatedAt:       res.UpdatedAt,
		TTLSeconds:      sandboxTTLSeconds(res.DeleteAfterStopSeconds, res.JSON.RawJSON()),
		IdleTTLSeconds:  res.IdleTtlSeconds,
		ExpiresAt:       sandboxExpiresAt(res.JSON.RawJSON()),
		SnapshotID:      res.SnapshotID,
		Vcpus:           res.Vcpus,
		MemBytes:        res.MemBytes,
		FsCapacityBytes: res.FsCapacityBytes,
		boxes:           boxes,
	}
}

func (s *Sandbox) applyGetResponse(res *SandboxBoxGetResponse) {
	s.ID = res.ID
	s.Name = res.Name
	s.DataplaneURL = res.DataplaneURL
	s.Status = res.Status
	s.StatusMessage = res.StatusMessage
	s.CreatedAt = res.CreatedAt
	s.UpdatedAt = res.UpdatedAt
	s.TTLSeconds = sandboxTTLSeconds(res.DeleteAfterStopSeconds, res.JSON.RawJSON())
	s.IdleTTLSeconds = res.IdleTtlSeconds
	s.ExpiresAt = sandboxExpiresAt(res.JSON.RawJSON())
	s.SnapshotID = res.SnapshotID
	s.Vcpus = res.Vcpus
	s.MemBytes = res.MemBytes
	s.FsCapacityBytes = res.FsCapacityBytes
}

func (s *Sandbox) applyUpdateResponse(res *SandboxBoxUpdateResponse) {
	s.ID = res.ID
	s.Name = res.Name
	s.DataplaneURL = res.DataplaneURL
	s.Status = res.Status
	s.StatusMessage = res.StatusMessage
	s.CreatedAt = res.CreatedAt
	s.UpdatedAt = res.UpdatedAt
	s.TTLSeconds = sandboxTTLSeconds(res.DeleteAfterStopSeconds, res.JSON.RawJSON())
	s.IdleTTLSeconds = res.IdleTtlSeconds
	s.ExpiresAt = sandboxExpiresAt(res.JSON.RawJSON())
	s.SnapshotID = res.SnapshotID
	s.Vcpus = res.Vcpus
	s.MemBytes = res.MemBytes
	s.FsCapacityBytes = res.FsCapacityBytes
}

func sandboxTTLSeconds(deleteAfterStopSeconds int64, rawJSON string) int64 {
	if deleteAfterStopSeconds != 0 {
		return deleteAfterStopSeconds
	}
	return sandboxLegacyInt64Field(rawJSON, "ttl_seconds")
}

func sandboxExpiresAt(rawJSON string) string {
	return sandboxLegacyStringField(rawJSON, "expires_at")
}

func sandboxLegacyInt64Field(rawJSON string, key string) int64 {
	if rawJSON == "" {
		return 0
	}
	payload := map[string]json.RawMessage{}
	if err := json.Unmarshal([]byte(rawJSON), &payload); err != nil {
		return 0
	}
	raw, ok := payload[key]
	if !ok {
		return 0
	}
	var value int64
	if err := json.Unmarshal(raw, &value); err != nil {
		return 0
	}
	return value
}

func sandboxLegacyStringField(rawJSON string, key string) string {
	if rawJSON == "" {
		return ""
	}
	payload := map[string]json.RawMessage{}
	if err := json.Unmarshal([]byte(rawJSON), &payload); err != nil {
		return ""
	}
	raw, ok := payload[key]
	if !ok {
		return ""
	}
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return ""
	}
	return value
}
