package mockllm

import "net/http"

// hijackAndClose abruptly closes the underlying TCP connection, simulating
// a network-level failure. If the ResponseWriter doesn't support hijacking
// (shouldn't happen with httptest.Server), it falls back to closing the
// response with no body.
func hijackAndClose(w http.ResponseWriter) {
	if hj, ok := w.(http.Hijacker); ok {
		conn, _, err := hj.Hijack()
		if err == nil {
			conn.Close()
			return
		}
	}
	// Fallback: write nothing and let the server close naturally
	w.WriteHeader(http.StatusBadGateway)
}
