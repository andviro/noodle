package wok_test

import (
	"fmt"
	"github.com/andviro/noodle"
	"github.com/andviro/wok"
	"golang.org/x/net/context"
	"gopkg.in/tylerb/is.v1"
	"net/http"
	"net/http/httptest"
	"testing"
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
		fmt.Fprintf(w, "[%s]", tag)
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
	wk.GET("/", handlerFactory("B"))
	is.Equal(testRequest(wk, "GET", "/"), "A>[B]")
}

func TestGroup(t *testing.T) {
	is := is.New(t)
	wk := wok.New(mwFactory("A"))
	g1 := wk.Group("/g1", mwFactory("G1"))
	g2 := wk.Group("/g2", mwFactory("G2"))
	g1.GET("/", noodle.New(mwFactory("G11")).Then(handlerFactory("B")))
	g2.GET("/", noodle.New(mwFactory("G21")).Then(handlerFactory("C")))

	is.Equal(testRequest(wk, "GET", "/g1/"), "A>G1>G11>[B]")
	is.Equal(testRequest(wk, "GET", "/g2/"), "A>G2>G21>[C]")
}

func TestVars(t *testing.T) {
	is := is.New(t)
	mw := func(next noodle.Handler) noodle.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			fmt.Fprint(w, "MW>")
			return next(context.WithValue(ctx, 0, "testValue"), w, r)
		}
	}
	wk := wok.New(mw)
	wk.GET("/:varA/:varB", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		fmt.Fprintf(w, "[%s][%s][%s]", wok.Var(ctx, "varA"), wok.Var(ctx, "varB"), ctx.Value(0).(string))
		return nil
	})
	is.Equal(testRequest(wk, "GET", "/1/2"), "MW>[1][2][testValue]")
}
