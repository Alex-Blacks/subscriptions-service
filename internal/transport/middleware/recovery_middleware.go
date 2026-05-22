package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/Alex-Blacks/subscriptions/internal/logging"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}

			logger := logging.LoggerFromContext(r.Context())
			if logger == nil {
				logger = slog.Default()
			}

			logger.Error("recovery panic",
				"error", rec,
				"stack", string(debug.Stack()),
			)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)

			_, _ = w.Write([]byte(`{"error":"internal server error"}`))

		}()
		next.ServeHTTP(w, r)
	})
}
