package noodle_test

import (
	"context"
	"fmt"
	"github.com/andviro/noodle"
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
	return func(next noodle.Handler) noodle.Handler {
		varName := "Var" + name
		varValue := name + "value"
		return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
			withVars := context.WithValue(c, varName, varValue)
			fmt.Fprintf(w, "%s>", name)
			return next(withVars, w, r)
		}
	}
}

func handlerFactory(name string, keys ...string) noodle.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		fmt.Fprintf(w, "%s ", name)
		for _, key := range keys {
			fmt.Fprintf(w, "[%s=%v]", key, ctx.Value(key))
		}
		return nil
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

	h := noodle.New().Then(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		fmt.Fprintf(w, "Abracadabra")
		return nil
	})
	is.Equal("Abracadabra", RunHTTP(h))
}
