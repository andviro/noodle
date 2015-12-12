# Wok

[![Build
Status](https://travis-ci.org/andviro/wok.svg?branch=master)](https://travis-ci.org/andviro/wok)
[![Coverage](http://gocover.io/_badge/github.com/andviro/wok?0)](http://gocover.io/github.com/andviro/wok)
[![GoDoc](http://godoc.org/github.com/andviro/wok?status.png)](http://godoc.org/github.com/andviro/wok)

A proof of concept web application router based on
[httprouter](https://github.com/julienschmidt/httprouter). Supports route
groups and per-group [noodle](https://github.com/andviro/noodle) middleware.
Sample application:

```go
package main

import (
	"fmt"
	mw "github.com/andviro/noodle/middleware"
	"github.com/andviro/noodle/render"
	"github.com/andviro/wok"
	"golang.org/x/net/context"
	"net/http"
)

func main() {
	// apiAuth guards access to api group
	apiAuth := mw.HTTPAuth("API", func(user, pass string) bool {
		return pass == "Secret"
	})
	// dashboardAuth guards access to dashboard group
	dashboardAuth := mw.HTTPAuth("Dashboard", func(user, pass string) bool {
		return pass == "Password"
	})

	// set up root router with Logger, Recovery and LocalStorage middleware
	w := wok.Default()

	// Index page
	w.GET("/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		fmt.Fprint(w, "Index page")
		return nil
	})

	// api is a group of routes with common authentication, result rendering and error handling
	api := w.Group("/api", apiAuth, render.JSON)
	{
		api.GET("/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			res := []int{1, 2, 3, 4, 5}
			return render.Yield(ctx, 200, res)
		})
		api.GET("/:id", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			id := wok.Var(ctx, "id")
			res := struct {
				ID string
			}{id}
			return render.Yield(ctx, 201, res)
		})
	}

	// dash is an example of another separate route group
	dash := w.Group("/dash", dashboardAuth)
	{
		dash.GET("/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			fmt.Fprintf(w, "Hello %s", mw.GetUser(ctx))
			return nil
		})
	}

	http.ListenAndServe(":8080", w)
}

```
