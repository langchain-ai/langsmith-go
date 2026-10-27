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
)

// A step is either an SSE event payload or a gate the server waits on before going on.
type step struct {
	patch string
	gate  chan struct{}
}

func ev(patch string) step { return step{patch: patch} }

func gate(c chan struct{}) step { return step{gate: c} }

// gatedServer answers /api/v2/runs/query with the steps keyed by the request's cursor, flushing each
// event and pausing on gates, and records every request body.
func gatedServer(t *testing.T, pages map[string][]step) (*httptest.Server, func() []map[string]any) {
	t.Helper()
	var mu sync.Mutex
	var bodies []map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		mu.Lock()
		bodies = append(bodies, body)
		mu.Unlock()
		cursor, _ := body["cursor"].(string)
		w.Header().Set("Content-Type", "text/event-stream")
		for _, s := range pages[cursor] {
			if s.gate != nil {
				select {
				case <-s.gate:
				case <-r.Context().Done():
					return
				}
				continue
			}
			fmt.Fprint(w, data(s.patch))
			w.(http.Flusher).Flush()
		}
	}))
	t.Cleanup(srv.Close)
	return srv, func() []map[string]any {
		mu.Lock()
		defer mu.Unlock()
		return append([]map[string]any(nil), bodies...)
	}
}

type runIter interface {
	Next() bool
	Current() langsmith.Run
	Err() error
	Close() error
}

func runPager(srv *httptest.Server) runIter {
	return newSSEClient(srv).Runs.Streaming.QueryV2StreamingAutoPaging(context.Background(), langsmith.RunStreamingQueryV2Params{
		AISearch: langsmith.F("q"),
	})
}

func drainIDs(t *testing.T, it runIter) string {
	t.Helper()
	var got []string
	for it.Next() {
		got = append(got, it.Current().ID)
	}
	if err := it.Err(); err != nil {
		t.Fatal(err)
	}
	return fmt.Sprint(got)
}

// nextAsync calls Next in the background so a test can tell "not yielded yet" from "yielded".
func nextAsync(it runIter) <-chan bool {
	c := make(chan bool, 1)
	go func() { c <- it.Next() }()
	return c
}

func rootDoc(items string) string {
	return `[{"op":"add","path":"","value":{"items":` + items + `,"next_cursor":null}}]`
}

func TestQueryStreamingAutoPagingYieldsAppendsBeforeTheStreamEnds(t *testing.T) {
	g := make(chan struct{})
	srv, _ := gatedServer(t, map[string][]step{"": {
		ev(rootDoc(`[{"id":"r1"}]`)),
		gate(g),
		ev(`[{"op":"add","path":"/items/-","value":{"id":"r2"}}]`),
	}})
	it := runPager(srv)
	defer it.Close()
	select {
	case ok := <-nextAsync(it):
		if !ok || it.Current().ID != "r1" {
			t.Fatalf("first item = %v %q", ok, it.Current().ID)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("first item not yielded while the stream is still open")
	}
	close(g)
	if got := drainIDs(t, it); got != "[r2]" {
		t.Fatalf("rest = %s", got)
	}
}

func TestQueryStreamingAutoPagingHoldsIndexedInsertsUntilThePageEnds(t *testing.T) {
	g := make(chan struct{})
	srv, _ := gatedServer(t, map[string][]step{"": {
		ev(rootDoc(`[]`)),
		ev(`[{"op":"add","path":"/items/0","value":{"id":"B"}}]`),
		gate(g),
		ev(`[{"op":"add","path":"/items/0","value":{"id":"A"}}]`),
		ev(`[{"op":"add","path":"/items/2","value":{"id":"C"}}]`),
	}})
	it := runPager(srv)
	defer it.Close()
	n := nextAsync(it)
	select {
	case <-n:
		t.Fatal("an item was yielded before the end of a page with indexed inserts")
	case <-time.After(300 * time.Millisecond):
	}
	close(g)
	if !<-n {
		t.Fatalf("no item, err = %v", it.Err())
	}
	got := []string{it.Current().ID}
	for it.Next() {
		got = append(got, it.Current().ID)
	}
	if it.Err() != nil || fmt.Sprint(got) != "[A B C]" {
		t.Fatalf("items = %v, err = %v", got, it.Err())
	}
}

func TestQueryStreamingAutoPagingFollowsCursorAcrossBufferedPages(t *testing.T) {
	srv, bodies := gatedServer(t, map[string][]step{
		"": {
			ev(rootDoc(`[{"id":"1"}]`)),
			ev(`[{"op":"add","path":"/items/-","value":{"id":"2"}},{"op":"replace","path":"/next_cursor","value":"abc"}]`),
		},
		"abc": {
			ev(rootDoc(`[]`)),
			ev(`[{"op":"add","path":"/items/0","value":{"id":"4"}}]`),
			ev(`[{"op":"add","path":"/items/0","value":{"id":"3"}},{"op":"replace","path":"/next_cursor","value":"def"}]`),
		},
		"def": {
			ev(rootDoc(`[{"id":"5"}]`)),
			ev(`[{"op":"replace","path":"/next_cursor","value":""}]`),
		},
	})
	it := runPager(srv)
	defer it.Close()
	if got := drainIDs(t, it); got != "[1 2 3 4 5]" {
		t.Fatalf("items = %s", got)
	}
	b := bodies()
	if len(b) != 3 || b[1]["cursor"] != "abc" || b[2]["cursor"] != "def" || b[1]["ai_search"] != "q" {
		t.Fatalf("bodies = %v", b)
	}
}

func TestQueryStreamingAutoPagingIgnoresDeeperPatches(t *testing.T) {
	srv, _ := gatedServer(t, map[string][]step{"": {
		ev(rootDoc(`[{"id":"1"}]`)),
		ev(`[{"op":"add","path":"/items/0/extra","value":1},{"op":"replace","path":"/items/0/id","value":"changed"},{"op":"add","path":"/items/-","value":{"id":"2"}}]`),
	}})
	if got := drainIDs(t, runPager(srv)); got != "[1 2]" {
		t.Fatalf("items = %s", got)
	}
}

func TestQueryStreamingAutoPagingAppliesReplaceAndRemoveOnlyToUnyieldedItems(t *testing.T) {
	// Buffered page: the indexes are in the document that is not yielded yet.
	srv, _ := gatedServer(t, map[string][]step{"": {
		ev(rootDoc(`[]`)),
		ev(`[{"op":"add","path":"/items/0","value":{"id":"B"}},{"op":"add","path":"/items/1","value":{"id":"C"}}]`),
		ev(`[{"op":"replace","path":"/items/0","value":{"id":"Bx"}},{"op":"remove","path":"/items/1"},{"op":"add","path":"/items/-","value":{"id":"D"}}]`),
	}})
	if got := drainIDs(t, runPager(srv)); got != "[Bx D]" {
		t.Fatalf("buffered items = %s", got)
	}

	// Append page: r1 was yielded, so a later replace or remove of it is ignored.
	g := make(chan struct{})
	srv2, _ := gatedServer(t, map[string][]step{"": {
		ev(rootDoc(`[{"id":"r1"}]`)),
		gate(g),
		ev(`[{"op":"replace","path":"/items/0","value":{"id":"r1x"}},{"op":"remove","path":"/items/0"},{"op":"add","path":"/items/-","value":{"id":"r2"}}]`),
	}})
	it := runPager(srv2)
	defer it.Close()
	if !<-nextAsync(it) || it.Current().ID != "r1" {
		t.Fatalf("first = %q", it.Current().ID)
	}
	close(g)
	if got := drainIDs(t, it); got != "[r2]" {
		t.Fatalf("after yield = %s", got)
	}
}

func TestQueryStreamingAutoPagingPlacesInsertsBeforeYieldedItemsAfterThem(t *testing.T) {
	g := make(chan struct{})
	srv, _ := gatedServer(t, map[string][]step{"": {
		ev(rootDoc(`[{"id":"r1"}]`)),
		gate(g),
		ev(`[{"op":"add","path":"/items/0","value":{"id":"X"}}]`),
	}})
	it := runPager(srv)
	defer it.Close()
	if !<-nextAsync(it) || it.Current().ID != "r1" {
		t.Fatalf("first = %q", it.Current().ID)
	}
	close(g)
	if got := drainIDs(t, it); got != "[X]" {
		t.Fatalf("after yield = %s", got)
	}
}
