package main

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

type paramsKey int

const myParamsKey paramsKey = 0

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		context.Set(r, myParamsKey, ps)
		h.ServeHTTP(w, r)
	}
}

type Router struct {
	router *httprouter.Router
}

func NewRouter() *Router {
	return &Router{httprouter.New()}
}

func (r *Router) Get(route string, h http.Handler) {
	r.router.GET(route, wrapHandler(h))
}

func (r *Router) ListenAndServe(addr string) {
	http.ListenAndServe(addr, r.router)
}
