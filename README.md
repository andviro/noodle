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

- Simple and minimalistic: <30 LOC in "core.go"
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
and `adaptors` subpackage.

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

Middleware chain is finalized and converted to `noodle.Handler` with
`Then()` method. Its first parameter is an application handler that consumes 
context and serves user requests. The resulting handler implements
`http.Handler` interface providing `ServeHTTP` method. When serving HTTP from 
`noodle.Handler` default empty context is created for each request. For further 
flexibility `noodle.Handler` can be provided with externally created `context`. This 
advanced usage is outlined in [httprouter adaptor 
example](https://github.com/andviro/noodle/blob/master/examples/httprouter/main.go).

```go
func index(c context.Context, w http.ResponseWriter, r *http.Request) error {
    fmt.Fprintln(w, "index")
    return nil
}

...

http.Handle("/", n.Then(index))
```

## Logging and recovery

Package `noodle` comes with baked-in `Logger` and `Recovery` middlewares that 
provide just that, logging and recovery. For convenience, root `noodle.Chain` 
with baked-in logging and recovery can be created with `noodle.Default()` 
constructor. Default recovery middleware wraps the panic object into `error` 
and passes it further for logger middleware to display.

```go

func panickyIndex(c context.Context, w http.ResponseWriter, r *http.Request) error {
    ...
    panic("Oh noes!!!")
    ...
}

...

n := noodle.Default()
http.Handle("/", n.Then(panickyIndex))
```

## Compatibility with third-party middleware

`noodle.Adapt` converts generic middleware constructor with signature 
`func(http.Handler) http.Handler` to Noodle Middleware. Resulting constructor 
can be easily integrated into existing `noodle.Chain` with `Use` method. While 
converted middleware can not consume request context and is not able to return 
any error, context propagation is not broken and error values will bubble up 
from further handlers in chain. This allows usage of various middlewares 
written for third-party middleware libraries, like 
[interpose](https://github.com/carbocation/interpose).

```go

func DumbMid(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.RequestWriter, r *http.Request){
        fmt.Fprintf(w, "I'm dumb and proud of it!!!")
        next.ServeHTTP(w, r)
    })
}

...
// AwareMid2 will consume context from AwareMid1
n := noodle.New(AwareMid1, noodle.Adapt(DumbMid), AwareMid2)
http.Handle("/", n.Then(indexHandler))
```

`noodle.AdaptNegroni` creates `noodle.Middleware` from function with signature 
`func(http.ResponseWriter, *http.Request, http.HandlerFunc)`. This adaptor 
simplifies integration of middlewares written for 
[negroni](https://github.com/codegangsta/negroni)
package.


```go

func NegroniHandler (w http.RequestWriter, r *http.Request, next http.HandlerFunc) {
    fmt.Fprintf(w, "Hi, I'm a Negroni middleware!!!")
    next(w, r)
}

...
// AwareMid2 will consume context from AwareMid1
n := noodle.New(AwareMid1, noodle.AdaptNegroni(NegroniHandler), AwareMid2)
http.Handle("/", n.Then(indexHandler))
```
## License

This code is released under 
[MIT](https://github.com/andviro/noodle/blob/master/LICENSE) license.
