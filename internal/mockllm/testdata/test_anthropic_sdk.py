# /// script
# requires-python = ">=3.11"
# dependencies = ["anthropic>=0.30"]
# ///
"""Test Eliza mock server with the Anthropic Python SDK."""

import json
import os
import sys

import anthropic

base_url = os.environ["ELIZA_BASE_URL"]
client = anthropic.Anthropic(base_url=base_url, api_key="fake")

results = {}

# --- Non-streaming ---
msg = client.messages.create(
    model="claude-sonnet-4-20250514",
    max_tokens=256,
    messages=[{"role": "user", "content": "hello"}],
)
results["nonstreaming_content"] = msg.content[0].text
results["nonstreaming_model"] = msg.model
results["nonstreaming_has_usage"] = msg.usage is not None
results["nonstreaming_stop_reason"] = msg.stop_reason

# --- Streaming ---
chunks = []
with client.messages.stream(
    model="claude-sonnet-4-20250514",
    max_tokens=256,
    messages=[{"role": "user", "content": "I am sad"}],
) as stream:
    for text in stream.text_stream:
        chunks.append(text)
results["streaming_text"] = "".join(chunks)
results["streaming_chunk_count"] = len(chunks)

# --- Tool use ---
msg = client.messages.create(
    model="claude-sonnet-4-20250514",
    max_tokens=256,
    messages=[{"role": "user", "content": "weather?"}],
    tools=[{
        "name": "get_weather",
        "description": "Get weather",
        "input_schema": {"type": "object", "properties": {}},
    }],
)
results["tool_stop_reason"] = msg.stop_reason
tool_blocks = [b for b in msg.content if b.type == "tool_use"]
results["tool_call_name"] = tool_blocks[0].name if tool_blocks else ""

# --- Error ---
try:
    bad = anthropic.Anthropic(base_url=base_url, api_key="sk-ant-invalid-test")
    bad.messages.create(
        model="claude-sonnet-4-20250514",
        max_tokens=10,
        messages=[{"role": "user", "content": "hi"}],
    )
    results["error_raised"] = False
except anthropic.AuthenticationError:
    results["error_raised"] = True

json.dump(results, sys.stdout)
