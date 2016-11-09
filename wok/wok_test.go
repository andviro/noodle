package wok_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/andviro/noodle"
	mw "github.com/andviro/noodle/middleware"
	"github.com/andviro/noodle/render"
	"github.com/andviro/noodle/wok"
	"gopkg.in/tylerb/is.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func mwFactory(tag string) noodle.Middleware {
	return func(next noodle.Handler) noodle.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			fmt.Fprintf(w, "%s>", tag)
			return next(ctx, w, r)
		}
	}
}

func handlerFactory(tag string) noodle.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		if val, ok := ctx.Value(tag).(string); ok {
			fmt.Fprintf(w, "[%s]", val)
		} else {
			fmt.Fprintf(w, "[%s]", tag)
		}
		return nil
	}
}

func TestNew(t *testing.T) {
	is := is.New(t)
	wok := wok.New()
	is.NotNil(wok)
}

func testRequest(wok *wok.Wok, method, path string) string {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(method, "http://localhost"+path, nil)
	wok.ServeHTTP(w, r)
	return w.Body.String()
}

func TestHandle(t *testing.T) {
	is := is.New(t)
	wk := wok.New(mwFactory("A"))
	wk.GET("/")(handlerFactory("B"))
	is.Equal(testRequest(wk, "GET", "/"), "A>[B]")
}

func TestGroup(t *testing.T) {
	is := is.New(t)
	wk := wok.New(mwFactory("A"))
	g1 := wk.Group("/g1", mwFactory("G1"))
	g2 := wk.Group("/g2", mwFactory("G2"))
	g1.GET("/", mwFactory("G11"))(handlerFactory("B"))
	g2.GET("/", mwFactory("G21"))(handlerFactory("C"))

	is.Equal(testRequest(wk, "GET", "/g1"), "A>G1>G11>[B]")
	is.Equal(testRequest(wk, "GET", "/g2"), "A>G2>G21>[C]")
}

func TestRouterVars(t *testing.T) {
	is := is.New(t)
	mw := func(next noodle.Handler) noodle.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			fmt.Fprint(w, "MW>")
			return next(context.WithValue(ctx, 0, "testValue"), w, r)
		}
	}
	wk := wok.New(mw)
	wk.GET("/:varA/:varB")(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		fmt.Fprintf(w, "[%s][%s][%s]", wok.Var(ctx, "varA"), wok.Var(ctx, "varB"), ctx.Value(0).(string))
		return nil
	})
	is.Equal(testRequest(wk, "GET", "/1/2"), "MW>[1][2][testValue]")
}

func ExampleApplication() {
	// globalErrorHandler receives all errors from all handlers
	// and tries to return meaningful HTTP status and message
	globalErrorHandler := func(next noodle.Handler) noodle.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			err := next(ctx, w, r)
			switch err {
			case mw.UnauthorizedRequest:
				w.WriteHeader(401)
				fmt.Fprint(w, "Please provide credentials")
			case nil:
				return nil
			default:
				w.WriteHeader(500)
				fmt.Fprintf(w, "There was an error: %v", err)
			}
			return err
		}
	}

	// apiErrorHandler is a specific error catcher that renders its messages into JSON
	apiErrorHandler := func(next noodle.Handler) noodle.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			err := next(ctx, w, r)
			switch err {
			case mw.UnauthorizedRequest:
				render.Yield(ctx, 401, map[string]interface{}{
					"error": err,
				})
			case nil:
				return nil
			default:
				render.Yield(ctx, 500, map[string]interface{}{
					"error": err,
				})
			}
			return err
		}
	}

	// apiAuth guards access to api group
	apiAuth := mw.HTTPAuth("API", func(user, pass string) bool {
		return pass == "Secret"
	})
	// apiAuth guards access to dashboard group
	dashboardAuth := mw.HTTPAuth("Dashboard", func(user, pass string) bool {
		return pass == "Password"
	})

	// w is the root router
	w := wok.Default(globalErrorHandler)

	// Handle index page
	w.GET("/")(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		fmt.Fprint(w, "Index page")
		return nil
	})

	// api is a group of routes with common authentication, result rendering and error handling
	api := w.Group("/api", render.JSON, apiErrorHandler, apiAuth)
	{
		api.GET("/")(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			res := []int{1, 2, 3, 4, 5}
			return render.Yield(ctx, 200, res)
		})
		api.GET("/:id")(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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
		dash.GET("/")(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			fmt.Fprintf(w, "Hello %s", mw.GetUser(ctx))
			return nil
		})
	}

	go http.ListenAndServe(":8989", w)
	time.Sleep(300 * time.Millisecond) // let it settle down

	// Here we will test webapp responses

	// index
	resp, _ := http.Get("http://localhost:8989/")
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println(string(data))

	// dashboard
	resp, _ = http.Get("http://user:Password@localhost:8989/dash")
	data, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println(string(data))

	// api index
	resp, _ = http.Get("http://user:Secret@localhost:8989/api")
	var lst []int
	json.NewDecoder(resp.Body).Decode(&lst)
	fmt.Println(lst)

	// api with parameter
	resp, _ = http.Get("http://user:Secret@localhost:8989/api/12")
	var obj map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&obj)
	fmt.Println(obj)

	// Output: Index page
	// Hello user
	// [1 2 3 4 5]
	// map[ID:12]
}
