package langsmith

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

type tracingTransportFunc func(*http.Request) (*http.Response, error)

func (f tracingTransportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestTracingHTTPClient(t *testing.T) {
	for _, status := range []int{http.StatusAccepted, http.StatusTemporaryRedirect} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			requests := make(chan *http.Request, 2)
			client := &http.Client{
				Timeout:       time.Second,
				CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
				Transport: tracingTransportFunc(func(r *http.Request) (*http.Response, error) {
					if _, err := io.Copy(io.Discard, r.Body); err != nil {
						return nil, err
					}
					requests <- r
					return &http.Response{StatusCode: status, Header: http.Header{"Location": {"https://redirect.example.com"}}, Body: io.NopCloser(strings.NewReader("")), Request: r}, nil
				}),
			}
			drain := DefaultDrainConfig()
			drain.DrainInterval, drain.CloseTimeout, drain.MaxWorkers = time.Hour, time.Second, 1
			tracer, err := NewTracingClient(context.Background(), WithTracingAPIKey("test-placeholder"),
				WithTracingAPIURL("https://example.com"), WithTracingProject("test-project"),
				WithSampleRate(1), WithTracingDrain(drain), WithTracingHTTPClient(client))
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(tracer.Close)
			id, err := uuid.NewV7()
			if err != nil {
				t.Fatal(err)
			}
			if err := tracer.CreateRun(&RunCreate{ID: id, TraceID: id, Name: "transport-test", RunType: "llm", StartTime: time.Now(), EndTime: time.Now()}); err != nil {
				t.Fatal(err)
			}
			tracer.Close()
			if len(requests) != 1 {
				t.Fatalf("expected one export, got %d", len(requests))
			}
			r := <-requests
			if r.URL.String() != "https://example.com/runs/multipart" {
				t.Fatalf("unexpected export URL: %s", r.URL)
			}
			if r.Header.Get("X-API-Key") != "test-placeholder" {
				t.Fatal("missing configured API key")
			}
		})
	}
}
