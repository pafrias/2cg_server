package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pafrias/2cgaming-api/app/trap"
	"github.com/pafrias/2cgaming-api/middleware"
	"github.com/pafrias/2cgaming-api/utils"
)

func (s *server) createMainRouter() {
	r := s.Router

	//Auth, requires heavy patching
	r.Path("/api/signin").HandlerFunc(s.signin()).Methods("GET")

	// create subrouters and delegate paths
	s.routeTrapCompendium(r.PathPrefix("/api/tc/").Subrouter())
	// s.routeEpicSpellCodex(r.PathPrefix("/api/epic/spell_codex").Subrouter())

	r.PathPrefix("/").HandlerFunc(s.serveFrontEnd()).Methods("GET")

	// middleware
	r.Use(middleware.LogRequests)
	r.Use(middleware.ParseForm)
}

func (s *server) serveFrontEnd() http.HandlerFunc {
	dir := http.Dir("./web")
	fs := http.FileServer(dir)

	tcBaseURL := "/trapcompendium/"
	tcRedirect := utils.HandleRedirect(tcBaseURL)

	return func(res http.ResponseWriter, req *http.Request) {
		path := req.RequestURI

		if utils.RequiresRedirect(path, tcBaseURL) {
			fmt.Println("redirecting")
			tcRedirect.ServeHTTP(res, req)
		} else {
			fs.ServeHTTP(res, req)
		}
	}
}

func (s *server) routeTrapCompendium(r *mux.Router) {
	api := trap.OpenService(&s.Connection)

	// should be query param?
	r.HandleFunc("/random/{budget}", api.HandleBuildTrap()).Methods("GET")

	r.HandleFunc("/components/last", api.GetLastUpdate("tc_component")).Methods("GET")
	r.HandleFunc("/components/{type}", api.GetComponents()).Methods("GET")
	r.HandleFunc("/components", api.GetComponents()).Methods("GET")
	// r.HandleFunc("/components/{type}", api.GetComponents()).Queries("fields", "").Methods("GET")
	r.HandleFunc("/components", s.checkAuth(2, api.PostComponent())).Methods("POST")
	// PATCH NEEDED

	r.HandleFunc("/upgrades/last", api.GetLastUpdate("tc_upgrade")).Methods("GET")
	r.HandleFunc("/upgrades", api.GetUpgrades()).Methods("GET")
	r.HandleFunc("/upgrades", s.checkAuth(2, api.PostUpgrade())).Methods("POST")

	// route may be removed to use at runtime
	r.HandleFunc("/upgrades/load", api.LoadUpgrades()).Methods("POST")
	// PATCH NEEDED
}

/*func (s *server) routeEpicSpellCodex (r *mux.Router) {
 // Fill me in
}*/
