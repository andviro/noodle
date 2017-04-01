package main

import (
	mw "gopkg.in/andviro/noodle.v2/middleware"
	"gopkg.in/andviro/noodle.v2/render"
	"gopkg.in/andviro/noodle.v2/wok"
	"html/template"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	// nothing to do here, everything is in the template
}

func apiIndex(w http.ResponseWriter, r *http.Request) {
	res := []int{1, 2, 3, 4, 5}
	render.Yield(r, 200, res)
}

func apiDetail(w http.ResponseWriter, r *http.Request) {
	id := wok.Var(r, "id")
	res := struct {
		ID string
	}{id}
	render.Yield(r, 201, res)
}

func dashIndex(w http.ResponseWriter, r *http.Request) {
	res := map[string]interface{}{
		"User": mw.GetUser(r),
	}
	render.Yield(r, 201, res)
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
