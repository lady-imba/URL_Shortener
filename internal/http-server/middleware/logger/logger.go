package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	log = log.With(
		slog.String("component", "middleware/logger"),
	)

	log.Info("logger middleware enabled")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()
			
			entry.Debug("request started")

			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.String("duration", time.Since(start).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}