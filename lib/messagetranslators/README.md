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
AnthropicRequestToResponses(body, model)
ResponsesRequestToAnthropic(body, model)
AnthropicRequestToChatCompletions(body, model)
ChatCompletionsRequestToAnthropic(body, model)
```

### Completed responses

```go
ResponsesResponseToAnthropic(body, model)
AnthropicResponseToResponses(body, model)
ChatCompletionsResponseToAnthropic(body, model)
AnthropicResponseToChatCompletions(body, model)
```

### Streaming responses

```go
NewResponsesToAnthropicStream(model)
NewAnthropicToResponsesStream(model)
NewChatCompletionsToAnthropicStream(model)
NewAnthropicToChatCompletionsStream(model)
```

Each stream converter accepts complete events through:

```go
Convert(event SSEEvent) ([]SSEEvent, error)
Finish() error
```

An empty destination model preserves the source model. See `doc.go` for the
supported semantic mappings and v0 limitations.

## Conversion warnings

Message translators are intentionally tolerant of additive source-schema drift.
The optional warning API makes that tolerance observable without adding
package-global state or a package logger.

```go
collector := &messagetranslators.WarningCollector{}
out, err := messagetranslators.ResponsesRequestToAnthropicWithOptions(
    body,
    model,
    messagetranslators.ConversionOptions{
        WarningHandler: collector.HandleWarning,
    },
)
warnings := collector.Warnings() // a concurrency-safe copy
```

The basic conversion functions and stream constructors use zero options.
Consequently, unknown fields remain silent unless a warning handler is installed.

### Tolerance and error policy

| Source condition | Result |
| --- | --- |
| Missing required field | `ErrInvalidWireData` |
| Missing optional field or documented default | Accepted without warning |
| Source `cache_control` hint without a destination equivalent | Accepted and omitted; destination cache behavior and billing may differ |
| Known, semantically significant unsupported field | `ErrUnsupported`; never downgraded to a warning |
| Unknown extra field in a recognized inspected object | `WarningUnknownField`; conversion continues |
| Unknown discriminator or SSE event type | `ErrUnsupported`, or `ErrInvalidWireData` for an invalid discriminator |
| Malformed JSON, type, range, token count, or function arguments | `ErrInvalidWireData` |
| Invalid stream ordering | `ErrInvalidSequence` |
| Validated mapping that intentionally discards information | May emit `WarningLossyConversion` |

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
options := messagetranslators.ConversionOptions{
    WarningHandler: func(w messagetranslators.Warning) {
        slog.Warn("message conversion schema drift",
            "code", w.Code,
            "path", w.Path,
            "field", w.Field,
            "message", w.Message,
        )
        conversionWarnings.WithLabelValues(string(w.Code), w.Field).Inc()
    },
}
stream := messagetranslators.NewResponsesToAnthropicStreamWithOptions(model, options)
```

Avoid placing unbounded paths or source values in metric labels. `Warning`
intentionally does not contain the unknown value.

## Provider golden captures

The dependency-free capture command records raw completed and streaming fixtures
from Anthropic Messages, OpenAI Responses, and OpenAI Chat Completions:

```bash
go run ./lib/messagetranslators/cmd/capturegoldens -dry-run
go run ./lib/messagetranslators/cmd/capturegoldens
```

See [`cmd/capturegoldens/README.md`](cmd/capturegoldens/README.md) for credentials,
selectors, output files, and cost warnings.
