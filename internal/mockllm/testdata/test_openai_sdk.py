# /// script
# requires-python = ">=3.11"
# dependencies = ["openai>=1.0"]
# ///
"""Test Eliza mock server with the OpenAI Python SDK."""

import json
import os
import sys

from openai import OpenAI

base_url = os.environ["ELIZA_BASE_URL"] + "/v1"
client = OpenAI(base_url=base_url, api_key="fake")

results = {}

# --- Non-streaming ---
resp = client.chat.completions.create(
    model="gpt-4o",
    messages=[{"role": "user", "content": "hello"}],
)
results["nonstreaming_content"] = resp.choices[0].message.content
results["nonstreaming_model"] = resp.model
results["nonstreaming_has_usage"] = resp.usage is not None

# --- Streaming ---
stream = client.chat.completions.create(
    model="gpt-4o",
    messages=[{"role": "user", "content": "I am sad"}],
    stream=True,
)
chunks = []
for chunk in stream:
    if chunk.choices and chunk.choices[0].delta.content:
        chunks.append(chunk.choices[0].delta.content)
results["streaming_text"] = "".join(chunks)
results["streaming_chunk_count"] = len(chunks)

# --- Tool calls ---
resp = client.chat.completions.create(
    model="gpt-4o",
    messages=[{"role": "user", "content": "weather?"}],
    tools=[{
        "type": "function",
        "function": {"name": "get_weather", "parameters": {}},
    }],
)
choice = resp.choices[0]
results["tool_finish_reason"] = choice.finish_reason
results["tool_call_name"] = choice.message.tool_calls[0].function.name if choice.message.tool_calls else ""

# --- Error ---
try:
    bad = OpenAI(base_url=base_url, api_key="sk-invalid-test")
    bad.chat.completions.create(
        model="gpt-4o",
        messages=[{"role": "user", "content": "hi"}],
    )
    results["error_raised"] = False
except Exception:
    results["error_raised"] = True

json.dump(results, sys.stdout)
