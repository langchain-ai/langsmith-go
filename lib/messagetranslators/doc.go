// Package messagetranslators provides provider-neutral translators between
// message-oriented wire formats, currently Anthropic Messages and OpenAI
// Responses and Chat Completions.
//
// The package is independent of provider SDKs. Its JSON functions operate on
// complete request or response bodies, while its stream converters operate on
// complete SSE events (they do not parse or frame an HTTP event stream).
//
// The mapping covers text, images, function tools and calls, token limits,
// sampling controls, tool choice, completion state, and token usage. It is
// necessarily lossy where the protocols differ: Responses annotations and
// reasoning-token counts have no Anthropic equivalent; Anthropic cache-creation
// tokens, stop sequences, tool-result is_error, and structured non-text tool
// results have no exact Responses equivalent. Semantically significant known
// request content that cannot be represented is rejected with ErrUnsupported
// rather than silently discarded. The v0 package does not translate audio, files,
// citations/refusals, extended-thinking or reasoning blocks, prompt-cache
// directives, or provider-specific conversation continuation controls. Tool
// results are limited to textual output. Source cache_control hints are accepted
// and omitted when the destination has no equivalent, so cache behavior and
// billing characteristics are not preserved. The converters reject tool-result
// errors, strict tools, parallel-tool controls, structured output, and unknown
// tool types rather than silently dropping them.
//
// A Responses stream normally reports input usage only in its terminal event,
// after the Anthropic message_start has already been emitted. Incremental
// conversion cannot amend that start event: message_start therefore contains
// converted input usage when response.created supplies it and zero otherwise;
// the terminal Anthropic message_delta contains output usage only. Consequently
// Responses-to-Anthropic streaming does not claim exact input usage when it was
// unavailable at response.created. Anthropic-to-Responses terminal events do
// include input_tokens, output_tokens, and total_tokens.
//
// Chat Completions support uses the v0 function-tool schema. Legacy
// functions/function_call, strict and parallel-tool controls, structured output,
// audio, files, reasoning, refusals, logprobs, prediction, and provider-specific
// persistence or service controls are rejected. Chat Completions requests must
// use n=1, and tool calls must be paired with textual role=tool messages. Because
// Chat Completions places assistant text before its tool_calls array, Anthropic
// text following a tool_use in the same completed message is rejected rather
// than reordered. Generated Chat Completions streams always include a terminal
// usage chunk and [DONE]; incoming Chat Completions streams
// accept an optional terminal usage-only chunk and require [DONE]. Input usage
// learned only in that terminal chunk cannot amend an already-emitted Anthropic
// message_start, matching the Responses streaming limitation described above.
//
// Provider stream error codes and messages are forwarded between wire formats;
// this low-level package does not redact them. A gateway that requires a
// different disclosure policy must apply it at its HTTP integration boundary.
//
// Option-bearing conversions and stream constructors can synchronously report
// unknown source fields through ConversionOptions.WarningHandler. The zero value
// keeps the historical silent tolerance. Warning callbacks must not panic; a
// warning may already have been delivered when later validation rejects an event,
// while failed stream events themselves do not advance converter state. See
// README.md for the tolerance table and inspection boundaries.
package messagetranslators
