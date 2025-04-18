package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(ms ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(ms) - 1; i >= 0; i-- {
			next = ms[i](next)
		}
		return next
	}
}
