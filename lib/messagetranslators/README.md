# Message translators

`messagetranslators` converts raw wire payloads between Anthropic Messages and
OpenAI Responses or Chat Completions. It has no dependency on either provider
SDK.

The JSON functions operate on complete request or response bodies. Stream
converters operate on complete SSE events; callers remain responsible for HTTP
and SSE framing.

## Main entrypoints

### Requests

```go
AnthropicRequestToResponses(body, model, ...Option)
ResponsesRequestToAnthropic(body, model, ...Option)
AnthropicRequestToChatCompletions(body, model, ...Option)
ChatCompletionsRequestToAnthropic(body, model, ...Option)
```

### Completed responses

```go
ResponsesResponseToAnthropic(body, model, ...Option)
AnthropicResponseToResponses(body, model, ...Option)
ChatCompletionsResponseToAnthropic(body, model, ...Option)
AnthropicResponseToChatCompletions(body, model, ...Option)
```

### Streaming responses

```go
NewResponsesToAnthropicStream(model, ...Option)
NewAnthropicToResponsesStream(model, ...Option)
NewChatCompletionsToAnthropicStream(model, ...Option)
NewAnthropicToChatCompletionsStream(model, ...Option)
```

Each stream converter accepts complete events through:

```go
Convert(event SSEEvent) ([]SSEEvent, error)
Finish() error
Usage() Usage
```

An empty destination model preserves the source model. See `doc.go` for the
supported semantic mappings and v0 limitations.

## Conversion warnings

Message translators are intentionally tolerant of additive source-schema drift.
The optional warning API makes that tolerance observable without adding
package-global state or a package logger.

```go
collector := &messagetranslators.WarningCollector{}
out, err := messagetranslators.ResponsesRequestToAnthropic(
    body,
    model,
    messagetranslators.WithWarningHandler(collector.HandleWarning),
)
warnings := collector.Warnings() // a concurrency-safe copy
```

Passing no options keeps the silent-tolerance default: source inspection only
runs when a warning handler is installed, so unknown fields stay silent
otherwise.

The other options are `WithClock`, which replaces `time.Now` for the timestamps
these translators must generate (Anthropic payloads carry none, Responses and
Chat Completions require them) and makes output deterministic under test, and
`WithUsageHandler`, described under [Token usage](#token-usage).

### Tolerance and error policy

| Source condition | Result |
| --- | --- |
| Missing required field | `ErrInvalidWireData` |
| Missing optional field | Accepted without warning |
| Field carrying the source provider's documented default, or an explicit `null` | Accepted without warning; see [Documented defaults](#documented-defaults) |
| Known, semantically significant unsupported field | `ErrUnsupported`; never downgraded to a warning |
| Validated mapping that intentionally discards information | Accepted; emits `WarningLossyConversion`. See [Known incompatibilities](#known-incompatibilities) |
| Unknown extra field in a recognized inspected object | `WarningUnknownField`; conversion continues |
| Unknown discriminator or SSE event type | `ErrUnsupported`, or `ErrInvalidWireData` for an invalid discriminator |
| Malformed JSON, type, range, token count, or function arguments | `ErrInvalidWireData` |
| Invalid stream ordering | `ErrInvalidSequence` |

#### Documented defaults

SDKs serialize fields the caller never set. `parallel_tool_calls: true`,
`truncation: "disabled"`, `response_format: {"type":"text"}`,
`text: {"format":{"type":"text"}}`, `logprobs: false`, and zero
frequency/presence penalties all mean "nothing unusual was requested", so they
are accepted rather than rejected. An explicit JSON `null` is treated as an
absent field for the same reason. The same field carrying any other value is
still `ErrUnsupported` — accepting the default never widens into silently
dropping a real request.

#### Terminal reasons

Where a source's advisory terminal reason disagrees with its own content, the
content wins: a completion carrying tool calls maps to a tool-call turn whatever
its `finish_reason` says. OpenAI reports `tool_calls`, but several
OpenAI-compatible servers report `stop` beside a populated `tool_calls` array.
`content_filter` and Anthropic's `refusal` map to each other.

Warnings are schema-drift telemetry, **not** a guarantee that an unknown field is
semantically harmless. A production gateway should monitor new field names and
paths before deciding whether a translator update or fail-closed policy is
needed.

### Warning coverage

Inspection is source- and API-context-specific. It currently covers:

- Anthropic Messages request and completed-response envelopes, system blocks,
  messages, content blocks, image sources, tool use/results, tools, tool choices,
  and usage;
- OpenAI Responses request and completed-response envelopes, input/output items,
  content parts, function calls/results, tools, tool choices, incomplete details,
  errors, usage, and token-detail objects;
- OpenAI Chat Completions request and completed-response envelopes, messages,
  content parts, image URL objects, tools, tool choices, choices, deltas, usage,
  and usage-detail objects;
- recognized Anthropic, Responses, and Chat Completions stream event objects.

Arbitrary user-defined JSON is deliberately opaque and is not inspected. This
includes JSON Schema bodies, function argument payloads, metadata values, and
arbitrary tool output values. Objects behind unknown discriminators are not
inspected because their schema is unknown. Destination-generated objects are
never inspected.

Known source fields that are accepted, validated, ignored as harmless defaults,
or rejected as unsupported are included in the relevant allowlist and do not
produce an unknown-field warning.

### Callback and stream semantics

Handlers run synchronously in deterministic field order before conversion. They
must return normally and should be fast. Handler panics are intentionally not
recovered.

A warning can be delivered even if later validation rejects the same body or
stream event. Failed stream events roll back converter state, but warnings already
delivered are observations and are not rolled back.

A Gateway can route warnings to logs and metrics without making this package
depend on a logging implementation:

```go
warn := messagetranslators.WithWarningHandler(func(w messagetranslators.Warning) {
    slog.Warn("message conversion schema drift",
        "code", w.Code,
        "path", w.Path,
        "field", w.Field,
        "message", w.Message,
    )
    conversionWarnings.WithLabelValues(string(w.Code), w.Field).Inc()
})
stream := messagetranslators.NewResponsesToAnthropicStream(model, warn)
```

Avoid placing unbounded paths or source values in metric labels. `Warning`
intentionally does not contain the unknown value.

## Known incompatibilities

Two policies govern anything that cannot be represented at the destination:

- **Throws.** Anything that changes what the model is asked to do, or what it
  actually did, is rejected with `ErrUnsupported`. The caller gets an error
  instead of a response that quietly means something different from the request.
- **Best effort.** Values that carry no instruction — provider bookkeeping,
  billing detail, cache hints — are dropped, and the conversion succeeds. These
  emit `WarningLossyConversion` where the loss is worth measuring.

The tables below are exhaustive for v0. Where a row applies to only one
destination, the destination is named.

### Anthropic source

| Source construct | Result |
| --- | --- |
| `thinking` / `redacted_thinking` content blocks | **Throws** |
| `document`, `search_result`, or any unrecognized content block | **Throws** |
| `image` with a `file` source | **Throws** |
| `tool_result` with `is_error: true` | **Throws** |
| `tool_result` containing non-text blocks (images, structured output) | **Throws** |
| `thinking`, `top_k`, `output_config`, `container` request fields | **Throws** |
| `mcp_servers`, `context_management` request fields | **Throws** |
| `metadata` request field, unless `{}` | **Throws** |
| `service_tier`, unless `"auto"` | **Throws** |
| `stop_sequences` → Responses | **Throws** (Responses has no stop control) |
| More than 4 `stop_sequences` → Chat Completions | **Throws** (the format caps at 4) |
| Text following a `tool_use` block in one message → Chat Completions | **Throws** (see [Ordering](#ordering-and-pairing)) |
| `pause_turn` stop reason | **Throws** |
| `refusal` stop reason → Responses | **Throws** |
| `refusal` stop reason → Chat Completions | Best effort: mapped to `content_filter` |
| `cache_control` hints anywhere | Best effort: dropped. Destination cache behavior and billing will differ |
| `cache_creation_input_tokens` | Best effort: folded into the OpenAI input total, losing its separately-billed identity |
| `stop_sequence` value on a completed response | Best effort: dropped; neither OpenAI format reports which sequence matched |
| Multiple `system` text blocks | Best effort: concatenated with no separator into one string |
| Multi-block `tool_result` text | Best effort: concatenated with no separator |
| `container` on a completed response | Best effort: dropped |

### Responses source

| Source construct | Result |
| --- | --- |
| `reasoning` items, and `reasoning`/`reasoning_summary_*` stream events | **Throws** |
| `web_search_call`, `file_search_call`, `computer_call`, `image_generation_call`, or any unrecognized item | **Throws** |
| `refusal` content parts and `response.refusal.*` events | **Throws** |
| `input_file` content parts | **Throws** |
| `input_image` with `detail` | **Throws** (no Anthropic equivalent for the fidelity hint) |
| Non-empty `annotations`, and `response.output_text.annotation.*` events | **Throws** |
| `function_call_output` whose `output` is not a string | **Throws** |
| `strict` on a tool | **Throws** |
| `system` / `developer` role input messages | **Throws** — use `instructions` |
| `previous_response_id`, `conversation`, `prompt_cache_key`, `prompt_cache_retention` | **Throws** |
| `background: true`, `store: false`, non-empty `include` | **Throws** |
| `truncation`, unless `"disabled"` | **Throws** |
| `text`, unless it requests plain text | **Throws** (structured output) |
| `parallel_tool_calls: false` | **Throws** |
| `max_tool_calls`, `user`, `safety_identifier`, non-`{}` `metadata` | **Throws** |
| `top_logprobs`, unless `0` | **Throws** |
| `service_tier`, unless `"auto"` or `"default"` | **Throws** |
| `incomplete_details.reason` other than `max_output_tokens` | **Throws** |
| A truncated function call in an `incomplete` response | **Throws** — the partial arguments are not valid JSON and cannot be represented |
| `reasoning_tokens` usage | Best effort: validated, then dropped. Anthropic has no reasoning usage category |
| Empty `annotations` / `logprobs` arrays | Best effort: dropped as SDK defaults |
| `obfuscation` on stream deltas | Best effort: ignored |

### Chat Completions source

| Source construct | Result |
| --- | --- |
| `input_audio`, `file`, and `refusal` content parts | **Throws** |
| `image_url.detail` | **Throws** |
| `audio` or non-empty `refusal` on a message | **Throws** |
| Non-empty `annotations` on a message | **Throws** |
| Legacy `functions` / `function_call` | **Throws** |
| `strict` on a tool | **Throws** |
| `n` other than `1`, or a response with more than one choice | **Throws** |
| `logprobs` on a response choice | **Throws** |
| `system` / `developer` messages after the conversation has begun | **Throws** |
| `name` on a message | **Throws** (Anthropic has no per-message name) |
| `response_format`, unless `{"type":"text"}` | **Throws** (structured output) |
| `parallel_tool_calls: false`, `logprobs: true`, `store: true` | **Throws** |
| Non-zero `frequency_penalty` / `presence_penalty` | **Throws** — Anthropic has no equivalent sampling control |
| `logit_bias`, unless `{}` | **Throws** |
| `modalities`, unless `["text"]` | **Throws** |
| `prediction`, `seed`, `user`, `reasoning_effort`, `web_search_options`, `safety_identifier` | **Throws** |
| `verbosity`, unless `"medium"`; `service_tier`, unless `"auto"`/`"default"` | **Throws** |
| Non-zero audio or prediction token counts in usage | **Throws** |
| `service_tier` on a response, or on a stream chunk | **Throws** |
| `content_filter` finish reason | Best effort: mapped to Anthropic `refusal` |
| `reasoning_tokens` usage | Best effort: validated, then dropped |
| `system_fingerprint` on a completed response | Best effort: dropped |
| `created` timestamp | Best effort: dropped; Anthropic responses carry no creation time |

### Ordering and pairing

These are structural, not field-level, and throw `ErrInvalidSequence` unless
noted:

| Constraint | Why |
| --- | --- |
| The conversation must begin with a `user` message | Anthropic requires it |
| Every tool call must be followed by its result before the next assistant turn | Anthropic and Chat Completions both require the pairing |
| Tool-call IDs must be unique and must resolve | Both formats key results by ID |
| Anthropic user content may not follow an unresolved tool call → Chat Completions | `role=tool` messages must come first (`ErrUnsupported`) |
| Anthropic text after `tool_use` in one message → Chat Completions | Chat Completions puts assistant text before `tool_calls`, so the order cannot be preserved and is not silently rearranged (`ErrUnsupported`). The same input converts cleanly to Responses |

A source's terminal reason is *not* cross-checked against its content — see
[Terminal reasons](#terminal-reasons).

### Generated, not translated

These destination fields have no source to translate from and are synthesized.
Do not treat them as provider data:

| Field | Value |
| --- | --- |
| Response and message IDs | Derived by hashing the source ID, so they are stable per source but are **not** the upstream provider's IDs. Correlate on your own request ID |
| `created` / `created_at` / `completed_at` | From the clock — see `WithClock` |
| Responses envelope fields (`temperature`, `top_p`, `parallel_tool_calls`, `service_tier`, `store`, `truncation`, `text`, `tool_choice`, `tools`, `top_logprobs`) | Schema-required placeholders holding API defaults, not the values the request actually used |
| `sequence_number` on generated Responses events | Renumbered from zero |
| Chat Completions terminal usage chunk | Always emitted, whether or not the client asked for `stream_options.include_usage` |

### Not covered at all

There is no direct Responses ↔ Chat Completions translator; only the four
Anthropic-paired directions exist. Routing between the two OpenAI formats means
pivoting through Anthropic, which applies Anthropic's stricter rules (see
[Ordering](#ordering-and-pairing)) to a pair of formats that would not otherwise
require them.

## Token usage

`Usage` reports token accounting in Anthropic terms: `InputTokens` excludes
cache reads, which are counted in `CacheReadInputTokens`.

Both OpenAI stream formats disclose input token counts only in a terminal event,
by which time the Anthropic `message_start` that would have carried them is long
gone. The emitted `message_delta` therefore reports output usage only. That is a
wire-format limitation, not a loss of information: the converter saw the counts
and exposes them.

```go
stream := messagetranslators.NewChatCompletionsToAnthropicStream(model)
// ... feed events ...
usage := stream.Usage() // input counts the translated stream could not carry
```

Use `WithUsageHandler` to be called once when a conversion or stream reaches a
terminal state, instead of polling.

## Errors and HTTP status

`ConversionError.Path` carries a JSON path (`$.messages[0].content`) or stream
event location (`event response.output_item.added`). Callers choose the status,
because the same sentinel means different things depending on which side of the
proxy produced the payload:

| Sentinel | Converting a request | Converting a response or stream |
| --- | --- | --- |
| `ErrInvalidWireData` | `400` — the client sent a malformed body | `502` — the upstream provider did |
| `ErrUnsupported` | `400` — the client asked for something unrepresentable | `502` — the provider replied with something unrepresentable |
| `ErrInvalidSequence` | `400` — malformed conversation history | `502` — the provider emitted events out of order |
| `ErrTruncatedStream` | n/a | `502` — the provider hung up before a terminal event |

Errors are also reachable **after** the response has been committed. A stream
converter can reject event N of a response whose first N-1 events have already
been flushed with a `200`, so there is no status left to change: an integration
needs a deliberate path that emits a terminal error event in the destination
format and closes the stream. This is reachable for every `ErrInvalidSequence`
case, for malformed tool-call arguments detected at the closing event, and for
Anthropic text following a `tool_use` block when the destination is Chat
Completions.

Provider error codes and messages are forwarded between wire formats; this
package does not redact them. Apply any disclosure policy at the HTTP boundary.

## Performance notes

Per-event cost is dominated by JSON decoding and re-encoding, which the
`map[string]any` representation makes unavoidable. Accumulated text and
tool-call arguments grow by append rather than string concatenation, so a stream
costs time linear in its length rather than quadratic.

Stream converters snapshot their state before each event so that a rejected
event leaves nothing partially applied. That costs roughly 14% of per-event time
and 8% of allocated bytes; it has been left in place because the rollback
guarantee is worth more than the margin.

## Provider golden captures

The dependency-free capture command records raw completed and streaming fixtures
from Anthropic Messages, OpenAI Responses, and OpenAI Chat Completions:

```bash
go run ./lib/messagetranslators/cmd/capturegoldens -dry-run
go run ./lib/messagetranslators/cmd/capturegoldens
```

See [`cmd/capturegoldens/README.md`](cmd/capturegoldens/README.md) for credentials,
selectors, output files, and cost warnings.
