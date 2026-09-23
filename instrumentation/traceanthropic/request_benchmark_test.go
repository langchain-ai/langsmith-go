package traceanthropic

import (
	"context"
	"fmt"
	"strings"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Compare prompt capture with the non-recording path on synthetic request bodies.
func BenchmarkRequestAttributes(b *testing.B) {
	for _, size := range []int{1024, 64 * 1024} {
		body := []byte(`{"model":"test-model","stream":true,"system":"Be helpful","messages":[{"role":"user","content":"` + strings.Repeat("x", size) + `"}]}`)
		for _, recording := range []bool{true, false} {
			b.Run(fmt.Sprintf("bytes=%d/recording=%t", size, recording), func(b *testing.B) {
				sampler := sdktrace.NeverSample()
				if recording {
					sampler = sdktrace.AlwaysSample()
				}
				tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sampler))
				defer tp.Shutdown(context.Background())
				_, span := tp.Tracer("benchmark").Start(context.Background(), "request")
				defer span.End()
				b.ReportAllocs()
				b.SetBytes(int64(len(body)))
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					extractRequestAttributes(span, body)
				}
			})
		}
	}
}
