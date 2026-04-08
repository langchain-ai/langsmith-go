# /// script
# requires-python = ">=3.11"
# dependencies = ["langchain-openai>=0.3", "langchain-anthropic>=0.3"]
# ///
"""Test Eliza mock server with LangChain ChatModels."""

import json
import os
import sys

base_url = os.environ["ELIZA_BASE_URL"]

results = {}

# --- ChatOpenAI invoke ---
from langchain_openai import ChatOpenAI

llm = ChatOpenAI(model="gpt-4o", base_url=base_url + "/v1", api_key="fake")
resp = llm.invoke("hello")
results["openai_invoke"] = resp.content

# --- ChatOpenAI stream ---
chunks = list(llm.stream("I am feeling anxious"))
results["openai_stream_chunks"] = len(chunks)
full_text = "".join(c.content for c in chunks)
results["openai_stream_text"] = full_text

# --- ChatAnthropic invoke ---
from langchain_anthropic import ChatAnthropic

llm2 = ChatAnthropic(
    model="claude-sonnet-4-20250514",
    base_url=base_url,
    api_key="fake",
)
resp2 = llm2.invoke("hello")
results["anthropic_invoke"] = resp2.content

# --- ChatAnthropic stream ---
chunks2 = list(llm2.stream("I need help"))
results["anthropic_stream_chunks"] = len(chunks2)
full_text2 = "".join(c.content for c in chunks2)
results["anthropic_stream_text"] = full_text2

json.dump(results, sys.stdout)
