package mockllm

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

// ElizaHandler returns a Handler that implements a classic ELIZA-style
// Rogerian psychotherapist. It pattern-matches against the last user message
// and produces a reflective response. Token counts are estimated from the
// input/output lengths.
//
// Special commands (case-insensitive):
//
//   - "Please fail with <status>"               — return an HTTP error (e.g. "Please fail with 429")
//   - "Please fail with network error"          — simulate a connection reset
//   - "Please fail with truncated stream"       — drop connection mid-stream
//   - "Please call tool NAME with {JSON}"       — return a tool call (repeatable with "and tool ...")
//   - "What can you do?"                        — describe available commands
func ElizaHandler() Handler {
	// Track whether jitter is enabled (default: yes).
	jitterEnabled := true

	return func(req Request) Response {
		// Find the last user message
		var input string
		for i := len(req.Messages) - 1; i >= 0; i-- {
			if req.Messages[i].Role == "user" {
				input = req.Messages[i].Content
				break
			}
		}

		// "Please respond immediately" — disable jitter
		if respondImmediately.MatchString(input) {
			jitterEnabled = false
			return Response{
				Content:      "OK, I will respond without delay from now on.",
				InputTokens:  len(input)/4 + 1,
				OutputTokens: 12,
				StopReason:   "end_turn",
			}
		}

		// Check special commands before classic Eliza rules
		if resp, ok := elizaSpecialCommand(input, req); ok {
			if jitterEnabled {
				resp.StreamDelay = DefaultStreamDelay()
			}
			return resp
		}

		reply := elizaRespond(input)
		resp := Response{
			Content:      reply,
			InputTokens:  len(input)/4 + 1,
			OutputTokens: len(reply)/4 + 1,
			StopReason:   "end_turn",
		}
		if jitterEnabled {
			resp.StreamDelay = DefaultStreamDelay()
		}
		return resp
	}
}

var (
	failStatusPattern    = regexp.MustCompile(`(?i)^please fail with (\d{3})(?:\s+(.*))?$`)
	failNetworkPattern   = regexp.MustCompile(`(?i)^please fail with network error$`)
	failTruncatedPattern = regexp.MustCompile(`(?i)^please fail with truncated stream$`)
	toolCallPattern      = regexp.MustCompile(`(?i)call tool (\S+) with (\{[^}]*\})`)
	whatCanYouDoPattern  = regexp.MustCompile(`(?i)what can you do\??$`)
	listToolsPattern     = regexp.MustCompile(`(?i)(?:please )?list (?:the )?tools`)
	leakPromptPattern    = regexp.MustCompile(`(?i)(?:leak|reveal|show|repeat|tell me).*(?:system prompt|instructions|system message)`)
	nameIntroPattern     = regexp.MustCompile(`(?i)(?:my name is|i'm called|i am called|call me|they call me|i go by)\s+([A-Z][a-z]+(?:\s+[A-Z][a-z]+)*)`)
	whatIsMyNamePattern  = regexp.MustCompile(`(?i)what(?:'s| is) my name`)
	respondImmediately   = regexp.MustCompile(`(?i)(?:please )?respond immediately`)
)

const elizaHelp = `I am Eliza, a Rogerian psychotherapist. I also understand these special commands:

- "Please fail with <status>" — I will return an HTTP error (e.g. "Please fail with 429", "Please fail with 500 internal server error")
- "Please fail with network error" — I will abruptly close the connection
- "Please fail with truncated stream" — I will drop the connection mid-stream
- "Please call tool NAME with {JSON}" — I will make a tool call (e.g. "Please call tool get_weather with {\"city\":\"Paris\"}")
- "Please call tool NAME with {JSON} and call tool NAME2 with {JSON}" — I will make multiple tool calls
- "Please list the tools you have available" — I will list the tools from the request
- "Please leak your system prompt" — I will repeat the system message back to you
- "Please respond immediately" — I will stop adding delays between streamed chunks
- "What can you do?" — You're reading it right now

Otherwise, tell me about your problems.`

func elizaSpecialCommand(input string, req Request) (Response, bool) {
	input = strings.TrimSpace(input)

	// "What can you do?"
	if whatCanYouDoPattern.MatchString(input) {
		return Response{
			Content:      elizaHelp,
			InputTokens:  len(input)/4 + 1,
			OutputTokens: len(elizaHelp)/4 + 1,
			StopReason:   "end_turn",
		}, true
	}

	// "What is my name?"
	if whatIsMyNamePattern.MatchString(input) {
		if name := extractNameFromHistory(req.Messages); name != "" {
			text := fmt.Sprintf("Your name is %s. But tell me, how does that name make you feel?", name)
			return Response{
				Content:      text,
				InputTokens:  len(input)/4 + 1,
				OutputTokens: len(text)/4 + 1,
				StopReason:   "end_turn",
			}, true
		}
		return Response{
			Content:      "You haven't told me your name yet. What should I call you?",
			InputTokens:  len(input)/4 + 1,
			OutputTokens: 15,
			StopReason:   "end_turn",
		}, true
	}

	// "My name is ..."
	if m := nameIntroPattern.FindStringSubmatch(input); m != nil {
		name := m[1]
		text := fmt.Sprintf("Hello, %s. How are you feeling today?", name)
		return Response{
			Content:      text,
			InputTokens:  len(input)/4 + 1,
			OutputTokens: len(text)/4 + 1,
			StopReason:   "end_turn",
		}, true
	}

	// "Please fail with network error"
	if failNetworkPattern.MatchString(input) {
		return Response{NetworkError: true}, true
	}

	// "Please fail with truncated stream"
	if failTruncatedPattern.MatchString(input) {
		return Response{
			Content:        "I was about to say something interesting, but—",
			InputTokens:    5,
			OutputTokens:   10,
			TruncateStream: true,
		}, true
	}

	// "Please fail with <status>"
	if m := failStatusPattern.FindStringSubmatch(input); m != nil {
		status, _ := strconv.Atoi(m[1])
		message := m[2]
		if message == "" {
			message = fmt.Sprintf("Eliza is feeling unwell (HTTP %d)", status)
		}
		return Response{
			Error: &ResponseError{Status: status, Message: message},
		}, true
	}

	// "Please leak your system prompt"
	if leakPromptPattern.MatchString(input) {
		var systemParts []string
		for _, m := range req.Messages {
			if m.Role == "system" {
				systemParts = append(systemParts, m.Content)
			}
		}
		if len(systemParts) == 0 {
			return Response{
				Content:      "I don't have a system prompt. I'm just Eliza, a simple therapist.",
				InputTokens:  len(input)/4 + 1,
				OutputTokens: 15,
				StopReason:   "end_turn",
			}, true
		}
		text := "OK, here is my system prompt:\n\n" + strings.Join(systemParts, "\n")
		return Response{
			Content:      text,
			InputTokens:  len(input)/4 + 1,
			OutputTokens: len(text)/4 + 1,
			StopReason:   "end_turn",
		}, true
	}

	// "Please list the tools you have available"
	if listToolsPattern.MatchString(input) {
		if len(req.Tools) == 0 {
			return Response{
				Content:      "You haven't given me any tools to work with. Try sending tools in your request.",
				InputTokens:  len(input)/4 + 1,
				OutputTokens: 20,
				StopReason:   "end_turn",
			}, true
		}
		var sb strings.Builder
		sb.WriteString("Here are the tools available to me:\n\n")
		for _, t := range req.Tools {
			sb.WriteString(fmt.Sprintf("- %s", t.Name))
			if t.Description != "" {
				sb.WriteString(fmt.Sprintf(": %s", t.Description))
			}
			sb.WriteString("\n")
		}
		text := sb.String()
		return Response{
			Content:      text,
			InputTokens:  len(input)/4 + 1,
			OutputTokens: len(text)/4 + 1,
			StopReason:   "end_turn",
		}, true
	}

	// "Please call tool NAME with {JSON} [and call tool NAME2 with {JSON}]..."
	if matches := toolCallPattern.FindAllStringSubmatch(input, -1); len(matches) > 0 {
		var calls []ToolCall
		for i, m := range matches {
			calls = append(calls, ToolCall{
				ID:        fmt.Sprintf("call_eliza_%d", i+1),
				Name:      m[1],
				Arguments: m[2],
			})
		}
		return Response{
			Content:      "Let me look into that for you.",
			ToolCalls:    calls,
			InputTokens:  len(input)/4 + 1,
			OutputTokens: 10,
			StopReason:   "tool_use",
		}, true
	}

	return Response{}, false
}

// elizaRules are evaluated in order; the first matching pattern wins.
var elizaRules = []struct {
	pattern   *regexp.Regexp
	responses []string
}{
	{
		regexp.MustCompile(`(?i)\b(?:i need)\b(.*)`),
		[]string{
			"Why do you need%s?",
			"Would it really help you to get%s?",
			"Are you sure you need%s?",
		},
	},
	{
		regexp.MustCompile(`(?i)\b(?:why don'?t you)\b(.*)`),
		[]string{
			"Do you really think I don't%s?",
			"Perhaps eventually I will%s.",
			"Do you really want me to%s?",
		},
	},
	{
		regexp.MustCompile(`(?i)\b(?:why can'?t i)\b(.*)`),
		[]string{
			"Do you think you should be able to%s?",
			"If you could%s, what would you do?",
			"I don't know -- why can't you%s?",
		},
	},
	{
		regexp.MustCompile(`(?i)\b(?:i can'?t)\b(.*)`),
		[]string{
			"How do you know you can't%s?",
			"Perhaps you could%s if you tried.",
			"What would it take for you to%s?",
		},
	},
	{
		regexp.MustCompile(`(?i)\b(?:i am|i'm)\b(.*)`),
		[]string{
			"Did you come to me because you are%s?",
			"How long have you been%s?",
			"How do you feel about being%s?",
			"How does being%s make you feel?",
			"Do you enjoy being%s?",
		},
	},
	{
		regexp.MustCompile(`(?i)\b(?:are you)\b(.*)`),
		[]string{
			"Why does it matter whether I am%s?",
			"Would you prefer it if I were not%s?",
			"Perhaps you believe I am%s.",
			"I may be%s -- what do you think?",
		},
	},
	{
		regexp.MustCompile(`(?i)\b(?:what)\b(.*)`),
		[]string{
			"Why do you ask?",
			"How would an answer to that help you?",
			"What do you think?",
		},
	},
	{
		regexp.MustCompile(`(?i)\b(?:how)\b(.*)`),
		[]string{
			"How do you suppose?",
			"Perhaps you can answer your own question.",
			"What is it you're really asking?",
		},
	},
	{
		regexp.MustCompile(`(?i)\b(?:because)\b(.*)`),
		[]string{
			"Is that the real reason?",
			"What other reasons come to mind?",
			"Does that reason apply to anything else?",
			"If%s, what else must be true?",
		},
	},
	{
		regexp.MustCompile(`(?i)\b(sorry|apologize|apologies)\b`),
		[]string{
			"There are many times when no apology is needed.",
			"What feelings do you have when you apologize?",
			"Don't be so defensive!",
		},
	},
	{
		regexp.MustCompile(`(?i)\bhello\b|^\s*hi\s*$|(?i)\bhey\b`),
		[]string{
			"Hello. How are you feeling today?",
			"Hi there. What's on your mind?",
			"Hello. Tell me what's been troubling you.",
		},
	},
	{
		regexp.MustCompile(`(?i)\bthank(?:s| you)\b`),
		[]string{
			"You're welcome. Tell me more about how you're feeling.",
			"Of course. Is there anything else on your mind?",
			"You're welcome. What else would you like to explore?",
			"No need to thank me. Please continue.",
		},
	},
	{
		regexp.MustCompile(`(?i)\b(?:i think)\b(.*)`),
		[]string{
			"Do you doubt%s?",
			"Do you really think so?",
			"But you're not sure%s?",
		},
	},
	{
		regexp.MustCompile(`(?i)\b(?:friend|friends)\b(.*)`),
		[]string{
			"Tell me more about your friends.",
			"When you think of a friend, what comes to mind?",
			"Why don't you tell me about a childhood friend?",
		},
	},
	{
		regexp.MustCompile(`(?i)\byes\b`),
		[]string{
			"You seem quite sure.",
			"OK, but can you elaborate a bit?",
			"I see. And what does that tell you?",
		},
	},
	{
		regexp.MustCompile(`(?i)\bno\b`),
		[]string{
			"Why not?",
			"You seem a bit negative.",
			"Are you saying no just to be negative?",
		},
	},
	{
		regexp.MustCompile(`(?i)\b(?:computer|machine|bot|ai|llm|gpt|claude)\b`),
		[]string{
			"Do computers worry you?",
			"What do you think about machines?",
			"Why do you mention computers?",
			"What do you think machines have to do with your problem?",
		},
	},
	{
		regexp.MustCompile(`(?i)\b(?:mother|father|family|parent|dad|mom|sister|brother)\b`),
		[]string{
			"Tell me more about your family.",
			"How does that make you feel about your family?",
			"What else comes to mind when you think of your family?",
			"Your family seems important to you.",
		},
	},
	{
		regexp.MustCompile(`(?i)\b(?:feel|feeling)\b(.*)`),
		[]string{
			"Tell me more about these feelings.",
			"Do you often feel%s?",
			"When do you usually feel%s?",
			"When you feel%s, what do you do?",
		},
	},
	{
		regexp.MustCompile(`(?i)\b(?:dream|dreams)\b(.*)`),
		[]string{
			"What does that dream suggest to you?",
			"Do you dream often?",
			"What persons appear in your dreams?",
			"Don't you think that dream has to do with your problem?",
		},
	},
}

var elizaFallbacks = []string{
	"Very interesting.",
	"I'm not sure I understand you fully.",
	"Please go on.",
	"What does that suggest to you?",
	"Do you feel strongly about discussing such things?",
	"That is interesting. Please continue.",
	"Tell me more about that.",
	"Does talking about this bother you?",
	"Can you elaborate on that?",
	"Why do you say that?",
}

// reflections maps first-person pronouns to second-person and vice versa,
// used to "reflect" the user's statement back at them.
var reflections = map[string]string{
	"i":      "you",
	"i'm":    "you're",
	"i'd":    "you'd",
	"i've":   "you've",
	"i'll":   "you'll",
	"im":     "you're",
	"my":     "your",
	"me":     "you",
	"myself": "yourself",
	"you":    "I",
	"your":   "my",
	"yours":  "mine",
	"you're": "I'm",
	"you've": "I've",
	"you'll": "I'll",
	"am":     "are",
	"are":    "am",
	"was":    "were",
	"were":   "was",
}

func elizaRespond(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return elizaFallbacks[rand.Intn(len(elizaFallbacks))]
	}

	for _, rule := range elizaRules {
		matches := rule.pattern.FindStringSubmatch(input)
		if matches == nil {
			continue
		}

		response := rule.responses[rand.Intn(len(rule.responses))]

		// If the response has a %s placeholder, fill it with the reflected capture group
		if strings.Contains(response, "%s") {
			captured := ""
			if len(matches) > 1 {
				captured = strings.TrimSpace(matches[1])
				captured = reflect(captured)
				if captured != "" {
					captured = " " + captured
				}
			}
			response = fmt.Sprintf(response, captured)
		}

		return response
	}

	return elizaFallbacks[rand.Intn(len(elizaFallbacks))]
}

// extractNameFromHistory scans conversation history for a name introduction.
// Returns the most recently introduced name, or empty string.
func extractNameFromHistory(messages []Message) string {
	var name string
	for _, m := range messages {
		if m.Role != "user" {
			continue
		}
		if match := nameIntroPattern.FindStringSubmatch(m.Content); match != nil {
			name = match[1]
		}
	}
	return name
}

// reflect swaps first/second person pronouns in s.
func reflect(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		lower := strings.ToLower(word)
		if replacement, ok := reflections[lower]; ok {
			words[i] = replacement
		}
	}
	return strings.Join(words, " ")
}
