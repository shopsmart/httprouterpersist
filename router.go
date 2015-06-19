/*
The httprouterpersist package is a wrapper for the httprouter package that
allows customized options of params persistence without using the non-standard
httprouter.Handle interface. This allows httprouter to be used with standard
http.HandlerFunc implementations.

By default, url params with disappear into a blackhole. However, that can be
persisted by using a context package, by attaching to form values, or other
ways.

A trivial example is:

	package main

	import (
		"fmt"
		"github.com/gorilla/context"
		"github.com/shopsmart/cms/router"
		"net/http"
		"log"
	)

	func Index(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome!\n")
	}

	func Hello(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, %s!\n", context.Get(r, "name"))
	}

	func main() {
		r := router.New()
		r.Persist = router.ContextPersistParamsFunc

		r.GET("/", Index)
		r.GET("/hello/:name", Hello)

		log.Fatal(http.ListenAndServe(":8080", r))
	}
*/
package httprouterpersist

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

/*
The Router type wraps an httprouter.Router type and provides new GET, POST, etc.
methods that wrap the httprouter calls with a persist params function.

The Persist attribute should be set to a function that can persist or discard
the httprouter params.
*/
type Router struct {
	*httprouter.Router
	Persist PersistParamsFunc
}

/*
Returns a new, intialized router that will discard httprouter params.
*/
func New() *Router {
	return &Router{httprouter.New(), BlackholePersist}
}

func (r *Router) Handle(method, path string, fn http.HandlerFunc) {
	r.Router.Handle(method, path, r.wrapHandler(fn))
}

func (r *Router) DELETE(path string, fn http.HandlerFunc) {
	r.Router.DELETE(path, r.wrapHandler(fn))
}

func (r *Router) GET(path string, fn http.HandlerFunc) {
	r.Router.GET(path, r.wrapHandler(fn))
}

func (r *Router) HEAD(path string, fn http.HandlerFunc) {
	r.Router.HEAD(path, r.wrapHandler(fn))
}

func (r *Router) OPTIONS(path string, fn http.HandlerFunc) {
	r.Router.OPTIONS(path, r.wrapHandler(fn))
}

func (r *Router) PATCH(path string, fn http.HandlerFunc) {
	r.Router.PATCH(path, r.wrapHandler(fn))
}

func (r *Router) POST(path string, fn http.HandlerFunc) {
	r.Router.POST(path, r.wrapHandler(fn))
}

func (r *Router) PUT(path string, fn http.HandlerFunc) {
	r.Router.PUT(path, r.wrapHandler(fn))
}

/*
The PersistParamsFunc type is the signature for functions that can be used
to persist httprouter params.
*/
type PersistParamsFunc func(*http.Request, httprouter.Params)

/*
A PersistParamsFunc implementation that discards httprouter params.
*/
func BlackholePersist(r *http.Request, ps httprouter.Params) {
	return
}

/*
A PersistParamsFunc implementation that assigns httprouter params to
a request context using gorilla context. The params will be attaches as
key, value pairs on the context.

	r.GET("/users/:id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "User ID: %s", context.Get("id"))
	})
*/
func ContextPersist(r *http.Request, ps httprouter.Params) {
	if len(ps) > 0 {
		for _, param := range ps {
			context.Set(r, param.Key, param.Value)
		}
	}
	return
}

/*
A PersistParamsFunc implementation that sets httprouter params to the
request url query params. This value will replace any existing value on
the url query params.

	r.GET("/users/:id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "User ID: %s", r.FormValue("id"))
	})
*/
func RequestPersist(r *http.Request, ps httprouter.Params) {
	if len(ps) > 0 {
		values := r.URL.Query()
		for _, param := range ps {
			values.Set(param.Key, param.Value)
		}
		r.URL.RawQuery = values.Encode()
		r.Form = nil
	}
	return
}

func (r *Router) wrapHandler(handlerFunc http.HandlerFunc) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		r.Persist(req, ps)
		handlerFunc(res, req)
	}
}
