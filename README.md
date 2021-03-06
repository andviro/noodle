# Noodle

## Migration to the new version

The release of Go 1.7 greatly simplified use of contexts in HTTP handlers.
Further development of Noodle is continued in separate Github organization
[go-noodle](https://github.com/go-noodle), making this project and its branches obsolete. It is currently
being maintained only to avoid breakage of old code. Please address your pull
requests and issues to the new and improved [Noodle
project](https://github.com/go-noodle/noodle).

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
* Comes with a [minimalistic web application framework](https://github.com/andviro/noodle/tree/master/wok)
  based on [httprouter](https://github.com/julienschmidt/httprouter) that has route
  groups and supports global, per-route and per-group noodle middleware.

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

Note that similar [HTTP Basic Auth middleware](https://godoc.org/github.com/andviro/noodle/middleware#HTTPAuth) is included in `middleware` subpackage discussed below.

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
`gorilla/mux` usage see [provided sample code](https://github.com/andviro/noodle/blob/master/examples/gorilla/main.go)
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
advanced usage is outlined in
[httprouter adaptor example](https://github.com/andviro/noodle/blob/master/examples/httprouter/main.go)
and put to use in `wok` router .

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

* HTTP Basic Auth middleware
* Logger
* Panic recovery
* Thread-safe request-global storage


### Logger and recovery

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


### LocalStore

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

## Rendering of handler results

Package [render](http://godoc.org/github.com/andviro/noodle/render) provides
basic middleware for serialization of handler-supplied values, similar to the
example above. The only difference is that handler must call `render.Yield`
function to pass HTTP status code and its data back to the render middleware
through context. Currently supported are `render.JSON` middleware for JSON
serialization, `render.XML` for XML output and `render.Template` that uses
pre-compiled `html/template` object to render data object into HTML.

```go
import (
    mw "github.com/andviro/noodle/middleware"
    "github.com/andviro/noodle/render"
    "html/template"
)


type TestStruct struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func index(c context.Context, w http.ResponseWriter, r *http.Request) error {
	testData := TestStruct{1, "Ohohoho"}
    return render.Yield(c, 201, &testData)
})

...

n := mw.Default()
http.Handle("/jsonEndpoint", n.Use(render.JSON).Then(index))

tpl, _ := template.New().Parse("<h1>Hello {{ .B }}</h1> {{ .A }}")
http.Handle("/htmlEndpoint", n.Use(render.Template(tpl)).Then(index))
```

## Rendering based on Accept header

`render.ContentType` analyzes request's `Accept` header and renders output data
to JSON, XML or HTML, setting `ContentType` appropriately. If template for HTML
rendering is `nil`, result is rendered to indented JSON inside HTML `PRE` tag,
which can be used for endpoint debugging. If `Accept` header is not specified
or not recognized, the result is rendered to JSON.

```go

tpl, _ := template.Must(template.New().Parse("<h1>Hello {{ .B }}</h1> {{ .A }}"))

// or
// tpl := nil

http.Handle("/anyContent", n.Use(render.ContentType(tpl)).Then(index))
```

## Request binding 

Package [bind](http://godoc.org/github.com/andviro/noodle/bind) provides
middleware for loading request body into supplied model. Handlers retrieve
bound objects using `bind.GetData` function.

```go
import (
    mw "github.com/andviro/noodle/middleware"
    "github.com/andviro/noodle/bind"
)

type TestStruct struct {
	A int    `json:"a" form:"a"`
	B string `json:"b" form:"b"`
}

func index(c context.Context, w http.ResponseWriter, r *http.Request) error {
    data := bind.GetData(c).(*TestStruct) // Get parsed data from context
    // Use model data
    ...
    return nil
})

...

n := mw.Default()
// The following handler will bind request body to TestStruct type
http.Handle("/jsonPostEndpoint", n.Use(bind.JSON(TestStruct{})).Then(index))

// The following handler will bind post form to TestStruct type
http.Handle("/formPostEndpoint", n.Use(bind.Form(TestStruct{})).Then(index))
```

Currently binding of JSON and web forms through
[agj/form](https://github.com/ajg/form) library is supported. XML etc is work
in progress, PRs are appreciated.

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

For compatibility with Gorilla [mux](https://github.com/gorilla/mux)  corresponding
[middleware](http://godoc.org/github.com/andviro/noodle/adapt/gorilla) is provided.

## License

This code is released under 
[MIT](https://github.com/andviro/noodle/blob/master/LICENSE) license.
