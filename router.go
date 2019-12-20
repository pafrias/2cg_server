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

	//Auth
	r.Path("/api/signin").HandlerFunc(s.signin()).Methods("GET")

	// create subrouters and delegate paths
	s.createTrapRouter(r.PathPrefix("/api/tc/").Subrouter())

	// serves ui
	r.PathPrefix("/").HandlerFunc(s.serveStaticFiles()).Methods("GET")

	// middleware
	r.Use(middleware.LogRequests)
	r.Use(middleware.ParseForm)
}

func (s *server) createTrapRouter(r *mux.Router) {
	api := trap.OpenService(&s.Connection)

	r.HandleFunc("/build/{budget}", api.HandleBuildTrap()).Methods("GET")

	r.HandleFunc("/components/last", api.GetLastUpdate("tc_component")).Methods("GET")
	r.HandleFunc("/components/{type}", api.GetComponents()).Methods("GET")
	r.HandleFunc("/components", api.GetComponents()).Methods("GET")
	// r.HandleFunc("/components/{type}", api.GetComponents()).Queries("fields", "").Methods("GET")
	r.HandleFunc("/components", s.checkAuth(2, api.PostComponent())).Methods("POST")
	// PATCH NEEDED

	r.HandleFunc("/upgrades/last", api.GetLastUpdate("tc_upgrade")).Methods("GET")
	r.HandleFunc("/upgrades", api.GetUpgrades()).Methods("GET")
	r.HandleFunc("/upgrades", s.checkAuth(2, api.PostUpgrade())).Methods("POST")
	// PATCH NEEDED
}

func (s *server) serveStaticFiles() http.HandlerFunc {
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
