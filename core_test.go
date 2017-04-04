package noodle_test

import (
	"context"
	"fmt"
	"gopkg.in/andviro/noodle.v2"
	"gopkg.in/tylerb/is.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

// RunHTTP imitates handling HTTP request, returns response body
func RunHTTP(h http.Handler) string {
	r, _ := http.NewRequest("GET", "http://localhost", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w.Body.String()
}

func mwFactory(name string) noodle.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		varName := "Var" + name
		varValue := name + "value"
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s>", name)
			next(w, noodle.WithValue(r, varName, varValue))
		}
	}
}

func handlerFactory(name string, keys ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s ", name)
		for _, key := range keys {
			fmt.Fprintf(w, "[%s=%v]", key, noodle.Value(r, key))
		}
	}
}

func TestNew(t *testing.T) {
	is := is.New(t)

	n := noodle.New(mwFactory("A"), mwFactory("B")).Then(handlerFactory("H1", "VarA", "VarB"))
	res := RunHTTP(n)
	is.Equal("A>B>H1 [VarA=Avalue][VarB=Bvalue]", res)
}

func TestUse(t *testing.T) {
	is := is.New(t)

	n := noodle.New(mwFactory("A"), mwFactory("B")).Use(mwFactory("C")).Then(handlerFactory("H1", "VarA", "VarB", "VarC"))
	res := RunHTTP(n)
	is.Equal("A>B>C>H1 [VarA=Avalue][VarB=Bvalue][VarC=Cvalue]", res)
}

func TestUseSeparates(t *testing.T) {
	is := is.New(t)

	root := noodle.New(mwFactory("A"), mwFactory("B"))
	chain1 := root.Use(mwFactory("C"))
	chain2 := root.Use(mwFactory("D"))
	h1 := root.Then(handlerFactory("H1", "VarA", "VarB"))
	h2 := chain1.Then(handlerFactory("H1", "VarA", "VarB", "VarC"))
	h3 := chain2.Then(handlerFactory("H1", "VarA", "VarB", "VarD"))
	res1 := RunHTTP(h1)
	res2 := RunHTTP(h2)
	res3 := RunHTTP(h3)
	is.Equal("A>B>H1 [VarA=Avalue][VarB=Bvalue]", res1)
	is.Equal("A>B>C>H1 [VarA=Avalue][VarB=Bvalue][VarC=Cvalue]", res2)
	is.Equal("A>B>D>H1 [VarA=Avalue][VarB=Bvalue][VarD=Dvalue]", res3)
}

func TestThen(t *testing.T) {
	is := is.New(t)

	h := noodle.New().Then(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Abracadabra")
	})
	is.Equal("Abracadabra", RunHTTP(h))
}

func TestWrap(t *testing.T) {
	is := is.New(t)

	r, _ := http.NewRequest("GET", "http://localhost", nil)
	w := httptest.NewRecorder()
	ctx := context.TODO()
	ctx = noodle.Wrap(ctx, w, r)
	is.NotNil(ctx)
	w1, r1 := noodle.Unwrap(ctx)
	is.Equal(r, r1)
	is.Equal(w, w1)
}
