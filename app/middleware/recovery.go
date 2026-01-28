package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/mytheresa/go-hiring-challenge/app/api"
)

func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					var errMsg string
					if e, ok := err.(error); ok {
						errMsg = e.Error()
					} else {
						errMsg = fmt.Sprint(err)
					}

					logger.Error("panic recovered",
						slog.String("error", errMsg),
						slog.String("stack", string(debug.Stack())),
						slog.String("method", r.Method),
						slog.String("path", r.URL.Path),
					)

					api.ErrorResponse(w, http.StatusInternalServerError, "Internal server error")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
