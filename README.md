# **DEPRECATED** httprouterpersist

httprouterpersist is a wrapper for httprouter that allows you to use httprouter with standard http.HandleFunc by providing an alternative method for persisting params.

## A trivial example is:

```
package main

import (
	"fmt"
	"github.com/gorilla/context"
	"net/http"
	"log"

	router "github.com/shopsmart/httprouterpersist"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, %s!\n", context.Get(r, "name"))
}

func main() {
	r := router.New()
	r.Persist = router.ContextPersist

	r.GET("/", Index)
	r.GET("/hello/:name", Hello)

	log.Fatal(http.ListenAndServe(":8080", r))
}
```