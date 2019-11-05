package main

import (
	"fmt"
	"net/http"

	"github.com/pafrias/2cgaming-api/app/trap"

	"github.com/gorilla/mux"
	"github.com/pafrias/2cgaming-api/middleware"
)

func (s *server) createMainRouter() {
	r := mux.NewRouter()

	// create subrouters and delegate paths
	s.applyTrapCompendiumRoutes(r.PathPrefix("/api/tc/").Subrouter())

	// serves ui
	r.PathPrefix("/").HandlerFunc(s.serveStaticFiles()).Methods("GET")

	// middleware
	r.Use(middleware.LogRequests)
	r.Use(middleware.ParseForm)

	s.router = r
}

func (s *server) applyTrapCompendiumRoutes(r *mux.Router) {
	trapAPI := trap.TCHandler{s.db}
	r.HandleFunc("/test", trapAPI.PrintForm()).Methods("POST")

	r.HandleFunc("/components", trapAPI.GetComponents()).Methods("GET")
	r.HandleFunc("/components/{type}", trapAPI.GetComponents()).Methods("GET")
	r.HandleFunc("/components", trapAPI.PostComponent()).Methods("POST")
	// PATCH NEEDED

	r.HandleFunc("/upgrades", trapAPI.GetUpgrades()).Methods("GET")
	r.HandleFunc("/upgrades", trapAPI.PostUpgrade()).Methods("POST")
	// PATCH NEEDED
}

// func spell codex api routes

func (s *server) serveStaticFiles() http.HandlerFunc {
	dir := http.Dir("./web")
	fs := http.FileServer(dir)
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("getting file -> ", req.URL)
		fs.ServeHTTP(res, req)
	}
}
