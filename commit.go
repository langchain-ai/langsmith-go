// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"time"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/option"
)

// CommitService contains methods and other services that help with interacting
// with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewCommitService] method instead.
type CommitService struct {
	Options []option.RequestOption
}

// NewCommitService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewCommitService(opts ...option.RequestOption) (r *CommitService) {
	r = &CommitService{}
	r.Options = opts
	return
}

// Response model for get_commit_manifest.
type CommitManifestResponse struct {
	CommitHash string                          `json:"commit_hash,required"`
	Manifest   map[string]interface{}          `json:"manifest,required"`
	Examples   []CommitManifestResponseExample `json:"examples,nullable"`
	JSON       commitManifestResponseJSON      `json:"-"`
}

// commitManifestResponseJSON contains the JSON metadata for the struct
// [CommitManifestResponse]
type commitManifestResponseJSON struct {
	CommitHash  apijson.Field
	Manifest    apijson.Field
	Examples    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CommitManifestResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r commitManifestResponseJSON) RawJSON() string {
	return r.raw
}

// Response model for example runs
type CommitManifestResponseExample struct {
	ID        string                            `json:"id,required" format:"uuid"`
	SessionID string                            `json:"session_id,required" format:"uuid"`
	Inputs    map[string]interface{}            `json:"inputs,nullable"`
	Outputs   map[string]interface{}            `json:"outputs,nullable"`
	StartTime time.Time                         `json:"start_time,nullable" format:"date-time"`
	JSON      commitManifestResponseExampleJSON `json:"-"`
}

// commitManifestResponseExampleJSON contains the JSON metadata for the struct
// [CommitManifestResponseExample]
type commitManifestResponseExampleJSON struct {
	ID          apijson.Field
	SessionID   apijson.Field
	Inputs      apijson.Field
	Outputs     apijson.Field
	StartTime   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CommitManifestResponseExample) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r commitManifestResponseExampleJSON) RawJSON() string {
	return r.raw
}
