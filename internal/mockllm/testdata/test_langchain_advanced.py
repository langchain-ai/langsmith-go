# /// script
# requires-python = ">=3.11"
# dependencies = ["langchain-openai>=0.3", "langchain-anthropic>=0.3", "langchain-core>=0.3"]
# ///
"""Advanced LangChain tests: multi-turn, tool calling, system messages."""

import json
import os
import sys

from langchain_core.messages import HumanMessage, SystemMessage, AIMessage
from langchain_core.tools import tool

base_url = os.environ["ELIZA_BASE_URL"]
results = {}

# --- OpenAI multi-turn ---
from langchain_openai import ChatOpenAI

llm = ChatOpenAI(model="gpt-4o", base_url=base_url + "/v1", api_key="fake")

resp = llm.invoke([
    SystemMessage(content="You are helpful."),
    HumanMessage(content="What is 2+2?"),
    AIMessage(content="4"),
    HumanMessage(content="And times 3?"),
])
results["openai_multiturn"] = resp.content
results["openai_multiturn_nonempty"] = len(resp.content) > 0

# --- Anthropic multi-turn ---
from langchain_anthropic import ChatAnthropic

llm2 = ChatAnthropic(
    model="claude-sonnet-4-20250514",
    base_url=base_url,
    api_key="fake",
)

resp2 = llm2.invoke([
    SystemMessage(content="You are helpful."),
    HumanMessage(content="hello"),
])
results["anthropic_system"] = resp2.content
results["anthropic_system_nonempty"] = len(resp2.content) > 0

# --- OpenAI with tools ---
@tool
def get_weather(location: str) -> str:
    """Get the weather for a location."""
    return f"Sunny in {location}"

llm_tools = ChatOpenAI(model="gpt-4o", base_url=base_url + "/v1", api_key="fake")
llm_with_tools = llm_tools.bind_tools([get_weather])
resp3 = llm_with_tools.invoke("What is the weather in Paris?")
results["openai_tool_calls"] = len(resp3.tool_calls) if hasattr(resp3, "tool_calls") else 0
if resp3.tool_calls:
    results["openai_tool_name"] = resp3.tool_calls[0]["name"]
else:
    results["openai_tool_name"] = ""

# --- Anthropic with tools ---
llm2_tools = ChatAnthropic(
    model="claude-sonnet-4-20250514",
    base_url=base_url,
    api_key="fake",
)
llm2_with_tools = llm2_tools.bind_tools([get_weather])
resp4 = llm2_with_tools.invoke("What is the weather in Paris?")
results["anthropic_tool_calls"] = len(resp4.tool_calls) if hasattr(resp4, "tool_calls") else 0
if resp4.tool_calls:
    results["anthropic_tool_name"] = resp4.tool_calls[0]["name"]
else:
    results["anthropic_tool_name"] = ""

# --- OpenAI streaming with system message ---
chunks = list(llm.stream([
    SystemMessage(content="Be concise."),
    HumanMessage(content="hello"),
]))
results["openai_system_stream_chunks"] = len(chunks)
results["openai_system_stream_text"] = "".join(c.content for c in chunks)

# --- Anthropic streaming multi-turn ---
chunks2 = list(llm2.stream([
    HumanMessage(content="My name is Alice"),
    AIMessage(content="Hello Alice"),
    HumanMessage(content="What is my name?"),
]))
results["anthropic_multiturn_stream_chunks"] = len(chunks2)
results["anthropic_multiturn_stream_text"] = "".join(c.content for c in chunks2)

json.dump(results, sys.stdout)
