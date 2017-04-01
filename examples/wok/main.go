package main

import (
	mw "gopkg.in/andviro/noodle.v2/middleware"
	"gopkg.in/andviro/noodle.v2/render"
	"gopkg.in/andviro/noodle.v2/wok"
	"html/template"
	"net/http"
)

// API is a simple service that provides some HTTP handler funcs
type API struct{}

// Auth callback for HTTPAuth middleware
func (*API) Auth(user, pass string) bool {
	return pass == "Secret"
}

// Index handler for API
func (*API) Index(w http.ResponseWriter, r *http.Request) {
	res := []int{1, 2, 3, 4, 5}
	render.Yield(r, 200, res)
}

// Detail on some id passed in URL
func (*API) Detail(w http.ResponseWriter, r *http.Request) {
	id := wok.Var(r, "id")
	res := struct {
		ID string
	}{id}
	render.Yield(r, 201, res)
}

// DashBoard is another service
type DashBoard struct{}

// Auth callback for HTTPAuth middleware
func (*DashBoard) Auth(user, pass string) bool {
	return pass == "Password"
}

// Index page
func (*DashBoard) Index(w http.ResponseWriter, r *http.Request) {
	res := map[string]interface{}{
		"User": mw.GetUser(r),
	}
	render.Yield(r, 201, res)
}

// Index is a generic HTTP handler function, it relies on template rendering and does nothing
func Index(w http.ResponseWriter, r *http.Request) {}

var (
	// Load and parse templates
	indexTpl     = template.Must(template.New("index").Parse("<h1>Hello</h1>"))
	dashboardTpl = template.Must(template.New("dash").Parse("<h1>Hello {{ .User }}</h1>"))
)

func main() {
	// set up root router with Logger, Recovery and LocalStorage middleware
	w := wok.Default()

	// Index page
	w.GET("/", render.Template(indexTpl))(Index)

	// API route group
	api := new(API)
	apiGrp := w.Group("/api", mw.HTTPAuth("API", api.Auth), render.JSON)
	apiGrp.GET("/")(api.Index)
	apiGrp.GET("/:id")(api.Detail)

	// Dashboard route group
	dash := new(DashBoard)
	dashGrp := w.Group("/dash", mw.HTTPAuth("Dashboard", dash.Auth))
	dashGrp.GET("/", render.Template(dashboardTpl))(dash.Index)

	http.ListenAndServe(":8080", w)
}
