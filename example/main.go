package main

import (
	mw "github.com/andviro/noodle/middleware"
	"github.com/andviro/noodle/render"
	"github.com/andviro/wok"
	"golang.org/x/net/context"
	"html/template"
	"net/http"
)

func index(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// nothing to do here, everything is in the template
	return nil
}

func apiIndex(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	res := []int{1, 2, 3, 4, 5}
	return render.Yield(ctx, 200, res)
}

func apiDetail(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := wok.Var(ctx, "id")
	res := struct {
		ID string
	}{id}
	return render.Yield(ctx, 201, res)
}

func dashIndex(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	res := map[string]interface{}{
		"User": mw.GetUser(ctx),
	}
	return render.Yield(ctx, 201, res)
}

func main() {
	// apiAuth guards access to api group
	apiAuth := mw.HTTPAuth("API", func(user, pass string) bool {
		return pass == "Secret"
	})
	// dashboardAuth guards access to dashboard group
	dashboardAuth := mw.HTTPAuth("Dashboard", func(user, pass string) bool {
		return pass == "Password"
	})

	// set up root router with Logger, Recovery and LocalStorage middleware
	w := wok.Default()

	// Index page
	idxTpl := template.Must(template.New("index").Parse("<h1>Hello</h1>"))
	w.GET("/", render.Template(idxTpl))(index)

	// api is a group of routes with common authentication and result rendering
	api := w.Group("/api", apiAuth, render.JSON)
	{
		api.GET("/")(apiIndex)
		api.GET("/:id")(apiDetail)
	}

	// dash is an example of another separate route group
	dash := w.Group("/dash", dashboardAuth)
	{
		tpl, _ := template.New("dash").Parse("<h1>Hello {{ .User }}</h1>")
		dash.GET("/", render.Template(tpl))(dashIndex)
	}

	http.ListenAndServe(":8080", w)
}
