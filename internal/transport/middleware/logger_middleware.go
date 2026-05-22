package middleware

import (
	"bufio"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/Alex-Blacks/subscriptions/internal/logging"
)

type wrapped struct {
	w       http.ResponseWriter
	status  int
	written int
}

func (wr *wrapped) Header() http.Header {
	return wr.w.Header()
}

func (wr *wrapped) Write(b []byte) (int, error) {
	if wr.status == 0 {
		wr.status = http.StatusOK
	}
	n, err := wr.w.Write(b)
	wr.written += n
	return n, err
}

func (wr *wrapped) WriteHeader(code int) {
	wr.status = code
	wr.w.WriteHeader(code)
}

func (wr *wrapped) Flush() {
	if f, ok := wr.w.(http.Flusher); ok {
		f.Flush()
	}
}

func (wr *wrapped) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := wr.w.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return h.Hijack()
}

var (
	_ http.Flusher  = (*wrapped)(nil)
	_ http.Hijacker = (*wrapped)(nil)
)

func LoggerMiddleware(baseLogger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			log := baseLogger.With(
				"method", r.Method,
				"path", r.URL.Path,
			)

			ctx := logging.LoggerWithContext(r.Context(), log)
			r = r.WithContext(ctx)

			wr := &wrapped{w: w}

			defer func() {
				log.Info("request finished",
					"status", wr.status,
					"bytes", wr.written,
					"duration", time.Since(start),
				)
			}()

			log.Info("request started")
			next.ServeHTTP(wr, r)
		})

	}
}
