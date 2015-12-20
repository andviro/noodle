# Wok

[![Build
Status](https://travis-ci.org/andviro/wok.svg?branch=master)](https://travis-ci.org/andviro/wok)
[![Coverage](http://gocover.io/_badge/github.com/andviro/wok?0)](http://gocover.io/github.com/andviro/wok)
[![GoDoc](http://godoc.org/github.com/andviro/wok?status.png)](http://godoc.org/github.com/andviro/wok)

A simple and minimalistic (51 LOC in wok.go) web application router based on
[httprouter](https://github.com/julienschmidt/httprouter). Supports route
groups, global, per-group and per-route
[noodle](https://github.com/andviro/noodle) middleware. For a quick start see
the [sample application](https://github.com/andviro/wok/blob/master/example/main.go).

## Root router object

Wok router is created by `wok.Default()` and `wok.New()` constructors that
accept arbitrary list of noodle middlewares. Note that the `Default`
constructor preloads standard logger, recovery and local storage middlewares
that come with [noodle middleware](https://github.com/andviro/noodle/tree/master/middleware)
package. Resulting middleware chain will be shared among all routes.

```go
func errorHandler(next noodle.Handler) noodle.Handler {
    return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
        err := next(c, w, r)
        if err != nil {
            // do something with it
        }
        return err // pass error to logger middleware
    }
}

// error handler will catch all errors from routes
w := wok.Default(errorHandler)
```

## Handling routes

Convenience methods `GET`, `POST`, `PUT`, `PATCH`, `DELETE` and `OPTIONS`
create routing entries with specific paths and middleware chains. Following is
an example of attaching hander to site root. Note that all of the methods
return a closure that accept a single `noodle.Handler` parameter.

```go

import "github.com/andviro/noodle/render"

func index(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// nothing to do here, everything is in the template
	return nil
}

...

idxTpl := template.Must(template.New("index").Parse("<h1>Hello</h1>"))
// Middlewares are passed as variadic parameter list to GET method.
// Resulting closure accepts a route handler parameter.
w.GET("/", render.Template(idxTpl))(index)
```

Named parameters such as `/:name` and catch-all parameters i.e. `/*pathList`
are supported in route path assignment. See the
[httprouter documentation](http://godoc.org/github.com/julienschmidt/httprouter) for the
parameter syntax reference. To get the value of a route parameter use `wok.Var`
function:

```go
func userDetail(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
    id := wok.Var(ctx, "id")
    // ... do something with the id
	return nil
}
```

## Route grouping

`Group` method creates a route group with the specific prefix. A middleware
variadic list can be supplied to `Group` function, then the resulting
middleware chain for a group will contain a router global middlewares *and*
group-specific middlewares.

```go

// apiAuth is some group specific middleware
// apiIndex and apiDetail are route handers for /api and /api/:id paths
api := w.Group("/api", apiAuth, render.JSON)
api.GET("/")(apiIndex)
api.GET("/:id")(apiDetail)
```

Note that you also can pass route-specific middleware lists to `GET` methods!


## Serving HTTP

The `Wok` router object implements `http.Handler` interface and can be directly
passed to `http.ListenAndServe` function.

```go 
w := wok.Default()
// setup routes and middlewares
// ...

// start server
http.ListenAndServe(":8080", w)
```
