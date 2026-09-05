# OpenAI tracing

Use `Client` or `WrapClient` for HTTP Chat Completions and Responses requests.
For the Responses WebSocket API, dial with `github.com/gorilla/websocket` and
wrap the connected socket:

```go
raw, _, err := websocket.DefaultDialer.DialContext(ctx, endpoint, headers)
if err != nil {
    return err
}
conn := traceopenai.WrapWebSocket(ctx, raw, traceopenai.WithTracerProvider(tp))
defer conn.Close()

if err := conn.WriteJSON(map[string]any{
    "type": "response.create",
    "model": "gpt-6-astra",
    "input": "Hello",
}); err != nil {
    return err
}
for {
    var event map[string]any
    if err := conn.ReadJSON(&event); err != nil {
        return err
    }
    switch event["type"] {
    case "response.completed", "response.incomplete", "response.failed", "error":
        // Handle the response or error in your application.
        return nil
    }
}
```

Use the wrapper's `ReadMessage`/`WriteMessage` or `ReadJSON`/`WriteJSON` methods
for all messages. One reader and one writer may run concurrently. The caller
configures the endpoint, authentication headers, deadlines, and read limits;
the wrapper does not dial or reconnect. Canceling `ctx` closes the socket.

Each response gets its own span, including queued requests, multiplexed
`stream_id` lanes, and automatic steering continuations. Keep reading after
`response.steer.accepted` and the original response's terminal event to trace
the continuation. Spans record prompts, final output, first-token timing, and
usage; no full delta history is buffered. Close, cancellation, and transport
errors end outstanding spans. As with HTTP tracing, prompts and outputs are
exported to the configured tracer provider.

Both transports preserve `configuration_update` reasoning details and map
`cache_write_tokens` to LangSmith's `cache_creation` usage detail, with service
tier and long-context prefixes where applicable.
