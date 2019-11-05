package main

import (
	"net/http"

	"github.com/pafrias/2cgaming-api/db"
)

func main() {

	s := server{}
	// test database connection
	s.db = db.Open()
	// test error

	s.createMainRouter()

	// listen to requests
	http.ListenAndServe(":3001", s.router)

}
