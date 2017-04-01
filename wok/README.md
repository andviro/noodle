# Wok

[![GoDoc](http://godoc.org/gopkg.in/andviro/noodle.v2/wok?status.png)](http://godoc.org/gopkg.in/andviro/noodle.v2/wok)

A simple and minimalistic (51 LOC in wok.go) web application router based on
[httprouter](https://github.com/julienschmidt/httprouter). Supports route
groups, global, per-group and per-route
[noodle](https://gopkg.in/andviro/noodle.v2) middleware. Compatible `http.HandlerFunc`.
For a quick start see the [sample application](https://github.com/andviro/noodle/blob/v2/examples/wok/main.go).

## Root router object

Wok router is created by `wok.Default()` and `wok.New()` constructors that
accept arbitrary list of noodle middlewares. Note that the `Default`
constructor preloads standard logger, recovery and local storage middlewares
that come with [noodle middleware](https://github.com/andviro/noodle/tree/v2/middleware)
package. Resulting middleware chain will be shared among all routes.

```go
w := wok.Default()
```


## Handling routes

Convenience methods `GET`, `POST`, `PUT`, `PATCH`, `DELETE` and `OPTIONS`
create routing entries with specific paths and middleware chains. Following is
an example of attaching hander to site root. Note that all of the methods
return a closure that accept a single `noodle.Handler` parameter.

```go

import "gopkg.in/andviro/noodle.v2/render"

func index(w http.ResponseWriter, r *http.Request) {
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
func userDetail(w http.ResponseWriter, r *http.Request) {
    id := wok.Var(r, "id")
    // ... do something with the id
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

## License

This code is released under
[MIT](https://github.com/andviro/noodle/blob/master/LICENSE) license.
