package middleware

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/jamal23041989/go-marketplace-inventory-service/internal/core/logger"
)

type Middleware struct {
	logger logger.Logger
}

func New(logger logger.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		m.logger.Info("Started %s %s", r.Method, r.RequestURI)

		rw := NewStatusRecorder(w)
		next.ServeHTTP(rw, r)

		m.logger.Info(
			"Method: %s, URI: %s, StatusCode: %d, Time: %s",
			r.Method,
			r.RequestURI,
			rw.statusCode,
			time.Since(start),
		)
	})
}

func (m *Middleware) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.logger.Error("panic recovered: %v\n%s", err, debug.Stack())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
