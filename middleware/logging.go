package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// Logging of device request time
func Logging(loggerFunc func(s string)) Middleware {
	return func(h http.Handler) http.Handler {
		funcWrapper := func(w http.ResponseWriter, r *http.Request) {
			defer logWrapper(loggerFunc, r)()
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(funcWrapper)
	}
}

func logWrapper(loggerFunc func(s string), r *http.Request) func() {
	start := time.Now()
	loggerFunc(fmt.Sprintf("Method: %s, url: %s, agent: %s started", r.Method, r.URL.Path, r.UserAgent()))
	return func() {
		loggerFunc(fmt.Sprintf("Method: %s, url: %s, agent: %s completed in %v", r.Method, r.URL.Path, r.UserAgent(), time.Since(start)))
	}
}
