package logmiddleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func ApplyHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)
		latency := time.Since(start)
		// 66.249.65.3 GET /hello?q=1 HTTP/1.1 200 30 "Mozilla/5.0"
		log.Info().Msgf("[HTTP] %s %s %s %s %d %s %s", r.RemoteAddr, r.Method, r.RequestURI, r.Proto, rw.statusCode, latency, r.UserAgent())
	})
}
