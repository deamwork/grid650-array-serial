package httpserver

import (
	"net/http"
)

type Router interface {
	RegisterRouteHandler(method string, pattern string, handler http.Handler) *HTTPServer
	RegisterRouteHandlerFunc(method string, pattern string, handlerFunc http.HandlerFunc) *HTTPServer
}
