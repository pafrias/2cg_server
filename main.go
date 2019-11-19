package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pafrias/2cgaming-api/db"
)

func main() {

	s := server{db.Open(), mux.NewRouter()}

	s.createMainRouter()

	// listen to requests
	http.ListenAndServe(":3001", s.router)

}
