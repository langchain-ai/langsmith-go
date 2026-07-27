# Provider golden capture

`capturegoldens` is a dependency-free, raw `net/http` utility for recording the
provider wire formats used by `lib/messagetranslators`. It does not use either
provider SDK. It is intended for deliberate, manual fixture refreshes—not for
normal test runs.

## Credentials, models, and endpoints

Set the credentials needed by the selected provider:

```sh
export ANTHROPIC_API_KEY='...'
export OPENAI_API_KEY='...'
```

Credentials are sent only in HTTP headers. They are never included in a request
snapshot or written to disk. Avoid putting credentials in base URLs; the command
rejects URLs containing user information.

The defaults are deliberately inexpensive models:

- `ANTHROPIC_MODEL=claude-haiku-4-5`
- `OPENAI_MODEL=gpt-5.6-luna`

Set both model variables explicitly when reproducibility matters. Provider aliases
such as `latest` can change. `ANTHROPIC_VERSION` defaults to the stable
`2023-06-01`; it can also be set with `-anthropic-version`.

`ANTHROPIC_BASE_URL` and `OPENAI_BASE_URL` override the provider roots for a
proxy or local `httptest`-style server. Both `https://host` and
`https://host/v1` forms are accepted. Defaults are the public provider hosts.

## Commands

Run from the repository root:

```sh
# Inspect all request bodies; no credentials, network, or files are needed.
go run ./lib/messagetranslators/cmd/capturegoldens -dry-run

# Capture all 12 fixtures (3 APIs x 2 scenarios x 2 modes).
go run ./lib/messagetranslators/cmd/capturegoldens

# Capture only streaming tool calls from OpenAI Responses.
go run ./lib/messagetranslators/cmd/capturegoldens \
  -provider openai -api responses -scenario tool -mode stream

# Deliberately replace selected existing fixtures.
go run ./lib/messagetranslators/cmd/capturegoldens \
  -provider anthropic -overwrite -timeout 90s
```

Flags:

- `-output` (default `lib/messagetranslators/testdata/provider_goldens`)
- `-provider all|anthropic|openai`
- `-api all|messages|chat-completions|responses`
- `-scenario all|text|tool`
- `-mode all|completed|stream`
- `-overwrite` (false by default)
- `-timeout` (per HTTP request, default `60s`)
- `-dry-run`
- `-anthropic-version`

Before making a network call, the command checks all selected destinations and
refuses to proceed if one exists unless `-overwrite` is set. Writes publish a
complete temporary file atomically. Execution is fail-fast; on an error, fixtures
already reported as written remain in place and the summary says how many were
completed.

## Expected files

The default directory receives these names:

```text
anthropic_messages_completed_text.json
anthropic_messages_completed_tool.json
anthropic_messages_stream_text.json
anthropic_messages_stream_tool.json
openai_chat_completions_completed_text.json
openai_chat_completions_completed_tool.json
openai_chat_completions_stream_text.json
openai_chat_completions_stream_tool.json
openai_responses_completed_text.json
openai_responses_completed_tool.json
openai_responses_stream_text.json
openai_responses_stream_tool.json
```

Each JSON document has `schema_version`, `provider`, `api`, `mode`, `scenario`,
and `request`, plus either the provider's completed `response` object or an
`events` array. SSE data is parsed as JSON where possible; `[DONE]` and other
non-JSON data remain strings, and an SSE `event` name is retained when present.
Provider-native response IDs and timestamps intentionally remain in payloads.
HTTP headers and capture timestamps are intentionally absent.

**Cost warning:** the live command makes billable API requests. Prompts and token
limits are small, but inspect the dry run and your selected models before capture.

**Inspect before commit:** provider payloads can contain generated content,
request/response IDs, timestamps, model revisions, moderation metadata, or proxy
extensions. Review every generated file before adding it to version control.
These are volatile, raw-provider evidence fixtures. They differ from normalized,
deterministic unit snapshots, which remove or canonicalize provider-specific
fields and should remain the primary stable regression tests.
