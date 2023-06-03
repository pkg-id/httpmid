package httpmid

import "net/http"

// Middleware is a function that wraps http.Handler.
type Middleware func(next http.Handler) http.Handler

// Then applies the middleware to the handler and returns the final http.Handler.
func (f Middleware) Then(h http.Handler) http.Handler { return f(h) }

// Reduce reduces multiple middlewares into one middleware.
// Reduce apply middlewares in reverse order.
// This will make the first middleware in the slice be the outermost
// middleware (i.e. the one that gets called first on a request).
// The last middleware in the slice will be the innermost middleware.
// The innermost will be called before the actual handler.
//
// Example:
//
//	Reduce(m1, m2, m3).Then(h)
//	will be equivalent to:
//	m1(m2(m3(h)))
func Reduce(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := range middlewares {
			next = middlewares[len(middlewares)-1-i].Then(next)
		}
		return next
	}
}
