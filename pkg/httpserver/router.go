package httpserver

import "github.com/gorilla/mux"

type Router interface {
	Register(router *mux.Router)
}