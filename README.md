# Noodle

[![Build Status](https://travis-ci.org/andviro/noodle.svg?branch=master)](https://travis-ci.org/andviro/noodle) 
[![Coverage](http://gocover.io/_badge/github.com/andviro/noodle?0)](http://gocover.io/github.com/andviro/noodle)  
[![GoDoc](http://godoc.org/github.com/andviro/noodle?status.png)](http://godoc.org/github.com/andviro/noodle)

Noodle is a tiny and (almost) unopinionated Golang middleware stack. It
borrows its ideas from [Stack](https://github.com/alexedwards/stack.git) 
package, but relies on Golang net 
[contexts](http://godoc.org/golang.org/x/net/context) for threading request 
environment through handler chains.

## Highlights

- Simple and minimalistic: <30 LOC in core package
- Strictly adheres to guidelines of [context](http://godoc.org/golang.org/x/net/context) package
- Noodle Handlers are context-aware and return error for easier error handling
- Finalized Noodle Handlers implement http.Handler interface, and easy to use 
  with routing library of choice
* Batteries included: contains collection of essential middlewares such as
  Logger, Recovery etc., all context-aware. 
* Includes adapter collection that allow integration of third-party
  middlewares without breaking of context and error propagation

## Middleware and handlers

`noodle.Middleware` is a generic `func(Handler) Handler` bidirectional 
middleware that accepts and returns `noodle.Handler`. Following is an example 
of HTTP basic auth middleware that stores user login in request `context`.

```go

// HTTPAuth is a middleware factory function that accepts `authFunc` for 
// username and password verification and returns `noodle.Middleware`
func HTTPAuth(authFunc func(username, password string) bool) noodle.Middleware 
{
    // The middleware function verifies user credentials and aborts chain 
    // execution if `authFunc` returns `false`
    return func(next noodle.Handler) noodle.Handler {
        return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
            username, password, ok := r.BasicAuth()
            if !ok || !authFunc(username, password) {
                w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
                w.WriteHeader(http.StatusUnauthorized)
                // Provide error for logging middleware then abort chain
                return fmt.Errorf("Unauthorized request")
            }
            // Inject user name into request context
            userContext := context.WithValue(c, "User", username)
            return next(userContext, w, r)
        }
    }
}
```

Note that similar HTTP Basic Auth middleware is included in noodle `middleware`
subpackage discussed below.

`noodle.Handler` provides context-aware `http.Handler` with `error` return 
value for enhanced chaining. Assuming that some middleware stored user login in 
request context, the following handler outputs personalized greeting:

```go
func index(c context.Context, w http.ResponseWriter, r *http.Request) error {
    user := c.Value("User").(string)
    fmt.Fprintf(w, "Hello %s", user)
    return nil
}
```

## Building noodle chains

Middleware chains are created with `noodle.New()`. Optionally middewares can be 
passed to this function to initalize chain.

```go
basicAuth := HTTPAuth(func(username, password string) bool {
    return password == "secret"
})

n := noodle.New(basicAuth)
```

At any moment `noodle.Chain` can be extended by calling `Use()` method with
some additional middlewares as arguments. Each `Use()` call creates new
middleware chain totally independent from parent. The following example extends
root chain with variables from `gorilla/mux` router. For standalone example of
`gorilla/mux` usage see [provided sample
code](https://github.com/andviro/noodle/blob/master/examples/gorilla/main.go)
and `adapt/gorilla` subpackage.

```go
func GorillaVars(next noodle.Handler) noodle.Handler {
    return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
        withVars := context.WithValue(c, "Vars", mux.Vars(r))
        return next(withVars, w, r)
    }
}

n = n.Use(GorillaVars)
```

## Handling HTTP requests

Middleware chain is finalized and converted to `noodle.Handler` with `Then()`
method. Its first parameter is an application handler that consumes context and
serves user requests. The resulting handler implements `http.Handler` interface
providing `ServeHTTP` method. When serving HTTP from `noodle.Handler` default
empty context is passed to each request. For further flexibility
`noodle.Handler` can be provided with externally created `context`. This
advanced usage is outlined in [httprouter adaptor
example](https://github.com/andviro/noodle/blob/master/examples/httprouter/main.go)
and put to use in `adapt/httprouter` subpackage.

```go
func index(c context.Context, w http.ResponseWriter, r *http.Request) error {
    fmt.Fprintln(w, "index")
    return nil
}

...

http.Handle("/", n.Then(index))
```

## Baked-in middlewares

Package `noodle` comes with a collection of essential middlewares organized
into the `middleware` package. Import "github.com/andviro/noodle/middleware" to
get access to following:

* Binding of request body to Golang structures
* HTTP Basic Auth middleware
* Logger
* Panic recovery
* Thread-safe request-global storage

`Logger` provides basic request logging with some useful info about response
timing. `Recovery` middleware wraps the panic object into `error` and passes it
further for logger middleware to display.


```go

import (
    "github.com/andviro/noodle/middleware"
)

func panickyIndex(c context.Context, w http.ResponseWriter, r *http.Request) error {
    ...
    panic("Oh noes!!!")
    ...
}

...

n := noodle.New(middleware.Logger, middleware.Recover)
http.Handle("/", n.Then(panickyIndex))
```

`LocalStore` middleware injects a thread-safe data store into the request
context. This store can be used to pass variables back and forth along the
middleware chain. Consider a handler that sets some variable in request-local
store and a middleware that wraps that handler and uses that variable after
handler execution to render output in JSON format:

```go

func RenderJSON(next noodle.Handler) noodle.Nandler {
    return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
        if err := next(c, w, r); err != nil {
            return err
        }
        // Expect some "data" value in local store
        data := middleware.GetStore(c).MustGet("data")
        json.NewEncoder(w).Encode(data)
        return nil
    }
}

func index(c context.Context, w http.ResponseWriter, r *http.Request) error {
    var res struct {
        A int
        B string
    }{1, "Ahahaha"}

    // This value will be caught by RenderJson middleware
    middleware.GetStore(c).Set("data", &res)
    return nil
}

...

n := noodle.New(middleware.LocalStore, RenderJSON)
http.Handle("/", n.Then(index))
```

For convenience, initial `noodle.Chain` with logging, recovery and
request-local store can be created with `middleware.Default()` constructor.

Refer to package [documentation](http://godoc.org/github.com/andviro/noodle/middleware) for
further information on provided middlewares.

## Compatibility with third-party middleware

Subpackage [adapt](http://godoc.org/github.com/andviro/noodle/adapt)
contains adaptors for third-party middleware libraries. `adapt.Http` converts
generic middleware constructor with signature `func(http.Handler) http.Handler`
to `noodle.Middleware`. Resulting constructor can be easily integrated into
existing `noodle.Chain` with `Use` method. While converted middleware can not
consume request context and is not able to return any error, context
propagation is not broken and error values will bubble up from further handlers
in chain. This allows usage of various middlewares written for third-party
middleware libraries, like
[interpose](https://github.com/carbocation/interpose).

```go

import (
    "github.com/andviro/noodle"
    "github.com/andviro/noodle/adapt"
)

func DumbMid(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.RequestWriter, r *http.Request){
        fmt.Fprintf(w, "I'm dumb and proud of it!!!")
        next.ServeHTTP(w, r)
    })
}

...
// AwareMid2 and indexHandler will consume context from AwareMid1
n := noodle.New(AwareMid1, adapt.Http(DumbMid), AwareMid2)
http.Handle("/", n.Then(indexHandler))
```

`adapt.Negroni` creates `noodle.Middleware` from function with signature
`func(http.ResponseWriter, *http.Request, http.HandlerFunc)`. This adaptor
simplifies integration of middlewares written for
[negroni](https://github.com/codegangsta/negroni) package.


```go

func NegroniHandler (w http.RequestWriter, r *http.Request, next http.HandlerFunc) {
    fmt.Fprintf(w, "Hi, I'm a Negroni middleware!!!")
    next(w, r)
}

...
// AwareMid2 will consume context from AwareMid1
n := noodle.New(AwareMid1, adapt.Negroni(NegroniHandler), AwareMid2)
http.Handle("/", n.Then(indexHandler))
```

## Convenience adaptors

For compatibility with popular router packages, such as [Gorilla
mux](https://github.com/gorilla/mux) and
[httprouter](https://github.com/julienschmidt/httprouter) corresponding
[middleware](http://godoc.org/github.com/andviro/noodle/adapt/gorilla) and
[adaptor struct](http://godoc.org/github.com/andviro/noodle/adapt/httprouter)
are included.

## License

This code is released under 
[MIT](https://github.com/andviro/noodle/blob/master/LICENSE) license.
