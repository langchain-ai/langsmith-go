//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/instrumentation/traceanthropic"
)

const testAnthropicModel = "claude-sonnet-4-5-20250929"
const anthropicRunNameBase = "anthropic.messages"

func newAnthropicClient(t *testing.T, tp *sdktrace.TracerProvider, runNameSuffix string) anthropic.Client {
	t.Helper()
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY not set")
	}
	opts := []traceanthropic.Option{traceanthropic.WithTracerProvider(tp)}
	if runNameSuffix != "" {
		opts = append(opts, traceanthropic.WithRunNameSuffix(runNameSuffix))
	}
	return anthropic.NewClient(
		option.WithAPIKey(apiKey),
		option.WithHTTPClient(traceanthropic.Client(opts...)),
	)
}

// --- Basic Non-Streaming ---

func TestAnthropic_NonStreaming(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newAnthropicClient(t, tt.TP, t.Name())
	expectedRunName := anthropicRunNameBase + "__" + t.Name()

	msg, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.Model(testAnthropicModel),
		MaxTokens: 10,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Say 'foo'")),
		},
	})
	if err != nil {
		t.Fatalf("Messages.New: %v", err)
	}
	if len(msg.Content) == 0 {
		t.Fatal("expected at least one content block")
	}
	if msg.Content[0].Type != "text" || msg.Content[0].Text == "" {
		t.Errorf("expected text content, got type=%s", msg.Content[0].Type)
	}

	tt.TP.ForceFlush(context.Background())
	spans := tt.Spans()
	if len(spans) == 0 {
		t.Fatal("expected at least one span")
	}

	found := false
	for _, s := range spans {
		if s.Name == expectedRunName || strings.HasPrefix(s.Name, anthropicRunNameBase) {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected span named " + anthropicRunNameBase + " or " + expectedRunName)
	}
	if v, ok := getSpanAttr(spans, "gen_ai.system"); !ok || v != "anthropic" {
		t.Errorf("gen_ai.system = %q, want 'anthropic'", v)
	}
	if v, ok := getSpanAttr(spans, "gen_ai.request.model"); !ok || !strings.Contains(v, "claude") {
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
				WantRunType: langsmith.RunQueryResponseRunsRunTypeLlm,
				WantInputs:  true,
				WantOutputs: true,
			})
		}
	}
}

// --- Streaming ---

func TestAnthropic_Streaming(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newAnthropicClient(t, tt.TP, t.Name())
	expectedRunName := anthropicRunNameBase + "__" + t.Name()

	stream := client.Messages.NewStreaming(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.Model(testAnthropicModel),
		MaxTokens: 50,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Count from 1 to 5")),
		},
	})

	var chunks int
	var fullText strings.Builder
	accumulated := anthropic.Message{}
	for stream.Next() {
		event := stream.Current()
		accumulated.Accumulate(event)
		chunks++

		switch variant := event.AsAny().(type) {
		case anthropic.ContentBlockDeltaEvent:
			switch delta := variant.Delta.AsAny().(type) {
			case anthropic.TextDelta:
				fullText.WriteString(delta.Text)
			}
		}
	}
	if stream.Err() != nil {
		t.Fatalf("streaming error: %v", stream.Err())
	}
	if chunks == 0 {
		t.Fatal("expected at least one streaming event")
	}
	if fullText.Len() == 0 {
		t.Error("expected non-empty streamed text")
	}

	tt.TP.ForceFlush(context.Background())
	spans := tt.Spans()
	if len(spans) == 0 {
		t.Fatal("expected at least one span")
	}
	found := false
	for _, s := range spans {
		if s.Name == expectedRunName || strings.HasPrefix(s.Name, anthropicRunNameBase) {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected span named " + anthropicRunNameBase + " or " + expectedRunName + " in streaming")
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
			// Reduced/assembled stream in outputs (behavior 2)
			assertLangSmithRunFields(t, &runs[0], LangSmithRunAssertions{
				WantName:    expectedRunName,
				WantRunType: langsmith.RunQueryResponseRunsRunTypeLlm,
				WantInputs:  true,
				WantOutputs: true,
			})
		}
	}
}

// --- Early Stream Termination ---

func TestAnthropic_EarlyStreamTermination(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newAnthropicClient(t, tt.TP, t.Name())
	expectedRunName := anthropicRunNameBase + "__" + t.Name()

	stream := client.Messages.NewStreaming(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.Model(testAnthropicModel),
		MaxTokens: 200,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Write a long essay about Go programming")),
		},
	})

	for i := 0; i < 5; i++ {
		if !stream.Next() {
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
			// Per spec: early termination should set error (e.g. "Cancelled") with partial output.
			assertLangSmithRunFields(t, &runs[0], LangSmithRunAssertions{
				WantName:    expectedRunName,
				WantRunType: langsmith.RunQueryResponseRunsRunTypeLlm,
				WantInputs:  true,
				WantOutputs: true,
				ExpectError: true,
			})
		}
	}
}

// --- Error Handling ---

func TestAnthropic_ErrorHandling(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	expectedRunName := anthropicRunNameBase + "__" + t.Name()

	client := anthropic.NewClient(
		option.WithAPIKey("sk-ant-invalid-key-for-testing"),
		option.WithHTTPClient(traceanthropic.Client(traceanthropic.WithTracerProvider(tt.TP), traceanthropic.WithRunNameSuffix(t.Name()))),
	)

	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.Model(testAnthropicModel),
		MaxTokens: 10,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("test")),
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
				WantRunType: langsmith.RunQueryResponseRunsRunTypeLlm,
				ExpectError: true,
			})
		}
	}
}

// --- Tool Use (Non-Streaming) ---

func TestAnthropic_ToolUse_NonStreaming(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newAnthropicClient(t, tt.TP, t.Name())
	expectedRunName := anthropicRunNameBase + "__" + t.Name()

	msg, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.Model(testAnthropicModel),
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("What's the weather in Paris?")),
		},
		Tools: []anthropic.ToolUnionParam{
			{
				OfTool: &anthropic.ToolParam{
					Name:        "get_weather",
					Description: anthropic.String("Get the current weather in a city"),
					InputSchema: anthropic.ToolInputSchemaParam{
						Properties: map[string]interface{}{
							"location": map[string]interface{}{
								"type":        "string",
								"description": "The city name",
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Messages.New with tools: %v", err)
	}

	hasToolUse := false
	for _, block := range msg.Content {
		if block.Type == "tool_use" {
			hasToolUse = true
			if block.Name != "get_weather" {
				t.Errorf("tool name = %q, want 'get_weather'", block.Name)
			}
		}
	}
	if !hasToolUse {
		t.Skip("model didn't use tool, skipping tool assertions")
	}

	tt.TP.ForceFlush(context.Background())
	spans := tt.Spans()

	if v, ok := getSpanAttr(spans, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion")
	} else if !strings.Contains(v, "tool_use") {
		t.Errorf("completion should contain 'tool_use': %s", v)
	}

	if tt.SendsToLangSmith() {
		runs := pollForRuns(t, tt.Project, 1, expectedRunName)
		if len(runs) == 0 {
			t.Error("expected at least one run in LangSmith for tool use")
		} else {
			// Inputs include tools; outputs include tool_calls (behavior 5)
			assertLangSmithRunFields(t, &runs[0], LangSmithRunAssertions{
				WantName:                expectedRunName,
				WantRunType:             langsmith.RunQueryResponseRunsRunTypeLlm,
				WantInputs:              true,
				WantOutputs:             true,
				WantOutputsContainTools: true,
			})
		}
	}
}

// --- Tool Use (Streaming) ---

func TestAnthropic_ToolUse_Streaming(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newAnthropicClient(t, tt.TP, t.Name())
	expectedRunName := anthropicRunNameBase + "__" + t.Name()

	stream := client.Messages.NewStreaming(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.Model(testAnthropicModel),
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("What's the weather in Paris?")),
		},
		Tools: []anthropic.ToolUnionParam{
			{
				OfTool: &anthropic.ToolParam{
					Name:        "get_weather",
					Description: anthropic.String("Get the current weather in a city"),
					InputSchema: anthropic.ToolInputSchemaParam{
						Properties: map[string]interface{}{
							"location": map[string]interface{}{
								"type":        "string",
								"description": "The city name",
							},
						},
					},
				},
			},
		},
	})

	hasToolUseEvent := false
	accumulated := anthropic.Message{}
	for stream.Next() {
		event := stream.Current()
		accumulated.Accumulate(event)
		switch variant := event.AsAny().(type) {
		case anthropic.ContentBlockStartEvent:
			if variant.ContentBlock.Type == "tool_use" {
				hasToolUseEvent = true
			}
		}
	}
	if stream.Err() != nil {
		t.Fatalf("streaming error: %v", stream.Err())
	}
	if !hasToolUseEvent {
		t.Skip("model didn't stream tool use, skipping")
	}

	tt.TP.ForceFlush(context.Background())
	spans := tt.Spans()

	if v, ok := getSpanAttr(spans, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion for streaming tool use")
	} else if !strings.Contains(v, "tool_use") {
		t.Errorf("streaming completion should contain 'tool_use': %s", v)
	}

	if tt.SendsToLangSmith() {
		runs := pollForRuns(t, tt.Project, 1, expectedRunName)
		if len(runs) == 0 {
			t.Error("expected at least one run in LangSmith for streaming tool use")
		} else {
			// Outputs include tool_calls (behavior 5)
			assertLangSmithRunFields(t, &runs[0], LangSmithRunAssertions{
				WantName:                expectedRunName,
				WantRunType:             langsmith.RunQueryResponseRunsRunTypeLlm,
				WantInputs:              true,
				WantOutputs:             true,
				WantOutputsContainTools: true,
			})
		}
	}
}

// --- System Messages ---

func TestAnthropic_SystemMessage(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newAnthropicClient(t, tt.TP, t.Name())
	expectedRunName := anthropicRunNameBase + "__" + t.Name()

	msg, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.Model(testAnthropicModel),
		MaxTokens: 10,
		System: []anthropic.TextBlockParam{
			{Text: "You always respond with exactly one word."},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("What color is the sky?")),
		},
	})
	if err != nil {
		t.Fatalf("Messages.New: %v", err)
	}
	if len(msg.Content) == 0 || msg.Content[0].Text == "" {
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
				WantName:                 expectedRunName,
				WantRunType:              langsmith.RunQueryResponseRunsRunTypeLlm,
				WantInputs:               true,
				WantOutputs:              true,
				WantInputsContainSystem:  true,
			})
		}
	}
}

// --- Multiple Messages ---

func TestAnthropic_MultipleMessages(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newAnthropicClient(t, tt.TP, t.Name())
	expectedRunName := anthropicRunNameBase + "__" + t.Name()

	msg, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.Model(testAnthropicModel),
		MaxTokens: 10,
		System: []anthropic.TextBlockParam{
			{Text: "You are a math tutor. Be concise."},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("What is 2+2?")),
			anthropic.NewAssistantMessage(anthropic.NewTextBlock("4")),
			anthropic.NewUserMessage(anthropic.NewTextBlock("And times 3?")),
		},
	})
	if err != nil {
		t.Fatalf("Messages.New: %v", err)
	}
	if len(msg.Content) == 0 || msg.Content[0].Text == "" {
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
			// Full conversation in run inputs (behavior 7)
			assertLangSmithRunFields(t, &runs[0], LangSmithRunAssertions{
				WantName:                expectedRunName,
				WantRunType:             langsmith.RunQueryResponseRunsRunTypeLlm,
				WantInputs:              true,
				WantOutputs:             true,
				WantInputsMessageCount:  4,
			})
		}
	}
}

// --- Token Usage ---

func TestAnthropic_TokenUsage_NonStreaming(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newAnthropicClient(t, tt.TP, t.Name())
	expectedRunName := anthropicRunNameBase + "__" + t.Name()

	msg, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.Model(testAnthropicModel),
		MaxTokens: 10,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Say hello")),
		},
	})
	if err != nil {
		t.Fatalf("Messages.New: %v", err)
	}
	if msg.Usage.InputTokens == 0 {
		t.Error("expected non-zero input tokens from Anthropic")
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
				WantName:      expectedRunName,
				WantRunType:   langsmith.RunQueryResponseRunsRunTypeLlm,
				WantInputs:    true,
				WantOutputs:   true,
				WantUsage:     true,
			})
		}
	}
}

func TestAnthropic_TokenUsage_Streaming(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newAnthropicClient(t, tt.TP, t.Name())
	expectedRunName := anthropicRunNameBase + "__" + t.Name()

	stream := client.Messages.NewStreaming(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.Model(testAnthropicModel),
		MaxTokens: 10,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Say hello")),
		},
	})

	for stream.Next() {
	}
	if stream.Err() != nil {
		t.Fatalf("streaming error: %v", stream.Err())
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
			// Run has usage (behavior 8)
			assertLangSmithRunFields(t, &runs[0], LangSmithRunAssertions{
				WantName:      expectedRunName,
				WantRunType:   langsmith.RunQueryResponseRunsRunTypeLlm,
				WantInputs:    true,
				WantOutputs:   true,
				WantUsage:     true,
			})
		}
	}
}

// --- Stop Reason ---

func TestAnthropic_StopReason(t *testing.T) {
	tt := newTracedTP(t, "")
	defer tt.Shutdown(context.Background())
	client := newAnthropicClient(t, tt.TP, t.Name())
	expectedRunName := anthropicRunNameBase + "__" + t.Name()

	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.Model(testAnthropicModel),
		MaxTokens: 5,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Say hello")),
		},
	})
	if err != nil {
		t.Fatalf("Messages.New: %v", err)
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
				WantRunType: langsmith.RunQueryResponseRunsRunTypeLlm,
				WantInputs:  true,
				WantOutputs: true,
			})
		}
	}
}
