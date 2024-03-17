package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func LoggingMiddleware(next http.Handler, log *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
