package router

import (
	"net/http"
	"backend/internal/middleware"
)

type Router struct {
	mux         *http.ServeMux
	middlewares []middleware.Middleware
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

func (r *Router) Use(mw ...middleware.Middleware) {
	r.middlewares = append(r.middlewares, mw...)
}

func (r *Router) HandleFunc(pattern string, handler http.HandlerFunc, mw ...middleware.Middleware) {
	finalHandler := http.Handler(handler)
	
	finalHandler = middleware.Chain(finalHandler, mw...)
	finalHandler = middleware.Chain(finalHandler, r.middlewares...)
	
	r.mux.Handle(pattern, finalHandler)
}

func (r *Router) Handle(pattern string, handler http.Handler, mw ...middleware.Middleware) {
	finalHandler := handler
	finalHandler = middleware.Chain(finalHandler, mw...)
	finalHandler = middleware.Chain(finalHandler, r.middlewares...)
	r.mux.Handle(pattern, finalHandler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
