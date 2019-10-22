package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *server) routes() {
	r := mux.NewRouter()
	// r.HandleFunc("/api/triggers", s.handleGetTriggers()).Methods("GET")
	// r.HandleFunc("/api/targets", s.handleGetTargets()).Methods("GET")
	// r.HandleFunc("/api/components", s.handleGetComponents()).Methods("GET")
	r.HandleFunc("/api/tc/ut", s.getUpgradeTypes()).Methods("GET")
	r.HandleFunc("/api/tc/components", s.postComponent()).Methods("POST")
	r.HandleFunc("/api/tc/components", s.getComponent()).Methods("GET")
	r.PathPrefix("/").HandlerFunc(s.serveStaticFiles()).Methods("GET")
	r.Use(logReq)
	s.router = r
}

func logReq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Println(req.Method, req.URL)
		next.ServeHTTP(res, req)
	})
}
