package middleware

import (
	"net/url"
	"strings"
	"net/http"
	"context"
	"fmt"
)

// URLQueryKey is the context key for the URL Query
const URLQueryKey contextKey = "urlquery"

// URLQuery is a middleware to parse the URL Query parameters just once,
// and store the resulting url.Values in the context.
func URLQuery() Middleware {
	return func(h http.Handler) http.Handler {
		funcWrapper := func(w http.ResponseWriter, r *http.Request) {
			q, err := extractURLQueries(r)
			if err != nil {
				r = r.WithContext(context.WithValue(r.Context(), URLQueryKey, err))
			} else {
				r = r.WithContext(context.WithValue(r.Context(), URLQueryKey, q))
			}

			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(funcWrapper)
	}
}

func extractURLQueries(req *http.Request) (*Queries, error) {
	queriesRaw, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		return nil, err
	}

	queries := &Queries{
		C: make(map[string][]string),
	}
	if 0 == len(queriesRaw) {
		return queries, nil
	}

	for k, v := range queriesRaw {
		for _, item := range v {
			values := strings.Split(item, ",")
			queries.C[k] = append(queries.C[k], values...)
		}
	}

	return queries, nil
}

type Queries struct {
	C map[string][]string
}

// Get return the key value, of the current *http.Request queries
func (q Queries) Get(key string) []string {
	if value, found := q.C[key]; found {
		return value
	}
	return make([]string, 0)
}

// Get returns all queries of the current *http.Request queries
func (q Queries) GetAll() map[string][]string {
	return q.C
}

// Count returns count of the current *http.Request queries
func (q Queries) Count() int {
	return len(q.C)
}

func GetQueries(ctx context.Context) (*Queries, error) {
	v := ctx.Value(URLQueryKey)
	queries, ok := v.(*Queries)
	if ok {
		return queries, nil
	}
	err, ok := v.(error)
	if ok {
		return nil, err
	}
	return nil, fmt.Errorf("type not supported (%v)", v)
}