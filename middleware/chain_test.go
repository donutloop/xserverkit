package middleware_test


import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/donutloop/xserverkit/middleware"
)

func testMiddleware(tag string) middleware.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(tag))
			h.ServeHTTP(w, r)
		})
	}
}

func TestThenOrdersHandlersCorrectly(t *testing.T) {

	t1 := testMiddleware("t1\n")
	t2 := testMiddleware("t2\n")
	t3 := testMiddleware("t3\n")

	testEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("endpoint\n"))
	})

	chained := middleware.New(t1, t2, t3).Then(testEndpoint)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	chained.ServeHTTP(w, r)
	if w.Body.String() != "t1\nt2\nt3\nendpoint\n" {
		t.Errorf("Then does not order handlers correctly (Order: %s)", w.Body.String())
		return
	}
}

func TestCopyOrdersHandlersCorrectly(t *testing.T) {

	t1 := testMiddleware("t1\n")
	t2 := testMiddleware("t2\n")
	t3 := testMiddleware("t3\n")

	testEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("endpoint\n"))
	})

	chained := middleware.New(t1, t2, t3)

	copyChained := chained.Copy().Then(testEndpoint)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	copyChained.ServeHTTP(w, r)
	if w.Body.String() != "t1\nt2\nt3\nendpoint\n" {
		t.Errorf("Then does not order handlers correctly (Order: %s)", w.Body.String())
		return
	}
}

func TestAddOrdersHandlersCorrectly(t *testing.T) {

	t1 := testMiddleware("t1\n")

	testEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("endpoint\n"))
	})

	chained := middleware.New(t1)

	copyChained := chained.Copy()

	t2 := testMiddleware("t2\n")
	t3 := testMiddleware("t3\n")

	copyChained.Add(t2, t3)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	copyChained.Then(testEndpoint).ServeHTTP(w, r)
	if w.Body.String() != "t1\nt2\nt3\nendpoint\n" {
		t.Errorf("Then does not order handlers correctly (Order: %s)", w.Body.String())
		return
	}
}
