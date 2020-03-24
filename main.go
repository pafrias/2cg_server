package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pafrias/2cgaming-api/db"
)

func main() {

	s := server{db.Open(), mux.NewRouter()}

	s.createMainRouter()

	// listen to requests
	err := http.ListenAndServe(":3001", s.Router)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
