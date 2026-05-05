//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"google.golang.org/genai"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/instrumentation/tracegemini"
)

const testGeminiModel = "gemini-2.5-flash"

const (
	runNameGeminiNonstreaming         = "gemini_nonstreaming"
	runNameGeminiStreaming            = "gemini_streaming"
	runNameGeminiError                = "gemini_error"
	runNameGeminiToolCallingNonstream = "gemini_tool_calling_nonstream"
	runNameGeminiToolCallingStream    = "gemini_tool_calling_stream"
	runNameGeminiSystemMessage        = "gemini_system_message"
	runNameGeminiMultipleMessages     = "gemini_multiple_messages"
	runNameGeminiTokenUsageNonstream  = "gemini_token_usage_nonstream"
	runNameGeminiTokenUsageStream     = "gemini_token_usage_stream"
	runNameGeminiStopReason           = "gemini_stop_reason"
)

func newGeminiClient(t *testing.T, tp *sdktrace.TracerProvider) *genai.Client {
	t.Helper()
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		t.Skip("GOOGLE_API_KEY not set")
	}
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:     apiKey,
		Backend:    genai.BackendGeminiAPI,
		HTTPClient: tracegemini.Client(tracegemini.WithTracerProvider(tp)),
	})
	if err != nil {
		t.Fatalf("genai.NewClient: %v", err)
	}
	return client
}

// --- Basic Non-Streaming ---

func TestGemini_NonStreaming(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newGeminiClient(t, tt.TP)
	expectedRunName := runNameGeminiNonstreaming
	ctx := tracegemini.WithRunNameContext(context.Background(), expectedRunName)

	resp, err := client.Models.GenerateContent(ctx, testGeminiModel, genai.Text("Say 'foo'"), &genai.GenerateContentConfig{
		MaxOutputTokens: 5,
		Temperature:     genai.Ptr[float32](0),
	})
	if err != nil {
		t.Fatalf("GenerateContent: %v", err)
	}
	if len(resp.Candidates) == 0 {
		t.Fatal("expected at least one candidate")
	}

	tt.TP.ForceFlush(context.Background())
	spans := tt.Spans()
	if len(spans) == 0 {
		t.Fatal("expected at least one span")
	}

	found := false
	for _, s := range spans {
		if s.Name == expectedRunName {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected span named " + expectedRunName)
	}
	if v, ok := getSpanAttr(spans, "gen_ai.provider.name"); !ok || v != "gcp.gemini" {
		t.Errorf("gen_ai.provider.name = %q, want 'gcp.gemini'", v)
	}
	if v, ok := getSpanAttr(spans, "gen_ai.request.model"); !ok || !strings.Contains(v, "gemini") {
		t.Errorf("gen_ai.request.model = %q", v)
	}
	if _, ok := getSpanAttr(spans, "gen_ai.prompt"); !ok {
		t.Error("expected gen_ai.prompt attribute")
	}
	if _, ok := getSpanAttr(spans, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion attribute")
	}
	if _, ok := getSpanAttr(spans, "gen_ai.response.model"); !ok {
		t.Error("expected gen_ai.response.model attribute")
	}

	if tt.SendsToLangSmith() {
		runs := pollForRuns(t, tt.Project, 1, expectedRunName)
		if len(runs) == 0 {
			t.Error("expected at least one run in LangSmith")
		} else {
			assertLangSmithRunFields(t, &runs[0], LangSmithRunAssertions{
				WantName:    expectedRunName,
				WantRunType: langsmith.RunTypeEnumLlm,
				WantInputs:  true,
				WantOutputs: true,
			})
		}
	}
}

// --- Streaming ---

func TestGemini_Streaming(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newGeminiClient(t, tt.TP)
	expectedRunName := runNameGeminiStreaming
	ctx := tracegemini.WithRunNameContext(context.Background(), expectedRunName)

	var chunks int
	var fullContent strings.Builder
	for result, err := range client.Models.GenerateContentStream(ctx, testGeminiModel, genai.Text("Count from 1 to 5"), &genai.GenerateContentConfig{
		MaxOutputTokens: 50,
		Temperature:     genai.Ptr[float32](0),
	}) {
		if err != nil {
			t.Fatalf("stream chunk: %v", err)
		}
		chunks++
		if len(result.Candidates) > 0 && result.Candidates[0].Content != nil {
			for _, part := range result.Candidates[0].Content.Parts {
				if part.Text != "" {
					fullContent.WriteString(part.Text)
				}
			}
		}
	}

	tt.TP.ForceFlush(context.Background())

	if chunks == 0 {
		t.Fatal("expected at least one chunk")
	}
	if fullContent.Len() == 0 {
		t.Error("expected non-empty streamed content")
	}

	spans := tt.Spans()
	if len(spans) == 0 {
		t.Fatal("expected at least one span")
	}
	found := false
	for _, s := range spans {
		if s.Name == expectedRunName {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected span named " + expectedRunName + " in streaming")
	}
	if _, ok := getSpanAttr(spans, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion attribute from streaming")
	}
	if _, ok := getSpanAttr(spans, "gen_ai.response.model"); !ok {
		t.Error("expected gen_ai.response.model from streaming")
	}

	if tt.SendsToLangSmith() {
		runs := pollForRuns(t, tt.Project, 1, expectedRunName)
		if len(runs) == 0 {
			t.Error("expected at least one run in LangSmith for streaming")
		} else {
			assertLangSmithRunFields(t, &runs[0], LangSmithRunAssertions{
				WantName:    expectedRunName,
				WantRunType: langsmith.RunTypeEnumLlm,
				WantInputs:  true,
				WantOutputs: true,
			})
		}
	}
}

// --- Error Handling ---

func TestGemini_ErrorHandling(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	expectedRunName := runNameGeminiError

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:     "invalid-key-for-testing",
		Backend:    genai.BackendGeminiAPI,
		HTTPClient: tracegemini.Client(tracegemini.WithTracerProvider(tt.TP)),
	})
	if err != nil {
		t.Fatalf("genai.NewClient: %v", err)
	}
	ctx := tracegemini.WithRunNameContext(context.Background(), expectedRunName)

	_, err = client.Models.GenerateContent(ctx, testGeminiModel, genai.Text("test"), nil)
	if err == nil {
		t.Fatal("expected error with invalid API key")
	}

	tt.TP.ForceFlush(context.Background())
	spans := tt.Spans()
	if len(spans) == 0 {
		t.Fatal("expected span even on error")
	}
	if tt.SendsToLangSmith() {
		runs := pollForRuns(t, tt.Project, 1, expectedRunName)
		if len(runs) == 0 {
			t.Error("expected at least one run in LangSmith even on error")
		} else {
			assertLangSmithRunFields(t, &runs[0], LangSmithRunAssertions{
				WantName:    expectedRunName,
				WantRunType: langsmith.RunTypeEnumLlm,
				ExpectError: true,
			})
		}
	}
}

// --- Tool/Function Calling (Non-Streaming) ---

func TestGemini_ToolCalling_NonStreaming(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newGeminiClient(t, tt.TP)
	expectedRunName := runNameGeminiToolCallingNonstream
	ctx := tracegemini.WithRunNameContext(context.Background(), expectedRunName)

	resp, err := client.Models.GenerateContent(ctx, testGeminiModel, genai.Text("What's the weather like in Paris?"), &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{
				FunctionDeclarations: []*genai.FunctionDeclaration{
					{
						Name:        "get_weather",
						Description: "Get the current weather in a given city.",
						Parameters: &genai.Schema{
							Type: genai.TypeObject,
							Properties: map[string]*genai.Schema{
								"location": {
									Type:        genai.TypeString,
									Description: "The city name",
								},
							},
							Required: []string{"location"},
						},
					},
				},
			},
		},
		ToolConfig: &genai.ToolConfig{
			FunctionCallingConfig: &genai.FunctionCallingConfig{
				Mode: genai.FunctionCallingConfigModeAny,
			},
		},
	})
	if err != nil {
		t.Fatalf("GenerateContent with tools: %v", err)
	}
	if len(resp.Candidates) == 0 {
		t.Fatal("expected at least one candidate")
	}

	hasFunctionCall := false
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.FunctionCall != nil {
			hasFunctionCall = true
			if part.FunctionCall.Name != "get_weather" {
				t.Errorf("function name = %q, want 'get_weather'", part.FunctionCall.Name)
			}
		}
	}
	if !hasFunctionCall {
		t.Skip("model didn't call function, skipping tool assertions")
	}

	tt.TP.ForceFlush(context.Background())
	spans := tt.Spans()

	if v, ok := getSpanAttr(spans, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion for tool call response")
	} else if !strings.Contains(v, "get_weather") {
		t.Errorf("completion should contain 'get_weather': %s", v)
	}

	if tt.SendsToLangSmith() {
		runs := pollForRuns(t, tt.Project, 1, expectedRunName)
		if len(runs) == 0 {
			t.Error("expected at least one run in LangSmith for tool calling")
		} else {
			assertLangSmithRunFields(t, &runs[0], LangSmithRunAssertions{
				WantName:                expectedRunName,
				WantRunType:             langsmith.RunTypeEnumLlm,
				WantInputs:              true,
				WantOutputs:             true,
				WantOutputsContainTools: true,
			})
		}
	}
}

// --- Tool/Function Calling (Streaming) ---

func TestGemini_ToolCalling_Streaming(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newGeminiClient(t, tt.TP)
	expectedRunName := runNameGeminiToolCallingStream
	ctx := tracegemini.WithRunNameContext(context.Background(), expectedRunName)

	hasFunctionCall := false
	for result, err := range client.Models.GenerateContentStream(ctx, testGeminiModel, genai.Text("What's the weather like in Paris?"), &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{
				FunctionDeclarations: []*genai.FunctionDeclaration{
					{
						Name:        "get_weather",
						Description: "Get the current weather in a given city.",
						Parameters: &genai.Schema{
							Type: genai.TypeObject,
							Properties: map[string]*genai.Schema{
								"location": {
									Type:        genai.TypeString,
									Description: "The city name",
								},
							},
							Required: []string{"location"},
						},
					},
				},
			},
		},
		ToolConfig: &genai.ToolConfig{
			FunctionCallingConfig: &genai.FunctionCallingConfig{
				Mode: genai.FunctionCallingConfigModeAny,
			},
		},
	}) {
		if err != nil {
			t.Fatalf("stream chunk: %v", err)
		}
		if len(result.Candidates) > 0 && result.Candidates[0].Content != nil {
			for _, part := range result.Candidates[0].Content.Parts {
				if part.FunctionCall != nil {
					hasFunctionCall = true
				}
			}
		}
	}

	tt.TP.ForceFlush(context.Background())

	if !hasFunctionCall {
		t.Skip("model didn't stream function calls, skipping")
	}

	spans := tt.Spans()
	if v, ok := getSpanAttr(spans, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion for streaming tool calls")
	} else if !strings.Contains(v, "get_weather") {
		t.Errorf("streaming completion should contain 'get_weather': %s", v)
	}

	if tt.SendsToLangSmith() {
		runs := pollForRuns(t, tt.Project, 1, expectedRunName)
		if len(runs) == 0 {
			t.Error("expected at least one run in LangSmith for streaming tool calls")
		} else {
			assertLangSmithRunFields(t, &runs[0], LangSmithRunAssertions{
				WantName:                expectedRunName,
				WantRunType:             langsmith.RunTypeEnumLlm,
				WantInputs:              true,
				WantOutputs:             true,
				WantOutputsContainTools: true,
			})
		}
	}
}

// --- System Messages ---

func TestGemini_SystemMessage(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newGeminiClient(t, tt.TP)
	expectedRunName := runNameGeminiSystemMessage
	ctx := tracegemini.WithRunNameContext(context.Background(), expectedRunName)

	resp, err := client.Models.GenerateContent(ctx, testGeminiModel, genai.Text("What color is the sky?"), &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: "You always respond with exactly one word."}}},
		MaxOutputTokens:   10,
		Temperature:       genai.Ptr[float32](0),
	})
	if err != nil {
		t.Fatalf("GenerateContent: %v", err)
	}
	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		t.Fatal("expected non-empty response")
	}

	tt.TP.ForceFlush(context.Background())
	spans := tt.Spans()

	if v, ok := getSpanAttr(spans, "gen_ai.prompt"); !ok {
		t.Error("expected gen_ai.prompt")
	} else if !strings.Contains(v, "system") {
		t.Errorf("prompt should contain system role: %s", v)
	}

	if tt.SendsToLangSmith() {
		runs := pollForRuns(t, tt.Project, 1, expectedRunName)
		if len(runs) == 0 {
			t.Error("expected at least one run in LangSmith for system message")
		} else {
			assertLangSmithRunFields(t, &runs[0], LangSmithRunAssertions{
				WantName:                expectedRunName,
				WantRunType:             langsmith.RunTypeEnumLlm,
				WantInputs:              true,
				WantOutputs:             true,
				WantInputsContainSystem: true,
			})
		}
	}
}

// --- Multiple Messages ---

func TestGemini_MultipleMessages(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newGeminiClient(t, tt.TP)
	expectedRunName := runNameGeminiMultipleMessages
	ctx := tracegemini.WithRunNameContext(context.Background(), expectedRunName)

	resp, err := client.Models.GenerateContent(ctx, testGeminiModel,
		[]*genai.Content{
			{Role: "user", Parts: []*genai.Part{{Text: "What is 2+2?"}}},
			{Role: "model", Parts: []*genai.Part{{Text: "4"}}},
			{Role: "user", Parts: []*genai.Part{{Text: "And what is that times 3?"}}},
		},
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: "You are a math tutor. Be concise."}}},
			MaxOutputTokens:   10,
			Temperature:       genai.Ptr[float32](0),
		},
	)
	if err != nil {
		t.Fatalf("GenerateContent: %v", err)
	}
	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		t.Fatal("expected non-empty response")
	}

	tt.TP.ForceFlush(context.Background())
	spans := tt.Spans()

	if v, ok := getSpanAttr(spans, "gen_ai.prompt"); !ok {
		t.Error("expected gen_ai.prompt")
	} else {
		var prompt map[string]interface{}
		if err := json.Unmarshal([]byte(v), &prompt); err != nil {
			t.Fatalf("prompt is not valid JSON: %v", err)
		}
		msgs, ok := prompt["messages"].([]interface{})
		if !ok {
			t.Fatal("expected messages array in prompt")
		}
		// system + 3 messages = 4
		if len(msgs) != 4 {
			t.Errorf("expected 4 messages in prompt (1 system + 3), got %d", len(msgs))
		}
	}

	if tt.SendsToLangSmith() {
		runs := pollForRuns(t, tt.Project, 1, expectedRunName)
		if len(runs) == 0 {
			t.Error("expected at least one run in LangSmith for multiple messages")
		} else {
			assertLangSmithRunFields(t, &runs[0], LangSmithRunAssertions{
				WantName:               expectedRunName,
				WantRunType:            langsmith.RunTypeEnumLlm,
				WantInputs:             true,
				WantOutputs:            true,
				WantInputsMessageCount: 4,
			})
		}
	}
}

// --- Token Usage ---

func TestGemini_TokenUsage_NonStreaming(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newGeminiClient(t, tt.TP)
	expectedRunName := runNameGeminiTokenUsageNonstream
	ctx := tracegemini.WithRunNameContext(context.Background(), expectedRunName)

	resp, err := client.Models.GenerateContent(ctx, testGeminiModel, genai.Text("Say hello"), &genai.GenerateContentConfig{
		MaxOutputTokens: 5,
	})
	if err != nil {
		t.Fatalf("GenerateContent: %v", err)
	}
	if resp.UsageMetadata == nil || resp.UsageMetadata.PromptTokenCount == 0 {
		t.Error("expected non-zero prompt token count from Gemini")
	}

	tt.TP.ForceFlush(context.Background())
	spans := tt.Spans()

	if v, ok := getSpanAttrInt(spans, "gen_ai.usage.input_tokens"); !ok || v == 0 {
		t.Errorf("expected non-zero gen_ai.usage.input_tokens, got %d", v)
	}
	if v, ok := getSpanAttrInt(spans, "gen_ai.usage.output_tokens"); !ok || v == 0 {
		t.Errorf("expected non-zero gen_ai.usage.output_tokens, got %d", v)
	}

	if tt.SendsToLangSmith() {
		runs := pollForRuns(t, tt.Project, 1, expectedRunName)
		if len(runs) == 0 {
			t.Error("expected at least one run in LangSmith for token usage")
		} else {
			assertLangSmithRunFields(t, &runs[0], LangSmithRunAssertions{
				WantName:    expectedRunName,
				WantRunType: langsmith.RunTypeEnumLlm,
				WantInputs:  true,
				WantOutputs: true,
				WantUsage:   true,
			})
		}
	}
}

func TestGemini_TokenUsage_Streaming(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newGeminiClient(t, tt.TP)
	expectedRunName := runNameGeminiTokenUsageStream
	ctx := tracegemini.WithRunNameContext(context.Background(), expectedRunName)

	for _, err := range client.Models.GenerateContentStream(ctx, testGeminiModel, genai.Text("Say hello"), &genai.GenerateContentConfig{
		MaxOutputTokens: 5,
	}) {
		if err != nil {
			t.Fatalf("stream chunk: %v", err)
		}
	}

	tt.TP.ForceFlush(context.Background())
	spans := tt.Spans()

	if v, ok := getSpanAttrInt(spans, "gen_ai.usage.input_tokens"); !ok || v == 0 {
		t.Errorf("expected non-zero gen_ai.usage.input_tokens in streaming, got %d", v)
	}
	if v, ok := getSpanAttrInt(spans, "gen_ai.usage.output_tokens"); !ok || v == 0 {
		t.Errorf("expected non-zero gen_ai.usage.output_tokens in streaming, got %d", v)
	}

	if tt.SendsToLangSmith() {
		runs := pollForRuns(t, tt.Project, 1, expectedRunName)
		if len(runs) == 0 {
			t.Error("expected at least one run in LangSmith for streaming token usage")
		} else {
			assertLangSmithRunFields(t, &runs[0], LangSmithRunAssertions{
				WantName:    expectedRunName,
				WantRunType: langsmith.RunTypeEnumLlm,
				WantInputs:  true,
				WantOutputs: true,
				WantUsage:   true,
			})
		}
	}
}

// --- Stop Reason ---

func TestGemini_StopReason(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newGeminiClient(t, tt.TP)
	expectedRunName := runNameGeminiStopReason
	ctx := tracegemini.WithRunNameContext(context.Background(), expectedRunName)

	_, err := client.Models.GenerateContent(ctx, testGeminiModel, genai.Text("Say hello"), &genai.GenerateContentConfig{
		MaxOutputTokens: 5,
	})
	if err != nil {
		t.Fatalf("GenerateContent: %v", err)
	}

	tt.TP.ForceFlush(context.Background())
	spans := tt.Spans()

	if _, ok := getSpanAttr(spans, "langsmith.metadata.stop_reason"); !ok {
		t.Error("expected langsmith.metadata.stop_reason attribute")
	}

	if tt.SendsToLangSmith() {
		runs := pollForRuns(t, tt.Project, 1, expectedRunName)
		if len(runs) == 0 {
			t.Error("expected at least one run in LangSmith for stop reason")
		} else {
			assertLangSmithRunFields(t, &runs[0], LangSmithRunAssertions{
				WantName:    expectedRunName,
				WantRunType: langsmith.RunTypeEnumLlm,
				WantInputs:  true,
				WantOutputs: true,
			})
		}
	}
}
