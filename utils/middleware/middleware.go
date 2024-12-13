package middleware

import "net/http"

type Middleware func(http.RoundTripper) http.RoundTripper

func ChainMiddleware(base http.RoundTripper, middlewares ...Middleware) http.RoundTripper {
	for i := len(middlewares) - 1; i >= 0; i-- {
		base = middlewares[i](base)
	}
	return base
}
