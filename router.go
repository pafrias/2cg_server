package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *server) routes() *mux.Router {
	r := mux.NewRouter()
	s.trapCompendiumAPI(r.PathPrefix("/api/tc").Subrouter())
	//s.epicSpellAPI(r.PathPrefix("/api/epic/spells").Subrouter())
	r.PathPrefix("/").HandlerFunc(s.serveStaticFiles()).Methods("GET")
	r.Use(logReq)
	return r
}

func (s *server) trapCompendiumAPI(r *mux.Router) {
	r.HandleFunc("/components", s.getComponents()).Methods("GET")
	r.HandleFunc("/upgrades", s.postUpgrade()).Methods("POST")
	// PATCH
	r.HandleFunc("/components", s.postComponent()).Methods("POST")
	// PATCH
}

func (s *server) epicSpellAPI(r *mux.Router) {
	// r.HandleFunc("/", s.getELSpells()).Methods("GET")
	// r.HandleFunc("/{spell_id}", s.getELSpells()).Methods("GET")
	// POST
	// PATCH
	// DELETE
}

func logReq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Println(req.Method, req.URL)
		next.ServeHTTP(res, req)
	})
}
