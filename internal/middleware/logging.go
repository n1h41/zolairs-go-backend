package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter is a wrapper for http.ResponseWriter that captures status code
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

// WrapResponseWriter wraps an http.ResponseWriter to capture status code
func WrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

// Status returns the status code
func (rw *responseWriter) Status() int {
	return rw.status
}

// LoggingMiddleware logs requests and their responses
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s - Request started", r.Method, r.URL.Path)

		// Capture response status
		ww := WrapResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(ww, r)

		// Log request details after completion
		log.Printf("%s %s - Status: %d - Duration: %v",
			r.Method, r.URL.Path, ww.Status(), time.Since(start))
	})
}

