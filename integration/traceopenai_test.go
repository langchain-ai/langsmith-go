package integration

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/sashabaranov/go-openai"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/instrumentation/traceopenai"
)

// Hardcoded run names per test type for identifying traces in the shared integration project.
const (
	runNameOpenAINonstreaming         = "openai_nonstreaming"
	runNameOpenAIStreaming            = "openai_streaming"
	runNameOpenAIEarlyTermination     = "openai_early_termination"
	runNameOpenAIError                = "openai_error"
	runNameOpenAIToolCallingNonstream = "openai_tool_calling_nonstream"
	runNameOpenAIToolCallingStream    = "openai_tool_calling_stream"
	runNameOpenAISystemMessage        = "openai_system_message"
	runNameOpenAIMultipleMessages     = "openai_multiple_messages"
	runNameOpenAITokenUsageNonstream  = "openai_token_usage_nonstream"
	runNameOpenAITokenUsageStream     = "openai_token_usage_stream"
)

func newOpenAIClient(t *testing.T, tp *sdktrace.TracerProvider) *openai.Client {
	t.Helper()
	apiKey := os.Getenv("OPENAI_API_KEY")
	mockURL, usingMock := mockBaseURL("openai")
	if usingMock {
		apiKey = "fake"
	}
	cfg := openai.DefaultConfig(apiKey)
	if usingMock {
		cfg.BaseURL = mockURL + "/v1"
	}

	cfg.HTTPClient = traceopenai.Client(traceopenai.WithTracerProvider(tp))
	return openai.NewClientWithConfig(cfg)
}

// --- Basic Non-Streaming ---

func TestOpenAI_NonStreaming(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newOpenAIClient(t, tt.TP)
	expectedRunName := runNameOpenAINonstreaming
	ctx := traceopenai.WithRunNameContext(context.Background(), expectedRunName)

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: "Say 'foo'"},
		},
		MaxTokens:   5,
		Temperature: 0,
	})
	if err != nil {
		t.Fatalf("CreateChatCompletion: %v", err)
	}
	if len(resp.Choices) == 0 {
		t.Fatal("expected at least one choice")
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
	if v, ok := getSpanAttr(spans, "gen_ai.system"); !ok || v != "openai" {
		t.Errorf("gen_ai.system = %q, want 'openai'", v)
	}
	if v, ok := getSpanAttr(spans, "gen_ai.request.model"); !ok || !strings.Contains(v, "gpt") {
		t.Errorf("gen_ai.request.model = %q", v)
	}
	if _, ok := getSpanAttr(spans, "gen_ai.prompt"); !ok {
		t.Error("expected gen_ai.prompt attribute")
	}
	if _, ok := getSpanAttr(spans, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion attribute")
	}
	// Note: traceopenai sets gen_ai.request.model but not gen_ai.response.model

	if tt.SendsToLangSmith() {
		runs := pollForRuns(t, tt.Project, 1, expectedRunName)
		if len(runs) == 0 {
			t.Error("expected at least one run in LangSmith")
		} else {
			// Full response is provider/model-dependent text; assert output presence and run shape.
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

func TestOpenAI_Streaming(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newOpenAIClient(t, tt.TP)
	expectedRunName := runNameOpenAIStreaming
	ctx := traceopenai.WithRunNameContext(context.Background(), expectedRunName)

	stream, err := client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: "Count from 1 to 5"},
		},
		MaxTokens:   50,
		Temperature: 0,
		StreamOptions: &openai.StreamOptions{
			IncludeUsage: true,
		},
	})
	if err != nil {
		t.Fatalf("CreateChatCompletionStream: %v", err)
	}

	var chunks int
	var fullContent strings.Builder
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("stream recv: %v", err)
		}
		chunks++
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			fullContent.WriteString(chunk.Choices[0].Delta.Content)
		}
	}
	stream.Close()
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

	if tt.SendsToLangSmith() {
		runs := pollForRuns(t, tt.Project, 1, expectedRunName)
		if len(runs) == 0 {
			t.Error("expected at least one run in LangSmith for streaming")
		} else {
			// Reduced/assembled stream in outputs (behavior 2)
			assertLangSmithRunFields(t, &runs[0], LangSmithRunAssertions{
				WantName:    expectedRunName,
				WantRunType: langsmith.RunTypeEnumLlm,
				WantInputs:  true,
				WantOutputs: true,
			})
		}
	}
}

// --- Early Stream Termination ---

func TestOpenAI_EarlyStreamTermination(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newOpenAIClient(t, tt.TP)
	expectedRunName := runNameOpenAIEarlyTermination
	ctx := traceopenai.WithRunNameContext(context.Background(), expectedRunName)

	stream, err := client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: "Write a long essay about Go programming"},
		},
		MaxTokens: 200,
	})
	if err != nil {
		t.Fatalf("CreateChatCompletionStream: %v", err)
	}

	for i := 0; i < 3; i++ {
		_, err := stream.Recv()
		if err != nil {
			break
		}
	}
	stream.Close()
	tt.TP.ForceFlush(context.Background())

	spans := tt.Spans()
	if len(spans) == 0 {
		t.Fatal("expected at least one span even with early termination")
	}
	if tt.SendsToLangSmith() {
		runs := pollForRuns(t, tt.Project, 1, expectedRunName)
		if len(runs) == 0 {
			t.Error("expected at least one run in LangSmith after early stream termination")
		} else {
			assertLangSmithRunFields(t, &runs[0], LangSmithRunAssertions{
				WantName:    expectedRunName,
				WantRunType: langsmith.RunTypeEnumLlm,
				WantInputs:  true,
				WantOutputs: true,
				ExpectError: true,
			})
		}
	}
}

// --- Error Handling ---

func TestOpenAI_ErrorHandling(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	expectedRunName := runNameOpenAIError

	cfg := openai.DefaultConfig("sk-invalid-key-for-testing")
	cfg.HTTPClient = traceopenai.Client(traceopenai.WithTracerProvider(tt.TP))
	client := openai.NewClientWithConfig(cfg)
	ctx := traceopenai.WithRunNameContext(context.Background(), expectedRunName)

	_, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: "test"},
		},
	})
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

func TestOpenAI_ToolCalling_NonStreaming(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newOpenAIClient(t, tt.TP)
	expectedRunName := runNameOpenAIToolCallingNonstream
	ctx := traceopenai.WithRunNameContext(context.Background(), expectedRunName)

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: "What's the weather like in Paris?"},
		},
		Tools: []openai.Tool{
			{
				Type: openai.ToolTypeFunction,
				Function: &openai.FunctionDefinition{
					Name:        "get_weather",
					Description: "Get the current weather in a given city.",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"location": map[string]interface{}{
								"type":        "string",
								"description": "The city name",
							},
						},
						"required": []string{"location"},
					},
				},
			},
		},
		ToolChoice: "auto",
	})
	if err != nil {
		t.Fatalf("CreateChatCompletion with tools: %v", err)
	}
	if len(resp.Choices) == 0 {
		t.Fatal("expected at least one choice")
	}

	choice := resp.Choices[0]
	if choice.FinishReason != "tool_calls" {
		t.Skipf("model didn't choose to call a tool (finish_reason=%s), skipping assertion", choice.FinishReason)
	}
	if len(choice.Message.ToolCalls) == 0 {
		t.Error("expected tool calls in response")
	} else if choice.Message.ToolCalls[0].Function.Name != "get_weather" {
		t.Errorf("tool name = %q, want 'get_weather'", choice.Message.ToolCalls[0].Function.Name)
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

func TestOpenAI_ToolCalling_Streaming(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newOpenAIClient(t, tt.TP)
	expectedRunName := runNameOpenAIToolCallingStream
	ctx := traceopenai.WithRunNameContext(context.Background(), expectedRunName)

	stream, err := client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: "What's the weather like in Paris?"},
		},
		Tools: []openai.Tool{
			{
				Type: openai.ToolTypeFunction,
				Function: &openai.FunctionDefinition{
					Name:        "get_weather",
					Description: "Get the current weather in a given city.",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"location": map[string]interface{}{
								"type":        "string",
								"description": "The city name",
							},
						},
						"required": []string{"location"},
					},
				},
			},
		},
		ToolChoice: "auto",
	})
	if err != nil {
		t.Fatalf("streaming with tools: %v", err)
	}

	var hasToolCallChunk bool
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("stream recv: %v", err)
		}
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.ToolCalls != nil {
			hasToolCallChunk = true
		}
	}
	stream.Close()
	tt.TP.ForceFlush(context.Background())

	if !hasToolCallChunk {
		t.Skip("model didn't stream tool calls, skipping")
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
			// Outputs include tool_calls (behavior 5)
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

func TestOpenAI_SystemMessage(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newOpenAIClient(t, tt.TP)
	expectedRunName := runNameOpenAISystemMessage
	ctx := traceopenai.WithRunNameContext(context.Background(), expectedRunName)

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "You always respond with exactly one word."},
			{Role: openai.ChatMessageRoleUser, Content: "What color is the sky?"},
		},
		MaxTokens:   10,
		Temperature: 0,
	})
	if err != nil {
		t.Fatalf("CreateChatCompletion: %v", err)
	}
	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
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
			// System content in run inputs (behavior 6)
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

func TestOpenAI_MultipleMessages(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newOpenAIClient(t, tt.TP)
	expectedRunName := runNameOpenAIMultipleMessages
	ctx := traceopenai.WithRunNameContext(context.Background(), expectedRunName)

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "You are a helpful math tutor."},
			{Role: openai.ChatMessageRoleUser, Content: "What is 2+2?"},
			{Role: openai.ChatMessageRoleAssistant, Content: "4"},
			{Role: openai.ChatMessageRoleUser, Content: "And what is that times 3?"},
		},
		MaxTokens:   10,
		Temperature: 0,
	})
	if err != nil {
		t.Fatalf("CreateChatCompletion: %v", err)
	}
	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
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
		if len(msgs) != 4 {
			t.Errorf("expected 4 messages in prompt, got %d", len(msgs))
		}
	}

	if tt.SendsToLangSmith() {
		runs := pollForRuns(t, tt.Project, 1, expectedRunName)
		if len(runs) == 0 {
			t.Error("expected at least one run in LangSmith for multiple messages")
		} else {
			// Full conversation in run inputs (behavior 7)
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

func TestOpenAI_TokenUsage_NonStreaming(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newOpenAIClient(t, tt.TP)
	expectedRunName := runNameOpenAITokenUsageNonstream
	ctx := traceopenai.WithRunNameContext(context.Background(), expectedRunName)

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: "Say hello"},
		},
		MaxTokens: 5,
	})
	if err != nil {
		t.Fatalf("CreateChatCompletion: %v", err)
	}
	if resp.Usage.PromptTokens == 0 {
		t.Error("expected non-zero prompt tokens from OpenAI")
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
			// Run has usage (behavior 8)
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

func TestOpenAI_TokenUsage_Streaming(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newOpenAIClient(t, tt.TP)
	expectedRunName := runNameOpenAITokenUsageStream
	ctx := traceopenai.WithRunNameContext(context.Background(), expectedRunName)

	stream, err := client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: "Say hello"},
		},
		MaxTokens: 5,
		StreamOptions: &openai.StreamOptions{
			IncludeUsage: true,
		},
	})
	if err != nil {
		t.Fatalf("CreateChatCompletionStream: %v", err)
	}

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("stream recv: %v", err)
		}
	}
	stream.Close()
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
			// Run has usage (behavior 8)
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
