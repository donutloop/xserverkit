package middleware

import "net/http"

type contextKey string

type Middleware func(http.Handler) http.Handler

// Chain acts as a list of middleware.
// Chain is effectively immutable:
// once created, it will always hold
// the same set of constructors in the same order.
type Chain struct {
	middleware []Middleware
}

// New creates a new chain,
// memorizing the given list of middleware.
// New serves no other function,
// constructors are only called upon a call to Then().
func New(middleware ...Middleware) Chain {
	return Chain{middleware: append(make([]Middleware, 0, len(middleware)), middleware...)}
}

// Copy middleware of existing chain and takes that as base for a new chain
// stdChain := easy_middleware.New(m1, m2)
// stdChainCopy := stdChain.Copy()
func (c Chain) Copy() Chain {
	return Chain{middleware: c.middleware}
}

// Add a middleware to a existing chain
// stdChain := easy_middleware.New(m1, m2)
// stdChainCopy := stdChain.Copy()
// stdChainCopy := stdChain.Add(m1)
// or
// stdChainCopy := stdChain.Add(m1, m2)
func (c *Chain) Add(middleware ...Middleware) {
	c.middleware = append(c.middleware, middleware...)
}

// Then chains the middleware and returns the final http.Handler.
//     New(m1, m2, m3).Then(h)
// is equivalent to:
//     m1(m2(m3(h)))
// When the request comes in, it will be passed to m1, then m2, then m3
// and finally, the given handler
// (assuming every middleware calls the following one).
//
// A chain can be safely reused by calling Then() several times.
//     stdStack := easy_middleware.New(ratelimitHandler, csrfHandler)
//     indexPipe = stdStack.Then(indexHandler)
//     authPipe = stdStack.Then(authHandler)
// Note that constructors are called on every call to Then()
// and thus several instances of the same middleware will be created
// when a chain is reused in this way.
// For proper middleware, this should cause no problems.
//
// Then() treats nil as http.DefaultServeMux.
func (c Chain) Then(endpoint http.Handler) http.Handler {

	if endpoint == nil {
		endpoint = http.DefaultServeMux
	}

	for i := range c.middleware {
		endpoint = c.middleware[len(c.middleware)-1-i](endpoint)
	}

	return endpoint
}

// ThenFunc works identically to Then, but takes
// a HandlerFunc instead of a Handler.
//
// The following two statements are equivalent:
//     c.Then(http.HandlerFunc(fn))
//     c.ThenFunc(fn)
//
// ThenFunc provides all the guarantees of Then.
func (c Chain) ThenFunc(endpointFunc http.HandlerFunc) http.Handler {

	if endpointFunc == nil {
		return c.Then(nil)
	}

	return c.Then(endpointFunc)
}
