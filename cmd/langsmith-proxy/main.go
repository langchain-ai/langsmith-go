// Command langsmith-proxy starts a reverse proxy that adds LangSmith tracing
// to Anthropic and OpenAI API calls.
//
// Usage:
//
//	go run ./cmd/langsmith-proxy [--addr :8090] [--anthropic-upstream URL] [--openai-upstream URL]
//
// Environment:
//
//	LANGSMITH_API_KEY   — LangSmith API key (required for real tracing)
//	LANGSMITH_PROJECT   — LangSmith project name (default: "default")
//	LANGSMITH_ENDPOINT  — LangSmith endpoint (default: api.smith.langchain.com)
//
// Then point any SDK at the proxy:
//
//	export ANTHROPIC_BASE_URL=http://localhost:8090
//	export OPENAI_BASE_URL=http://localhost:8090/v1
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/langchain-ai/langsmith-go/instrumentation/proxy"
)

func main() {
	addr := flag.String("addr", ":8090", "listen address")
	anthropicUpstream := flag.String("anthropic-upstream", "", "Anthropic upstream URL (default: https://api.anthropic.com)")
	openaiUpstream := flag.String("openai-upstream", "", "OpenAI upstream URL (default: https://api.openai.com)")
	flag.Parse()

	p, err := proxy.New(proxy.Config{
		AnthropicUpstream: *anthropicUpstream,
		OpenAIUpstream:    *openaiUpstream,
	})
	if err != nil {
		log.Fatalf("Failed to create proxy: %v", err)
	}

	aUp := *anthropicUpstream
	if aUp == "" {
		aUp = "https://api.anthropic.com"
	}
	oUp := *openaiUpstream
	if oUp == "" {
		oUp = "https://api.openai.com"
	}

	log.Printf("LangSmith tracing proxy listening on %s", *addr)
	log.Printf("  Anthropic: /v1/messages          -> %s", aUp)
	log.Printf("  OpenAI:    /v1/chat/completions  -> %s", oUp)
	log.Printf("  OpenAI:    /v1/responses         -> %s", oUp)
	log.Printf("  Passthru:  /v1/models            -> %s", oUp)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	server := &http.Server{Addr: *addr, Handler: p}
	go func() {
		<-ctx.Done()
		log.Println("Shutting down...")
		p.Shutdown(context.Background())
		server.Shutdown(context.Background())
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
