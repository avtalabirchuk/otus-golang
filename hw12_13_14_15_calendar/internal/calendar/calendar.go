package calendar

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/repository"
)

type Calendar struct {
	r repository.Base
}

func New(r repository.Base) (*Calendar, error) {
	return &Calendar{r: r}, nil
}

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)
		latency := time.Since(start)
		// 66.249.65.3 GET /hello?q=1 HTTP/1.1 200 30 "Mozilla/5.0"
		log.Info().Msgf("%s %s %s %s %d %d %s", r.RemoteAddr, r.Method, r.RequestURI, r.Proto, rw.statusCode, latency, r.UserAgent())
	})
}

func (c *Calendar) dummyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func (c *Calendar) Run(addr string) error {
	mux := http.NewServeMux()

	mux.Handle("/hello", logMiddleware(http.HandlerFunc(c.dummyHandler)))

	log.Info().Msgf("Listening at %s", addr)

	return http.ListenAndServe(addr, mux)
}
