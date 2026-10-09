package traceopenai_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/langchain-ai/langsmith-go/instrumentation/traceanthropic"
	"github.com/langchain-ai/langsmith-go/instrumentation/traceopenai"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

type requestCaptureTransport func(*http.Request) (*http.Response, error)

func (f requestCaptureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestRequestCaptureSampling(t *testing.T) {
	for _, tc := range []struct {
		name, path, request, prompt, response, streamResponse string
		wrap                                                  func(*http.Client, trace.TracerProvider) *http.Client
	}{
		{
			name: "anthropic", path: "/v1/messages",
			request:  `"model":"test-model","max_tokens":16,"temperature":0.5,"system":[{"type":"text","text":"Be helpful <&>"}],"messages":[{"role":"user","content":"hello"}]`,
			prompt:   `{"messages":[{"role":"system","content":[{"type":"text","text":"Be helpful <&>"}]},{"role":"user","content":"hello"}]}`,
			response: `{"model":"test-model","content":[{"type":"text","text":"hi"}],"usage":{"input_tokens":3,"output_tokens":2}}`,
			streamResponse: "data: {\"type\":\"message_start\",\"message\":{\"model\":\"test-model\",\"usage\":{\"input_tokens\":3}}}\n\n" +
				"data: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"text_delta\",\"text\":\"hi\"}}\n\n" +
				"data: {\"type\":\"message_delta\",\"usage\":{\"output_tokens\":2}}\n\n" +
				"data: {\"type\":\"message_stop\"}\n\n",
			wrap: func(c *http.Client, tp trace.TracerProvider) *http.Client {
				return traceanthropic.WrapClient(c, traceanthropic.WithTracerProvider(tp))
			},
		},
		{
			name: "openai-chat", path: "/v1/chat/completions",
			request:        `"model":"test-model","messages":[{"role":"system","content":"Be helpful <&>"},{"role":"user","content":"hello"}]`,
			prompt:         `{"messages":[{"role":"system","content":"Be helpful <&>"},{"role":"user","content":"hello"}]}`,
			response:       `{"choices":[{"message":{"role":"assistant","content":"hi"}}],"usage":{"prompt_tokens":3,"completion_tokens":2}}`,
			streamResponse: "data: {\"choices\":[{\"delta\":{\"content\":\"hi\"}}],\"usage\":{\"prompt_tokens\":3,\"completion_tokens\":2}}\n\ndata: [DONE]\n\n",
			wrap: func(c *http.Client, tp trace.TracerProvider) *http.Client {
				return traceopenai.WrapClient(c, traceopenai.WithTracerProvider(tp))
			},
		},
		{
			name: "openai-responses", path: "/v1/responses",
			request:        `"model":"test-model","instructions":"Be helpful <&>","input":[{"role":"user","content":[{"type":"input_text","text":"hello"}]}]`,
			prompt:         `{"messages":[{"role":"system","content":"Be helpful <&>"},{"role":"user","content":"hello"}]}`,
			response:       `{"output":[{"type":"message","content":[{"type":"output_text","text":"hi"}]}],"usage":{"input_tokens":3,"output_tokens":2}}`,
			streamResponse: "data: {\"type\":\"response.completed\",\"response\":{\"output\":[{\"type\":\"message\",\"content\":[{\"type\":\"output_text\",\"text\":\"hi\"}]}],\"usage\":{\"input_tokens\":3,\"output_tokens\":2}}}\n\n",
			wrap: func(c *http.Client, tp trace.TracerProvider) *http.Client {
				return traceopenai.WrapClient(c, traceopenai.WithTracerProvider(tp))
			},
		},
	} {
		for _, recording := range []bool{false, true} {
			for _, streaming := range []bool{false, true} {
				for _, replayable := range []bool{false, true} {
					t.Run(fmt.Sprintf("%s/recording=%t/streaming=%t/replayable=%t", tc.name, recording, streaming, replayable), func(t *testing.T) {
						exporter := tracetest.NewInMemoryExporter()
						parentTP := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
						t.Cleanup(func() { require.NoError(t, parentTP.Shutdown(context.Background())) })
						ctx, parent := parentTP.Tracer("test").Start(context.Background(), "parent")
						sampler := sdktrace.NeverSample()
						if recording {
							sampler = sdktrace.AlwaysSample()
						}
						childTP := sdktrace.NewTracerProvider(sdktrace.WithSampler(sampler), sdktrace.WithSyncer(exporter))
						t.Cleanup(func() { require.NoError(t, childTP.Shutdown(context.Background())) })
						body := fmt.Sprintf("{\n  %s, \"stream\":%t\n}\n", tc.request, streaming)
						response := tc.response
						if streaming {
							response = tc.streamResponse
						}
						client := tc.wrap(&http.Client{Transport: requestCaptureTransport(func(req *http.Request) (*http.Response, error) {
							got, err := io.ReadAll(req.Body)
							require.NoError(t, err)
							require.Equal(t, body, string(got))
							require.NoError(t, req.Body.Close())
							require.Equal(t, int64(len(body)), req.ContentLength)
							if replayable {
								replay, err := req.GetBody()
								require.NoError(t, err)
								got, err = io.ReadAll(replay)
								require.NoError(t, err)
								require.Equal(t, body, string(got))
								require.NoError(t, replay.Close())
							} else {
								require.Nil(t, req.GetBody)
							}
							return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(response))}, nil
						})}, childTP)
						req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://example.com"+tc.path, strings.NewReader(body))
						require.NoError(t, err)
						if !replayable {
							req.GetBody = nil
						}
						resp, err := client.Do(req)
						require.NoError(t, err)
						got, err := io.ReadAll(resp.Body)
						require.NoError(t, err)
						require.Equal(t, response, string(got))
						require.NoError(t, resp.Body.Close())
						parent.End()
						spans := exporter.GetSpans()
						wantSpans := 1
						if recording {
							wantSpans++
						}
						require.Len(t, spans, wantSpans)
						for _, span := range spans {
							attrs := make(map[string]attribute.Value)
							for _, attr := range span.Attributes {
								attrs[string(attr.Key)] = attr.Value
							}
							require.Equal(t, int64(3), attrs["gen_ai.usage.input_tokens"].AsInt64())
							require.Equal(t, int64(2), attrs["gen_ai.usage.output_tokens"].AsInt64())
							if span.Name != "parent" {
								require.JSONEq(t, tc.prompt, attrs["gen_ai.prompt"].AsString())
								require.Equal(t, "test-model", attrs["gen_ai.request.model"].AsString())
								require.JSONEq(t, `{"messages":[{"role":"assistant","content":`+completionContent(tc.name)+`}]}`, attrs["gen_ai.completion"].AsString())
								require.JSONEq(t, `{"input_tokens":3,"output_tokens":2,"total_tokens":5}`, attrs["langsmith.usage_metadata"].AsString())
								if tc.name == "anthropic" {
									require.Equal(t, int64(16), attrs["gen_ai.request.max_tokens"].AsInt64())
									require.Equal(t, 0.5, attrs["gen_ai.request.temperature"].AsFloat64())
								}
							}
						}
					})
				}
			}
		}
	}
}

func completionContent(provider string) string {
	if provider == "anthropic" {
		return `[{"type":"text","text":"hi"}]`
	}
	return `"hi"`
}
