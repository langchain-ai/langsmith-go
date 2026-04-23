/**
 * Test Eliza mock server with the Vercel AI SDK (@ai-sdk/openai).
 *
 * Run via: npx --yes tsx test_ai_sdk.mjs
 * Requires: ELIZA_BASE_URL env var
 */

import { createOpenAI } from "@ai-sdk/openai";
import { generateText, streamText, tool } from "ai";
import { z } from "zod";

const baseURL = process.env.ELIZA_BASE_URL + "/v1";
const provider = createOpenAI({ baseURL, apiKey: "fake" });

const results = {};

// --- generateText (non-streaming) ---
const gen = await generateText({
  model: provider("gpt-4o"),
  prompt: "hello",
});
results.generate_text = gen.text;
results.generate_has_usage = gen.usage != null;

// --- streamText ---
const stream = streamText({
  model: provider("gpt-4o"),
  prompt: "I am sad",
});
let streamedText = "";
for await (const chunk of stream.textStream) {
  streamedText += chunk;
}
results.stream_text = streamedText;
results.stream_has_text = streamedText.length > 0;

// --- tool calls ---
const toolGen = await generateText({
  model: provider("gpt-4o"),
  prompt: "weather?",
  tools: {
    get_weather: tool({
      description: "Get weather",
      parameters: z.object({ location: z.string().optional() }),
      execute: async () => ({ temp: 20 }),
    }),
  },
});
results.tool_text = toolGen.text;
results.tool_calls_count = toolGen.toolCalls?.length ?? 0;
results.tool_call_name = toolGen.toolCalls?.[0]?.toolName ?? "";

console.log(JSON.stringify(results));
