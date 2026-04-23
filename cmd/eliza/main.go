// Command eliza starts an HTTP server that serves both the Anthropic and
// OpenAI chat APIs backed by a classic ELIZA psychotherapist.
//
// Usage:
//
//	go run ./cmd/eliza [--addr :8080]
//
// Then point any Anthropic or OpenAI SDK at the server:
//
//	# Anthropic SDK
//	export ANTHROPIC_BASE_URL=http://localhost:8080
//	# OpenAI SDK
//	export OPENAI_BASE_URL=http://localhost:8080/v1
//	# Codex
//	codex -c 'openai_base_url="http://localhost:8080/v1"' -m eliza
//
// All models are accepted. Streaming works. Token counts are approximated.
package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/langchain-ai/langsmith-go/internal/mockllm"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	srv := mockllm.NewCombinedServer(mockllm.WithHandler(mockllm.ElizaHandler()))

	// Copy the handler from the httptest server — we want to listen on a
	// real address, not httptest's random port.
	handler := srv.Server.Config.Handler
	srv.Close() // close the httptest server; we'll use our own

	log.Printf("Eliza is listening on %s", *addr)
	log.Printf("  Anthropic API: POST http://localhost%s/v1/messages", *addr)
	log.Printf("  OpenAI API:    POST http://localhost%s/v1/chat/completions", *addr)
	log.Printf("  Responses API: POST http://localhost%s/v1/responses", *addr)
	log.Printf("  Models:        GET  http://localhost%s/v1/models", *addr)

	wrapped := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		handler.ServeHTTP(w, r)
	})

	// Serve root with info
	mux := http.NewServeMux()
	mux.Handle("/v1/", wrapped)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			wrapped.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"name":    "eliza",
			"version": "1.0.0",
			"status":  "Tell me about your problems.",
		})
	})

	log.Fatal(http.ListenAndServe(*addr, mux))
}
