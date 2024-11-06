package server

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// LoggingMiddleware logs incoming HTTP requests
func LoggingMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			logrus.WithFields(logrus.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
			}).Info("Request started")

			next.ServeHTTP(w, r)

			logrus.WithFields(logrus.Fields{
				"duration": time.Since(start),
			}).Info("Request completed")
		})
	}
}
