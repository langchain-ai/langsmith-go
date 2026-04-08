# langsmith-proxy

A reverse proxy that transparently adds [LangSmith](https://smith.langchain.com) tracing to Anthropic and OpenAI API calls.

## Install

```bash
go install github.com/langchain-ai/langsmith-go/cmd/langsmith-proxy@latest
```

## Quick start

### Wrap a command

The simplest way to use the proxy — wrap any command and it automatically sets `ANTHROPIC_BASE_URL` and `OPENAI_BASE_URL`:

```bash
# Wrap Claude Code
langsmith-proxy claude

# Wrap a Python script
langsmith-proxy python my_agent.py

# Wrap Codex
langsmith-proxy codex exec "hello"
```

The proxy starts on a random port, runs the command with the SDK env vars pointing at it, and shuts down when the command exits. All LLM calls are traced to LangSmith.

### Standalone server

Run the proxy as a long-lived server:

```bash
langsmith-proxy
```

Then point any SDK at it:

```bash
export ANTHROPIC_BASE_URL=http://localhost:14355
export OPENAI_BASE_URL=http://localhost:14355/v1
python my_agent.py
```

## Configuration

### Environment variables

The proxy reads `PROXY_`-prefixed env vars first, falling back to the unprefixed version. This lets the proxy have its own credentials when the application being proxied also uses `LANGSMITH_API_KEY`.

| Variable | Default | Description |
|---|---|---|
| `PROXY_LANGSMITH_API_KEY` / `LANGSMITH_API_KEY` | _(none)_ | LangSmith API key. If unset, traces are printed to stdout. |
| `PROXY_LANGSMITH_PROJECT` / `LANGSMITH_PROJECT` | `default` | Project name. Supports Go templates (see below). |
| `PROXY_LANGSMITH_ENDPOINT` / `LANGSMITH_ENDPOINT` | `api.smith.langchain.com` | LangSmith API endpoint. |

### Flags

```
--addr ADDRESS              Listen address (default "localhost:14355", or ":0" in wrap mode)
--anthropic-upstream URL    Anthropic upstream (default: https://api.anthropic.com)
--openai-upstream URL       OpenAI upstream (default: https://api.openai.com)
```

### Project name templates

The project name supports Go `text/template` syntax with these fields:

| Field | Description |
|---|---|
| `{{.User}}` | Current OS username |
| `{{.Dir}}` | Basename of the working directory |
| `{{.Cmd}}` | Wrapped command name (empty in standalone mode) |

Examples:

```bash
# Static
LANGSMITH_PROJECT=my-traces langsmith-proxy claude

# Per-user
LANGSMITH_PROJECT='proxy-{{.User}}' langsmith-proxy claude

# Per-command
LANGSMITH_PROJECT='{{.Cmd}}-traces' langsmith-proxy python my_agent.py
# → project name: "python-traces"

# Combined
LANGSMITH_PROJECT='{{.User}}/{{.Dir}}/{{.Cmd}}' langsmith-proxy claude
# → project name: "ramon/my-project/claude"
```

## How it works

The proxy uses `httputil.ReverseProxy` with the existing `traceanthropic` and `traceopenai` OpenTelemetry instrumentation as its HTTP transport. Requests are routed by path:

| Path | Provider | Traced |
|---|---|---|
| `/v1/messages` | Anthropic | Yes |
| `/v1/chat/completions` | OpenAI | Yes |
| `/v1/responses` | OpenAI | Yes |
| `/v1/completions` | OpenAI | Yes |
| `/v1/embeddings` | OpenAI | Yes |
| `/v1/models` | OpenAI | No (passthrough) |

Streaming SSE responses are passed through without buffering — clients see tokens as they arrive. The tracing instrumentation captures the full request/response in the background via a `BufferedReader` that fires on stream EOF.

## Testing against Eliza

The repo includes a mock LLM server backed by a classic ELIZA psychotherapist:

```bash
# Terminal 1: start Eliza
go run ./cmd/eliza --addr :8080

# Terminal 2: proxy through Eliza
langsmith-proxy --anthropic-upstream=http://localhost:8080 --openai-upstream=http://localhost:8080

# Terminal 3: talk to Eliza through the proxy
curl http://localhost:14355/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"gpt-4o","messages":[{"role":"user","content":"I am sad"}]}'
```

Or in wrap mode:

```bash
langsmith-proxy --anthropic-upstream=http://localhost:8080 --openai-upstream=http://localhost:8080 -- python -c "
from openai import OpenAI
client = OpenAI()  # picks up OPENAI_BASE_URL automatically
print(client.chat.completions.create(model='gpt-4o', messages=[{'role':'user','content':'hello'}]).choices[0].message.content)
"
```
