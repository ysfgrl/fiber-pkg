package routes

import (
	"github.com/gofiber/fiber/v2"
)

var MethodPost = "POST"
var MethodGet = "GET"
var MethodPut = "PUT"
var MethodDelete = "DELETE"
var MethodPatch = "PATCH"
var MethodOptions = "OPTIONS"

type Route struct {
	Path    string
	Method  string
	Handler fiber.Handler
	Role    []string
	Auth    bool
	Static  bool
	Log     bool
}

func Static(path string, handler fiber.Handler) Route {
	return Route{
		Method:  MethodGet,
		Path:    path,
		Handler: handler,
		Auth:    false,
		Role:    nil,
		Static:  true,
	}
}

func Post(path string, handler fiber.Handler, auth bool, roles ...string) *Route {
	return &Route{
		Method:  MethodPost,
		Path:    path,
		Handler: handler,
		Auth:    auth,
		Role:    roles,
		Log:     true,
	}
}

func Get(path string, handler fiber.Handler, auth bool, roles ...string) *Route {
	return &Route{
		Method:  MethodGet,
		Path:    path,
		Handler: handler,
		Auth:    auth,
		Role:    roles,
		Log:     true,
	}
}

func Put(path string, handler fiber.Handler, auth bool, roles ...string) *Route {
	return &Route{
		Method:  MethodPut,
		Path:    path,
		Handler: handler,
		Auth:    auth,
		Role:    roles,
		Log:     true,
	}
}
func Delete(path string, handler fiber.Handler, auth bool, roles ...string) *Route {
	return &Route{
		Method:  MethodDelete,
		Path:    path,
		Handler: handler,
		Auth:    auth,
		Role:    roles,
		Log:     true,
	}
}

type Interface interface {
	Routes(app *fiber.App)
	AddRoute(route *Route)
}
