package middleware_test

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"github.com/donutloop/xserverkit/middleware"
)

func TestURLQuery(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		q, err :=  middleware.GetQueries(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		limit := q.Get("limit")[0]
		if limit != "10" {
			w.WriteHeader(http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusOK)
	}

	testHandler := http.HandlerFunc(handler)

	test := httptest.NewServer(middleware.URLQuery()(testHandler))
	defer test.Close()

	response, err := http.Get(test.URL + "?limit=10")
	if err != nil {
		t.Errorf("url query middleware request: %v", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("url query middleware request: unexpected bad request (%d)", response.StatusCode)
	}
}
