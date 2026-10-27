package langsmith_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
)

// sseServer answers /api/v2/{runs,threads}/query with the page keyed by the
// request's cursor and records every request body.
func sseServer(t *testing.T, pages map[string][]string) (*httptest.Server, *[]map[string]any) {
	t.Helper()
	var mu sync.Mutex
	var bodies []map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "text/event-stream" {
			t.Errorf("Accept = %q, want text/event-stream", r.Header.Get("Accept"))
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		mu.Lock()
		bodies = append(bodies, body)
		mu.Unlock()
		cursor, _ := body["cursor"].(string)
		events, ok := pages[cursor]
		if !ok {
			t.Errorf("no page for cursor %q", cursor)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		for _, e := range events {
			fmt.Fprint(w, e)
		}
	}))
	t.Cleanup(srv.Close)
	return srv, &bodies
}

func data(patches string) string { return "event: data\ndata: " + patches + "\n\n" }

const end = "event: end\ndata: \n\n"

func newSSEClient(srv *httptest.Server) *langsmith.Client {
	return langsmith.NewClient(
		option.WithBaseURL(srv.URL),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)
}

func TestQueryStreamingAutoPagingFollowsCursor(t *testing.T) {
	srv, bodies := sseServer(t, map[string][]string{
		"": {
			data(`[{"op":"add","path":"","value":{"items":[{"id":"run-1","name":"a"}],"next_cursor":null}}]`),
			data(`[{"op":"add","path":"/items/-","value":{"id":"run-2","name":"b"}},{"op":"add","path":"/items/0/extra","value":1}]`),
			data(`[{"op":"add","path":"/items/-","value":{"id":"run-3","name":"c"}},{"op":"replace","path":"/next_cursor","value":"abc"}]`),
			end,
		},
		"abc": {
			data(`[{"op":"add","path":"","value":{"items":[{"id":"run-4","name":"d"}],"next_cursor":null}}]`),
			data(`[{"op":"replace","path":"/next_cursor","value":null}]`),
			end,
		},
	})

	iter := newSSEClient(srv).Runs.Streaming.QueryV2StreamingAutoPaging(context.Background(), langsmith.RunStreamingQueryV2Params{
		AISearch: langsmith.F("question about langchain"),
	})
	var got []string
	for iter.Next() {
		got = append(got, iter.Current().ID+":"+iter.Current().Name)
	}
	if err := iter.Err(); err != nil {
		t.Fatalf("err = %v", err)
	}
	want := []string{"run-1:a", "run-2:b", "run-3:c", "run-4:d"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("items = %v, want %v", got, want)
	}
	if len(*bodies) != 2 {
		t.Fatalf("requests = %d, want 2", len(*bodies))
	}
	if c := (*bodies)[1]["cursor"]; c != "abc" {
		t.Fatalf("second request cursor = %v, want abc", c)
	}
	if q := (*bodies)[1]["ai_search"]; q != "question about langchain" {
		t.Fatalf("second request lost ai_search: %v", q)
	}
}

func TestQueryStreamingAutoPagingYieldsIndexedInsertsInFinalOrder(t *testing.T) {
	// The threads stream inserts each item at its sorted position
	// (`add /items/<n>`) instead of appending; deeper paths are ignored.
	srv, _ := sseServer(t, map[string][]string{
		"": {
			data(`[{"op":"add","path":"","value":{"items":[],"next_cursor":null}}]`),
			data(`[{"op":"add","path":"/items/0","value":{"thread_id":"t-1"}},{"op":"add","path":"/items/0","value":{"thread_id":"t-2"}}]`),
			data(`[{"op":"add","path":"/items/1/feedback_stats","value":{}},{"op":"replace","path":"/next_cursor","value":""}]`),
			end,
		},
	})
	iter := newSSEClient(srv).Threads.Streaming.QueryStreamingAutoPaging(context.Background(), langsmith.ThreadStreamingQueryParams{})
	var got []string
	for iter.Next() {
		got = append(got, iter.Current().ThreadID)
	}
	if err := iter.Err(); err != nil {
		t.Fatalf("err = %v", err)
	}
	if fmt.Sprint(got) != "[t-2 t-1]" {
		t.Fatalf("items = %v, want [t-2 t-1]", got)
	}
}

func TestQueryStreamingAutoPagingYieldsTraces(t *testing.T) {
	srv, _ := sseServer(t, map[string][]string{
		"": {
			data(`[{"op":"add","path":"","value":{"items":[{"root_run":{"id":"run-1"}}],"next_cursor":null}}]`),
			data(`[{"op":"add","path":"/items/-","value":{"root_run":{"id":"run-2"}}},{"op":"add","path":"/items/0/trace_aggregates","value":{"total_tokens":3}}]`),
			end,
		},
	})
	iter := newSSEClient(srv).Traces.Streaming.QueryStreamingAutoPaging(context.Background(), langsmith.TraceStreamingQueryParams{})
	var got []string
	for iter.Next() {
		got = append(got, iter.Current().RootRun.ID)
	}
	if err := iter.Err(); err != nil {
		t.Fatalf("err = %v", err)
	}
	if fmt.Sprint(got) != "[run-1 run-2]" {
		t.Fatalf("items = %v, want [run-1 run-2]", got)
	}
}

func TestQueryStreamingAutoPagingStopsOnEmptyPage(t *testing.T) {
	srv, bodies := sseServer(t, map[string][]string{
		"": {data(`[{"op":"add","path":"","value":{"items":[],"next_cursor":null}}]`), end},
	})
	iter := newSSEClient(srv).Runs.Streaming.QueryV2StreamingAutoPaging(context.Background(), langsmith.RunStreamingQueryV2Params{})
	if iter.Next() {
		t.Fatal("Next() = true on an empty page")
	}
	if err := iter.Err(); err != nil {
		t.Fatalf("err = %v", err)
	}
	if len(*bodies) != 1 {
		t.Fatalf("requests = %d, want 1", len(*bodies))
	}
}

func TestQueryStreamingAutoPagingReportsStreamError(t *testing.T) {
	srv, _ := sseServer(t, map[string][]string{
		"": {
			data(`[{"op":"add","path":"","value":{"items":[{"id":"run-1"}],"next_cursor":null}}]`),
			"event: error\ndata: {\"detail\":\"boom\"}\n\n",
		},
	})
	iter := newSSEClient(srv).Runs.Streaming.QueryV2StreamingAutoPaging(context.Background(), langsmith.RunStreamingQueryV2Params{})
	if !iter.Next() || iter.Current().ID != "run-1" {
		t.Fatalf("first item not delivered before the error")
	}
	if iter.Next() {
		t.Fatal("Next() = true after an error event")
	}
	if iter.Err() == nil {
		t.Fatal("err = nil, want the stream error")
	}
}

func TestQueryStreamingAutoPagingYieldsDecodedItemsBeforeError(t *testing.T) {
	srv, _ := sseServer(t, map[string][]string{
		"": {data(`[{"op":"add","path":"/items/-","value":{"id":"run-1"}},{"op":"add","path":"","value":"not a document"}]`), end},
	})
	iter := newSSEClient(srv).Runs.Streaming.QueryV2StreamingAutoPaging(context.Background(), langsmith.RunStreamingQueryV2Params{})
	if !iter.Next() || iter.Current().ID != "run-1" {
		t.Fatal("run-1 decoded before the bad op was not yielded")
	}
	if iter.Next() {
		t.Fatal("Next() = true after the decode error")
	}
	if iter.Err() == nil {
		t.Fatal("err = nil, want a decode error")
	}
}

func TestQueryStreamingAutoPagingStopsAfterClose(t *testing.T) {
	srv, bodies := sseServer(t, map[string][]string{
		"": {data(`[{"op":"add","path":"","value":{"items":[{"id":"run-1"},{"id":"run-2"}],"next_cursor":"abc"}}]`), end},
	})
	iter := newSSEClient(srv).Runs.Streaming.QueryV2StreamingAutoPaging(context.Background(), langsmith.RunStreamingQueryV2Params{})
	if !iter.Next() || iter.Current().ID != "run-1" {
		t.Fatal("first item not delivered")
	}
	if err := iter.Close(); err != nil {
		t.Fatalf("Close() = %v", err)
	}
	if iter.Next() {
		t.Fatalf("Next() = true after Close(), yielded %s", iter.Current().ID)
	}
	if err := iter.Err(); err != nil {
		t.Fatalf("err = %v", err)
	}
	if len(*bodies) != 1 {
		t.Fatalf("requests = %d, want 1", len(*bodies))
	}
}

func TestQueryStreamingAutoPagingClosesOnDecodeError(t *testing.T) {
	closed := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, data(`[{"op":"add","path":"","value":"not a document"}]`))
		w.(http.Flusher).Flush()
		<-r.Context().Done() // the handler returns only once the client closes the body
		close(closed)
	}))
	t.Cleanup(srv.Close)
	iter := newSSEClient(srv).Runs.Streaming.QueryV2StreamingAutoPaging(context.Background(), langsmith.RunStreamingQueryV2Params{})
	if iter.Next() {
		t.Fatal("Next() = true on a bad document")
	}
	if iter.Err() == nil {
		t.Fatal("err = nil, want a decode error")
	}
	select {
	case <-closed:
	case <-time.After(2 * time.Second):
		t.Fatal("response body still open after the decode error")
	}
}

func TestQueryStreamingAutoPagingReportsHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotAcceptable)
	}))
	t.Cleanup(srv.Close)
	iter := newSSEClient(srv).Runs.Streaming.QueryV2StreamingAutoPaging(context.Background(), langsmith.RunStreamingQueryV2Params{})
	if iter.Next() {
		t.Fatal("Next() = true on a 406")
	}
	if iter.Err() == nil {
		t.Fatal("err = nil, want the HTTP error")
	}
}
