package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gorilla/mux"
)

type server struct {
	db     *sql.DB
	router *mux.Router
	// auth something figured out later
}

func (s *server) testFatalError(e error, message string) {
	if e != nil {
		fmt.Println(message)
		log.Fatal(e.Error())
	}
}

func (s *server) testQueryError(e error, message string) {
	if e != nil {
		fmt.Println(message, e.Error())
	}
}
