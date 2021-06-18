package patriot_router

import (
	"fmt"
	"net/http"
)

type Router struct {
	routes map[*Route]func(http.ResponseWriter, *http.Request)
}

func (self *Router) RegisterRoute(route *Route, handler func(http.ResponseWriter, *http.Request)) {
	self.routes[route] = handler
}

func (self *Router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	request_path := request.URL.Path
	fmt.Printf("Serving '%s' for %s by %s\n", request.Method, request_path, request.RemoteAddr)
	for route, handler := range self.routes {
		if route.Match(request_path) {
			handler(response, request)
			return
		}
	}

	response.WriteHeader(404)
	response.Write([]byte("not found"))
}

func CreateRouter() *Router {
	var new_router *Router = new(Router)
	new_router.routes = make(map[*Route]func(http.ResponseWriter, *http.Request))

	return new_router
}
