package main

import (
	"fmt"
	"github.com/andviro/noodle"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
)

// HR adapts noodle middleware chain to httprouter handler signature.
// httprouter.Params are injected into context under key "Params"
func HR(h noodle.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(noodle.Factory(), "Params", p)
		h(ctx, w, r)
	}
}

func HTTPAuth(authFunc func(username, password string) bool) noodle.Middleware {
	return func(next noodle.Handler) noodle.Handler {
		return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
			username, password, ok := r.BasicAuth()
			if !ok || !authFunc(username, password) {
				w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
				w.WriteHeader(http.StatusUnauthorized)
				return fmt.Errorf("Unauthorized request")
			}
			userContext := context.WithValue(c, "User", username)
			return next(userContext, w, r)
		}
	}
}

func index(c context.Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintln(w, "index")
	return nil
}

func products(c context.Context, w http.ResponseWriter, r *http.Request) error {
	params := c.Value("Params").(httprouter.Params)
	user := c.Value("User").(string)
	fmt.Fprintf(w, "user: %s, products: %s", user, params.ByName("id"))
	return nil
}

func main() {
	r := httprouter.New()
	basicAuth := HTTPAuth(func(username, password string) bool {
		return password == "secret"
	})
	n := noodle.Default()
	r.GET("/", HR(n.Then(index)))
	r.GET("/products/:id", HR(n.Use(basicAuth).Then(products)))
	http.ListenAndServe(":8080", r)
}
