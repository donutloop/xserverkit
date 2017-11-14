package middleware_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/donutloop/xserverkit/middleware"
)

func TestLogging(t *testing.T) {

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(map[string]string{"echo": "test"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	b := new(bytes.Buffer)
	logger := log.New(b, "", 0)

	loggingCallback := func(s string) {
		logger.Println(s)
	}

	testHandler := http.HandlerFunc(handler)
	test := httptest.NewServer(middleware.Logging(loggingCallback)(testHandler))
	defer test.Close()

	response, err := http.Get(test.URL)
	if err != nil {
		t.Errorf("logging middleware request: %v", err)
		return
	}
	defer response.Body.Close()

	if !strings.Contains(b.String(), "completed in") || !strings.Contains(b.String(), "started") {
		t.Errorf("logging middleware request: log output should match %q is a string", b.String())
	}
}
