package router

import (
	"api/src/router/routes"

	"github.com/gorilla/mux"
)

//Generate return new router with configs routes
func Generate() *mux.Router {
	r := mux.NewRouter()
	return routes.Configure(r)
}
