package cmd

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	langsmith "github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/cli/internal/client"
	"github.com/langchain-ai/langsmith-go/cli/internal/output"
	"github.com/langchain-ai/langsmith-go/option"
	"github.com/langchain-ai/langsmith-go/shared"
)

// exampleBulkBatchSize caps each POST /examples/bulk request so a long stdin
// list doesn't become one oversized body.
const exampleBulkBatchSize = 100

// collectRunIDs expands --from-run values: "-" reads run IDs from stdin, one per
// line, ignoring blank lines. Every ID must be a UUID.
func collectRunIDs(values []string, stdin io.Reader) ([]string, error) {
	var ids []string
	readStdin := false
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v != "-" {
			ids = append(ids, v)
			continue
		}
		if readStdin {
			continue
		}
		readStdin = true
		scanner := bufio.NewScanner(stdin)
		for scanner.Scan() {
			if line := strings.TrimSpace(scanner.Text()); line != "" {
				ids = append(ids, line)
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("reading run IDs from stdin: %w", err)
		}
	}
	for _, id := range ids {
		if _, err := uuid.Parse(id); err != nil {
			return nil, fmt.Errorf("invalid run ID %q for --from-run: must be a UUID", id)
		}
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("--from-run got no run IDs")
	}
	return ids, nil
}

// canonicalInputs returns a stable string for a raw inputs object: encoding/json
// sorts map keys, numbers keep their exact digits, and null and empty inputs
// compare equal.
func canonicalInputs(raw json.RawMessage) (string, error) {
	var v any
	if len(raw) > 0 {
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.UseNumber()
		if err := dec.Decode(&v); err != nil {
			return "", err
		}
	}
	if v == nil {
		v = map[string]any{}
	}
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// rawField returns a response field's JSON as the server sent it, or nil when
// it is missing or null. The SDK decodes arbitrary JSON numbers through
// float64, so copied data and timestamps are taken from the raw bytes.
func rawField(raw string) json.RawMessage {
	if raw == "" || raw == "null" {
		return nil
	}
	return json.RawMessage(raw)
}

// existingInputKeys returns the canonical inputs of every example in the dataset.
func existingInputKeys(ctx context.Context, c *client.Client, datasetID string) (map[string]bool, error) {
	keys := map[string]bool{}
	pager := c.SDK.Examples.ListAutoPaging(ctx, langsmith.ExampleListParams{
		Dataset: langsmith.F(datasetID),
		Limit:   langsmith.F(int64(100)),
		Select:  langsmith.F([]langsmith.ExampleSelect{langsmith.ExampleSelectInputs}),
	})
	for pager.Next() {
		key, err := canonicalInputs(rawField(pager.Current().JSON.Inputs.Raw()))
		if err != nil {
			return nil, err
		}
		keys[key] = true
	}
	if err := pager.Err(); err != nil {
		return nil, fmt.Errorf("listing existing examples: %w", err)
	}
	return keys, nil
}

// sourceRun is the part of a run copied into an example, kept as the raw JSON
// the server returned.
type sourceRun struct {
	id, sessionID, traceID string
	startTime              json.RawMessage
	inputs, outputs        json.RawMessage
	attachments            []string // attachment names, without the "attachment." prefix
}

// attachmentNames returns the run's attachment names from a map keyed by
// attachment name, with or without the "attachment." prefix, in sorted order.
func attachmentNames[V any](m map[string]V, onlyPrefixed bool) []string {
	var names []string
	for k := range m {
		name, prefixed := strings.CutPrefix(k, "attachment.")
		if onlyPrefixed && !prefixed {
			continue
		}
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

// fetchSourceRun reads a run through the v2 API when the deployment supports
// it and a project is given (v2 reads are scoped to a project), else through
// the v1 single-run read.
func fetchSourceRun(ctx context.Context, c *client.Client, runID, sessionID string) (sourceRun, error) {
	if sessionID != "" {
		useV2, err := c.UseV2API(ctx)
		if err != nil {
			return sourceRun{}, err
		}
		if useV2 {
			run, err := c.SDK.Runs.GetV2(ctx, runID, langsmith.RunGetV2Params{
				ProjectID: langsmith.F(sessionID),
				Selects: langsmith.F([]langsmith.RunGetV2ParamsSelect{
					langsmith.RunGetV2ParamsSelectID,
					langsmith.RunGetV2ParamsSelectTraceID,
					langsmith.RunGetV2ParamsSelectProjectID,
					langsmith.RunGetV2ParamsSelectStartTime,
					langsmith.RunGetV2ParamsSelectInputs,
					langsmith.RunGetV2ParamsSelectOutputs,
					langsmith.RunGetV2ParamsSelectAttachments,
				}),
			})
			if err != nil {
				return sourceRun{}, err
			}
			return sourceRun{
				id: run.ID, sessionID: run.ProjectID, traceID: run.TraceID,
				startTime:   rawField(run.JSON.StartTime.Raw()),
				inputs:      rawField(run.JSON.Inputs.Raw()),
				outputs:     rawField(run.JSON.Outputs.Raw()),
				attachments: attachmentNames(run.Attachments, false),
			}, nil
		}
	}
	run, err := c.SDK.Runs.Get(ctx, runID, langsmith.RunGetParams{})
	if err != nil {
		if sessionID == "" && isHTTPStatus(err, http.StatusNotImplemented) {
			return sourceRun{}, fmt.Errorf("%w\nhint: this deployment reads runs by project; pass --project or --project-id", err)
		}
		return sourceRun{}, err
	}
	return sourceRun{
		id: run.ID, sessionID: run.SessionID, traceID: run.TraceID,
		startTime:   rawField(run.JSON.StartTime.Raw()),
		inputs:      rawField(run.JSON.Inputs.Raw()),
		outputs:     rawField(run.JSON.Outputs.Raw()),
		attachments: attachmentNames(run.S3URLs, true),
	}, nil
}

// needsConversion reports whether a run's inputs are a traced model call that
// the server must convert to a chat or LLM dataset's format. The server reads
// chat inputs from "messages" and LLM inputs from "prompt" or "prompts"; runs
// whose inputs are already in the dataset's format are copied as they are,
// since converting them would drop the inputs.
func needsConversion(ds *langsmith.Dataset, inputs json.RawMessage) bool {
	var keys []string
	switch ds.DataType {
	case langsmith.DataTypeChat:
		keys = []string{"messages"}
	case langsmith.DataTypeLlm:
		keys = []string{"prompt", "prompts"}
	default:
		return false
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(inputs, &fields) != nil {
		return false
	}
	for _, k := range keys {
		if _, ok := fields[k]; ok {
			return true
		}
	}
	return false
}

// exampleFromRun builds the create request for one run. When convert is set
// the server converts the run's inputs and outputs to the dataset's format
// (use_source_run_io), which also copies the run's attachments; otherwise they
// are copied byte for byte and the attachments are named explicitly, since the
// server copies only the attachments a request names.
func exampleFromRun(datasetID string, run sourceRun, convert bool, metadata map[string]any, split string) (langsmith.ExampleBulkNewParamsBody, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return langsmith.ExampleBulkNewParamsBody{}, err
	}
	body := langsmith.ExampleBulkNewParamsBody{
		ID:          langsmith.F(id.String()),
		DatasetID:   langsmith.F(datasetID),
		SourceRunID: langsmith.F(run.id),
	}
	if convert {
		body.UseSourceRunIo = langsmith.F(true)
	} else {
		inputs := run.inputs
		if inputs == nil {
			inputs = json.RawMessage("{}")
		}
		body.Inputs = langsmith.Raw[map[string]any](inputs)
		if run.outputs != nil {
			body.Outputs = langsmith.Raw[map[string]any](run.outputs) //nolint:staticcheck // older servers reject reference_outputs
		}
		if len(run.attachments) > 0 {
			body.UseSourceRunAttachments = langsmith.F(run.attachments)
		}
	}
	// The session and start time let the server find the source run without a
	// scan. The start time is sent as the server returned it: the server
	// compares it exactly, and the SDK's date-time encoding drops fractional
	// seconds.
	if run.sessionID != "" {
		body.SourceSessionID = langsmith.F(run.sessionID)
	}
	if run.traceID != "" {
		body.SourceTraceID = langsmith.F(run.traceID)
	}
	if run.startTime != nil {
		body.SourceRunStartTime = langsmith.Raw[time.Time](run.startTime)
	}
	if metadata != nil {
		body.Metadata = langsmith.F(metadata)
	}
	if split != "" {
		body.Split = langsmith.F[langsmith.ExampleBulkNewParamsBodySplitUnion](shared.UnionString(split))
	}
	return body, nil
}

// singleExampleParams converts a bulk body to the single-create request, used
// when only one example is created.
func singleExampleParams(b langsmith.ExampleBulkNewParamsBody) langsmith.ExampleNewParams {
	params := langsmith.ExampleNewParams{
		ID:                      b.ID,
		DatasetID:               b.DatasetID,
		Inputs:                  b.Inputs,
		Outputs:                 b.Outputs, //nolint:staticcheck // older servers reject reference_outputs
		Metadata:                b.Metadata,
		SourceRunID:             b.SourceRunID,
		SourceSessionID:         b.SourceSessionID,
		SourceTraceID:           b.SourceTraceID,
		SourceRunStartTime:      b.SourceRunStartTime,
		UseSourceRunIo:          b.UseSourceRunIo,
		UseSourceRunAttachments: b.UseSourceRunAttachments,
	}
	if s, ok := b.Split.Value.(shared.UnionString); ok && b.Split.Present {
		params.Split = langsmith.F[langsmith.ExampleNewParamsSplitUnion](s)
	}
	return params
}

// convertedInputs asks the server to convert each body's source run to the
// dataset's format, without creating anything, and returns the inputs it would
// store, in order.
func convertedInputs(ctx context.Context, c *client.Client, bodies []langsmith.ExampleBulkNewParamsBody) ([]json.RawMessage, error) {
	var out []json.RawMessage
	for start := 0; start < len(bodies); start += exampleBulkBatchSize {
		end := min(start+exampleBulkBatchSize, len(bodies))
		req, err := json.Marshal(bodies[start:end])
		if err != nil {
			return nil, err
		}
		var results []struct {
			Inputs json.RawMessage `json:"inputs"`
		}
		if _, err := c.SDK.Examples.Validate.Bulk(ctx,
			option.WithRequestBody("application/json", req),
			option.WithResponseBodyInto(&results),
		); err != nil {
			return nil, fmt.Errorf("converting runs to the dataset format: %w", err)
		}
		if len(results) != end-start {
			return nil, fmt.Errorf("converting runs to the dataset format: got %d results for %d runs", len(results), end-start)
		}
		for _, r := range results {
			out = append(out, r.Inputs)
		}
	}
	return out, nil
}

// sourceRunFetchConcurrency bounds the run reads in flight at once, so a long
// stdin list is read in parallel without flooding the API.
const sourceRunFetchConcurrency = 8

// fetchSourceRuns reads the runs with a bounded pool and returns them in the
// order of runIDs. After a read fails no new reads start, and the error
// returned is that of the earliest failing run ID, so it doesn't depend on
// scheduling.
func fetchSourceRuns(ctx context.Context, c *client.Client, runIDs []string, sessionID string) ([]sourceRun, error) {
	if sessionID != "" {
		// Resolve the deployment version once here: the client caches it
		// without locking, so the concurrent reads below only read the cache.
		if _, err := c.UseV2API(ctx); err != nil {
			return nil, err
		}
	}
	runs := make([]sourceRun, len(runIDs))
	errs := make([]error, len(runIDs))
	var (
		wg     sync.WaitGroup
		failed atomic.Bool
	)
	sem := make(chan struct{}, sourceRunFetchConcurrency)
	for i, runID := range runIDs {
		sem <- struct{}{}
		if failed.Load() {
			<-sem
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			run, err := fetchSourceRun(ctx, c, runID, sessionID)
			if err != nil {
				errs[i] = fmt.Errorf("fetching run %s: %w", runID, err)
				failed.Store(true)
				return
			}
			runs[i] = run
		}()
	}
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	return runs, nil
}

// createExamplesFromRuns fetches each run, then creates one example per run in
// the dataset. All runs are fetched before anything is created, so a missing
// run leaves the dataset unchanged. sessionID, when set, is the project the
// runs belong to.
func createExamplesFromRuns(ctx context.Context, c *client.Client, datasetName string, runIDs []string, sessionID string, metadata map[string]any, split string, skipDuplicates bool) error {
	ds, err := resolveDataset(ctx, c, datasetName)
	if err != nil {
		return err
	}
	runs, err := fetchSourceRuns(ctx, c, runIDs, sessionID)
	if err != nil {
		return err
	}
	var (
		bodies []langsmith.ExampleBulkNewParamsBody
		keys   []json.RawMessage
	)
	for _, run := range runs {
		body, err := exampleFromRun(ds.ID, run, needsConversion(ds, run.inputs), metadata, split)
		if err != nil {
			return err
		}
		bodies = append(bodies, body)
		keys = append(keys, run.inputs)
	}

	createdRunIDs := runIDs
	skippedRunIDs := []string{}
	if skipDuplicates {
		seen, err := existingInputKeys(ctx, c, ds.ID)
		if err != nil {
			return err
		}
		// Converted runs are stored in the dataset's format, so compare what
		// the server would store rather than the run's own inputs.
		var toConvert []langsmith.ExampleBulkNewParamsBody
		var at []int
		for i, b := range bodies {
			if b.UseSourceRunIo.Value {
				toConvert = append(toConvert, b)
				at = append(at, i)
			}
		}
		if len(toConvert) > 0 {
			converted, err := convertedInputs(ctx, c, toConvert)
			if err != nil {
				return err
			}
			for j, i := range at {
				keys[i] = converted[j]
			}
		}
		kept := bodies[:0]
		createdRunIDs = nil
		for i, runID := range runIDs {
			key, err := canonicalInputs(keys[i])
			if err != nil {
				return fmt.Errorf("encoding inputs of run %s: %w", runID, err)
			}
			if seen[key] {
				skippedRunIDs = append(skippedRunIDs, runID)
				continue
			}
			seen[key] = true
			kept = append(kept, bodies[i])
			createdRunIDs = append(createdRunIDs, runID)
		}
		bodies = kept
	}
	exampleIDs := []string{}
	switch {
	case len(bodies) == 1:
		ex, err := c.SDK.Examples.New(ctx, singleExampleParams(bodies[0]))
		if err != nil {
			return fmt.Errorf("creating example: %w", err)
		}
		exampleIDs = append(exampleIDs, ex.ID)
	case len(bodies) > 1:
		for start := 0; start < len(bodies); start += exampleBulkBatchSize {
			end := min(start+exampleBulkBatchSize, len(bodies))
			if _, err := c.SDK.Examples.Bulk.New(ctx, langsmith.ExampleBulkNewParams{Body: bodies[start:end]}); err != nil {
				return fmt.Errorf("creating examples (%d of %d created before the error): %w", start, len(bodies), err)
			}
			for _, b := range bodies[start:end] {
				exampleIDs = append(exampleIDs, b.ID.Value)
			}
		}
	}

	if GetFormat() == "pretty" {
		var rows [][]string
		for i, runID := range createdRunIDs {
			rows = append(rows, []string{runID, "created", exampleIDs[i]})
		}
		for _, runID := range skippedRunIDs {
			rows = append(rows, []string{runID, "skipped (duplicate inputs)", ""})
		}
		output.OutputTable([]string{"Run ID", "Result", "Example ID"}, rows,
			fmt.Sprintf("Created %d example(s) in %s, skipped %d", len(exampleIDs), ds.Name, len(skippedRunIDs)))
		return nil
	}
	return output.OutputJSON(map[string]any{
		"status":          "created",
		"dataset_id":      ds.ID,
		"created":         len(exampleIDs),
		"skipped":         len(skippedRunIDs),
		"example_ids":     exampleIDs,
		"skipped_run_ids": skippedRunIDs,
	}, "")
}
