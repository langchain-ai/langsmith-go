// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"net/http"
	"slices"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// SqlService contains methods and other services that help with interacting with
// the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSqlService] method instead.
type SqlService struct {
	Options []option.RequestOption
}

// NewSqlService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewSqlService(opts ...option.RequestOption) (r *SqlService) {
	r = &SqlService{}
	r.Options = opts
	return
}

// **Beta:** This endpoint is in active development and may change without notice.
//
// GET with body payload — no resources created. Returns the rows matching the
// read-only SQL query passed as the request payload, limited to the given tracing
// projects. Results are limited to a configured number of rows, 1000 by default.
func (r *SqlService) Query(ctx context.Context, body SqlQueryParams, opts ...option.RequestOption) (res *SqlQueryResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v2/sql/query"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type SqlQueryResponse struct {
	Columns         []SqlQueryResponseColumn `json:"columns"`
	ExecutionTimeMs float64                  `json:"execution_time_ms"`
	// IsTruncated reports that the query matched rows not included in Rows.
	IsTruncated bool  `json:"is_truncated"`
	RowCount    int64 `json:"row_count"`
	// Rows holds at most the configured row limit (1000 by default), each ordered like
	// Columns.
	Rows [][]interface{} `json:"rows"`
	// TruncationReason is "row_limit" when the query matched more rows than the row
	// limit, "size_limit" when the rows past about 10 MB were left out, or null.
	TruncationReason SqlQueryResponseTruncationReason `json:"truncation_reason" api:"nullable"`
	// Warnings contains query planning warnings. Empty when there are no warnings.
	Warnings []string             `json:"warnings"`
	JSON     sqlQueryResponseJSON `json:"-"`
}

// sqlQueryResponseJSON contains the JSON metadata for the struct
// [SqlQueryResponse]
type sqlQueryResponseJSON struct {
	Columns          apijson.Field
	ExecutionTimeMs  apijson.Field
	IsTruncated      apijson.Field
	RowCount         apijson.Field
	Rows             apijson.Field
	TruncationReason apijson.Field
	Warnings         apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SqlQueryResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sqlQueryResponseJSON) RawJSON() string {
	return r.raw
}

type SqlQueryResponseColumn struct {
	// DataType is how the column's values are encoded in JSON: uuid, string, integer,
	// float, boolean, timestamp and date (RFC 3339 strings), binary (base64 strings),
	// list, object or null. Float values JSON cannot represent are the strings "NaN",
	// "Infinity" and "-Infinity".
	DataType SqlQueryResponseColumnsDataType `json:"data_type"`
	Name     string                          `json:"name"`
	JSON     sqlQueryResponseColumnJSON      `json:"-"`
}

// sqlQueryResponseColumnJSON contains the JSON metadata for the struct
// [SqlQueryResponseColumn]
type sqlQueryResponseColumnJSON struct {
	DataType    apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SqlQueryResponseColumn) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sqlQueryResponseColumnJSON) RawJSON() string {
	return r.raw
}

// DataType is how the column's values are encoded in JSON: uuid, string, integer,
// float, boolean, timestamp and date (RFC 3339 strings), binary (base64 strings),
// list, object or null. Float values JSON cannot represent are the strings "NaN",
// "Infinity" and "-Infinity".
type SqlQueryResponseColumnsDataType string

const (
	SqlQueryResponseColumnsDataTypeUuid      SqlQueryResponseColumnsDataType = "uuid"
	SqlQueryResponseColumnsDataTypeString    SqlQueryResponseColumnsDataType = "string"
	SqlQueryResponseColumnsDataTypeInteger   SqlQueryResponseColumnsDataType = "integer"
	SqlQueryResponseColumnsDataTypeFloat     SqlQueryResponseColumnsDataType = "float"
	SqlQueryResponseColumnsDataTypeBoolean   SqlQueryResponseColumnsDataType = "boolean"
	SqlQueryResponseColumnsDataTypeTimestamp SqlQueryResponseColumnsDataType = "timestamp"
	SqlQueryResponseColumnsDataTypeDate      SqlQueryResponseColumnsDataType = "date"
	SqlQueryResponseColumnsDataTypeBinary    SqlQueryResponseColumnsDataType = "binary"
	SqlQueryResponseColumnsDataTypeList      SqlQueryResponseColumnsDataType = "list"
	SqlQueryResponseColumnsDataTypeObject    SqlQueryResponseColumnsDataType = "object"
	SqlQueryResponseColumnsDataTypeNull      SqlQueryResponseColumnsDataType = "null"
)

func (r SqlQueryResponseColumnsDataType) IsKnown() bool {
	switch r {
	case SqlQueryResponseColumnsDataTypeUuid, SqlQueryResponseColumnsDataTypeString, SqlQueryResponseColumnsDataTypeInteger, SqlQueryResponseColumnsDataTypeFloat, SqlQueryResponseColumnsDataTypeBoolean, SqlQueryResponseColumnsDataTypeTimestamp, SqlQueryResponseColumnsDataTypeDate, SqlQueryResponseColumnsDataTypeBinary, SqlQueryResponseColumnsDataTypeList, SqlQueryResponseColumnsDataTypeObject, SqlQueryResponseColumnsDataTypeNull:
		return true
	}
	return false
}

// TruncationReason is "row_limit" when the query matched more rows than the row
// limit, "size_limit" when the rows past about 10 MB were left out, or null.
type SqlQueryResponseTruncationReason string

const (
	SqlQueryResponseTruncationReasonRowLimit  SqlQueryResponseTruncationReason = "row_limit"
	SqlQueryResponseTruncationReasonSizeLimit SqlQueryResponseTruncationReason = "size_limit"
)

func (r SqlQueryResponseTruncationReason) IsKnown() bool {
	switch r {
	case SqlQueryResponseTruncationReasonRowLimit, SqlQueryResponseTruncationReasonSizeLimit:
		return true
	}
	return false
}

type SqlQueryParams struct {
	// ProjectIDs are the tracing projects the query can read.
	ProjectIDs param.Field[[]string] `json:"project_ids" api:"required" format:"uuid"`
	// Query is a single read-only SQL statement. Queries over runs need a start_time
	// lower bound.
	Query param.Field[string] `json:"query" api:"required"`
}

func (r SqlQueryParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
