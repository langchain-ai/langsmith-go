// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package shared

import (
	"github.com/langchain-ai/langsmith-go/internal/apijson"
)

type QueryPatchEvent []QueryPatchOperation

type QueryPatchOperation struct {
	Op    QueryPatchOperationOp   `json:"op" api:"required"`
	Path  string                  `json:"path" api:"required"`
	Value interface{}             `json:"value"`
	JSON  queryPatchOperationJSON `json:"-"`
}

// queryPatchOperationJSON contains the JSON metadata for the struct
// [QueryPatchOperation]
type queryPatchOperationJSON struct {
	Op          apijson.Field
	Path        apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *QueryPatchOperation) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r queryPatchOperationJSON) RawJSON() string {
	return r.raw
}

type QueryPatchOperationOp string

const (
	QueryPatchOperationOpAdd     QueryPatchOperationOp = "add"
	QueryPatchOperationOpReplace QueryPatchOperationOp = "replace"
)

func (r QueryPatchOperationOp) IsKnown() bool {
	switch r {
	case QueryPatchOperationOpAdd, QueryPatchOperationOpReplace:
		return true
	}
	return false
}
