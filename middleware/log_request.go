package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			status:          http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		// log base (limpo)
		logFn := slog.Info
		msg := "request"

		// destaque para erros
		if rw.status >= 400 {
			logFn = slog.Error
			msg = "request_error"
		}

		logFn(msg,
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"duration_ms", duration.Milliseconds(),
		)
	})
}