package middleware

import (
	"net/http"
)

type middleware interface {
	Attach(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request)
}

func Attach(handler func(w http.ResponseWriter, r *http.Request), middleware middleware) func(w http.ResponseWriter, r *http.Request) {
	return AttachMultiple(handler, Middlewares{
		middleware,
	})
}

type Middlewares []middleware

func AttachMultiple(handler func(w http.ResponseWriter, r *http.Request), middlewares Middlewares) func(w http.ResponseWriter, r *http.Request) {
	result := handler
	len := len(middlewares)

	for i := len - 1; i >= 0; i-- {
		middleware := middlewares[i]
		result = middleware.Attach(result)
	}

	return result
}

type HttpHandler struct {
	handleFunc func(w http.ResponseWriter, r *http.Request)
}

func (h HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handleFunc(w, r)
}

func AttachToHandler(handler http.Handler, middleware middleware) http.Handler {
	return AttachMultipleToHandler(handler, Middlewares{
		middleware,
	})
}

func AttachMultipleToHandler(handler http.Handler, middlewares Middlewares) http.Handler {
	result := handler
	len := len(middlewares)

	for i := len - 1; i >= 0; i-- {
		middleware := middlewares[i]

		innerResult := result

		handleFunc := func(w http.ResponseWriter, r *http.Request) {
			innerResult.ServeHTTP(w, r)
		}

		httpHandler := middleware.Attach(handleFunc)

		result = HttpHandler{
			handleFunc: httpHandler,
		}
	}

	return result
}
