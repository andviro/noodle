package wok_test

import (
	"encoding/json"
	"fmt"
	"gopkg.in/andviro/noodle.v2"
	mw "gopkg.in/andviro/noodle.v2/middleware"
	"gopkg.in/andviro/noodle.v2/render"
	"gopkg.in/andviro/noodle.v2/wok"
	"gopkg.in/tylerb/is.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func mwFactory(tag string) noodle.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s>", tag)
			next(w, r)
		}
	}
}

func handlerFactory(tag string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if val, ok := noodle.Value(r, tag).(string); ok {
			fmt.Fprintf(w, "[%s]", val)
		} else {
			fmt.Fprintf(w, "[%s]", tag)
		}
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
	mw := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "MW>")
			next(w, noodle.WithValue(r, 0, "testValue"))
		}
	}
	wk := wok.New(mw)
	wk.GET("/:varA/:varB")(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "[%s][%s][%s]", wok.Var(r, "varA"), wok.Var(r, "varB"), noodle.Value(r, 0).(string))
	})
	is.Equal(testRequest(wk, "GET", "/1/2"), "MW>[1][2][testValue]")
}

func ExampleApplication() {
	// globalErrorHandler receives all errors from all handlers
	// and tries to return meaningful HTTP status and message
	globalErrorHandler := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r)
			//switch err {
			//case mw.UnauthorizedRequest:
			//w.WriteHeader(401)
			//fmt.Fprint(w, "Please provide credentials")
			//case nil:
			//return
			//default:
			//w.WriteHeader(500)
			//fmt.Fprintf(w, "There was an error: %v", err)
			//}
			//return err
		}
	}

	// apiErrorHandler is a specific error catcher that renders its messages into JSON
	apiErrorHandler := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r)
			//switch err {
			//case mw.UnauthorizedRequest:
			//render.Yield(ctx, 401, map[string]interface{}{
			//"error": err,
			//})
			//case nil:
			//return nil
			//default:
			//render.Yield(ctx, 500, map[string]interface{}{
			//"error": err,
			//})
			//}
			//return err
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
	w.GET("/")(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Index page")
	})

	// api is a group of routes with common authentication, result rendering and error handling
	api := w.Group("/api", render.JSON, apiErrorHandler, apiAuth)
	{
		api.GET("/")(func(w http.ResponseWriter, r *http.Request) {
			res := []int{1, 2, 3, 4, 5}
			render.Yield(r, 200, res)
		})
		api.GET("/:id")(func(w http.ResponseWriter, r *http.Request) {
			id := wok.Var(r, "id")
			res := struct {
				ID string
			}{id}
			render.Yield(r, 201, res)
		})
	}

	// dash is an example of another separate route group
	dash := w.Group("/dash", dashboardAuth)
	{
		dash.GET("/")(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello %s", mw.GetUser(r))
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
