package noodle

import (
	"bytes"
	"fmt"
	"golang.org/x/net/context"
	"gopkg.in/tylerb/is.v1"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func serveAndRequest(h http.Handler) string {
	ts := httptest.NewServer(h)
	defer ts.Close()
	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	resBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(resBody)
}

func mwFactory(name string) Middleware {
	return func(next Handler) Handler {
		varName := "Var" + name
		varValue := name + "value"
		return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
			withVars := context.WithValue(c, varName, varValue)
			fmt.Fprintf(w, "%s>", name)
			return next(withVars, w, r)
		}
	}
}

func handlerFactory(name string, keys ...string) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		fmt.Fprintf(w, "%s ", name)
		for _, key := range keys {
			fmt.Fprintf(w, "[%s=%v]", key, ctx.Value(key))
		}
		return nil
	}
}

func dumbMwFactory(name string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s?>", name)
			next.ServeHTTP(w, r)
		})
	}
}

func negroniMwFactory(name string) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		fmt.Fprintf(w, "%s?>", name)
		next(w, r)
	}
}

func TestNew(t *testing.T) {
	is := is.New(t)

	n := New(mwFactory("A"), mwFactory("B")).Then(handlerFactory("H1", "VarA", "VarB"))
	res := serveAndRequest(n)
	is.Equal("A>B>H1 [VarA=Avalue][VarB=Bvalue]", res)
}

func TestUse(t *testing.T) {
	is := is.New(t)

	n := New(mwFactory("A"), mwFactory("B")).Use(mwFactory("C")).Then(handlerFactory("H1", "VarA", "VarB", "VarC"))
	res := serveAndRequest(n)
	is.Equal("A>B>C>H1 [VarA=Avalue][VarB=Bvalue][VarC=Cvalue]", res)
}

func TestUseSeparates(t *testing.T) {
	is := is.New(t)

	root := New(mwFactory("A"), mwFactory("B"))
	chain1 := root.Use(mwFactory("C"))
	chain2 := root.Use(mwFactory("D"))
	h1 := root.Then(handlerFactory("H1", "VarA", "VarB"))
	h2 := chain1.Then(handlerFactory("H1", "VarA", "VarB", "VarC"))
	h3 := chain2.Then(handlerFactory("H1", "VarA", "VarB", "VarD"))
	res1 := serveAndRequest(h1)
	res2 := serveAndRequest(h2)
	res3 := serveAndRequest(h3)
	is.Equal("A>B>H1 [VarA=Avalue][VarB=Bvalue]", res1)
	is.Equal("A>B>C>H1 [VarA=Avalue][VarB=Bvalue][VarC=Cvalue]", res2)
	is.Equal("A>B>D>H1 [VarA=Avalue][VarB=Bvalue][VarD=Dvalue]", res3)
}

func TestThen(t *testing.T) {
	is := is.New(t)

	h := New().Then(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		fmt.Fprintf(w, "Abracadabra")
		return nil
	})
	is.Equal("Abracadabra", serveAndRequest(h))
}

func panickyHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "Abracadabra")
	panic("whoopsie!")
}

func runErrorPropagationTest(is *is.Is, extras ...Middleware) {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	h := Default(extras...).Then(panickyHandler)
	_ = serveAndRequest(h)
	log := buf.String()
	is.Equal("error = panic: whoopsie!", strings.TrimSpace(log[len(log)-25:]))
}

func TestDefaultMiddlewares(t *testing.T) {
	is := is.New(t)
	runErrorPropagationTest(is)
}

func runAdaptorTest(is *is.Is, converted Middleware) {
	root := New(mwFactory("A"), converted, mwFactory("C"))
	res := serveAndRequest(root.Then(handlerFactory("H1", "VarA", "VarC")))
	is.Equal("A>B?>C>H1 [VarA=Avalue][VarC=Cvalue]", res)
}

func TestAdaptContextKeeps(t *testing.T) {
	is := is.New(t)
	runAdaptorTest(is, Adapt(dumbMwFactory("B")))
}

func TestAdaptErrorPropagates(t *testing.T) {
	is := is.New(t)
	runErrorPropagationTest(is, Adapt(dumbMwFactory("A")))
}

func TestAdaptNegroniContextKeeps(t *testing.T) {
	is := is.New(t)
	runAdaptorTest(is, AdaptNegroni(negroniMwFactory("B")))
}

func TestAdaptNegroniErrorPropagates(t *testing.T) {
	is := is.New(t)
	runErrorPropagationTest(is, AdaptNegroni(negroniMwFactory("A")))
}
