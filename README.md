# Noodle


[![Coverage](http://gocover.io/_badge/github.com/andviro/noodle?0)](http://gocover.io/github.com/andviro/noodle)  
[![GoDoc](http://godoc.org/github.com/andviro/noodle?status.png)](http://godoc.org/github.com/andviro/noodle)

Noodle is a tiny and (almost) unopinionated Golang middleware stack. It
borrows its ideas from [Stack](https://github.com/alexedwards/stack.git) 
package, but relies on Golang net 
[contexts](http://godoc.org/golang.org/x/net/context) for threading request 
environment through handler chains.

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
with additional middlewares as arguments. Each `Use()` call creates new 
middleware chain totally independent from parent. The following example extends 
root chain with variables from `gorilla/mux` router. For standalone example of 
`gorilla/mux` usage see [provided sample 
code](https://github.com/andviro/noodle/blob/master/examples/gorilla/main.go).

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
context variables and serves user requests. The resulting handler implements
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

## License and credits

This code is released under 
[MIT](https://github.com/andviro/noodle/blob/master/LICENSE) license.
