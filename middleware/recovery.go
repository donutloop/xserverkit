package middleware

import (
	"net/http"
	"net/http/httputil"
	"runtime/debug"

)

// Recovery middleware for panic
func Recovery(loggerFunc  func(requestDump []byte, stackDump []byte)) Middleware {
	return func(h http.Handler) http.Handler {
		funcWrapper := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {

					requestDump, err := httputil.DumpRequest(r, true)
					if err != nil {
						requestDump = make([]byte, 0)
					}

					loggerFunc(requestDump, debug.Stack())

					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(funcWrapper)
	}
}
