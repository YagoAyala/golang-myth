package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic recovered",
						slog.String("error", err.(error).Error()),
						slog.String("stack", string(debug.Stack())),
						slog.String("method", r.Method),
						slog.String("path", r.URL.Path),
					)

					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"error":"Internal server error"}`))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
