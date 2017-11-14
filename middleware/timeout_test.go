package middleware_test

import (
	"testing"
	"net/http/httptest"
	"time"
	"net/http"
	"github.com/donutloop/xserverkit/middleware"
)

func TestTimeout(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		<- time.After(2 * time.Second)

		select {
		case <- r.Context().Done():
			return
		}

		w.WriteHeader(http.StatusOK)
	}

	testHandler := http.HandlerFunc(handler)
	test := httptest.NewServer(middleware.Timeout(1*time.Second)(testHandler))
	defer test.Close()

	response, err := http.Get(test.URL)
	if err != nil {
		t.Errorf("timeout middleware request: %v", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusGatewayTimeout {
		t.Error("timeout middleware request: Unexpected good request")
	}
}
