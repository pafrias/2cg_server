package main

import (
	"fmt"
	"log"

	"github.com/pafrias/2cgaming-api/db"

	"github.com/gorilla/mux"
)

type server struct {
	db     *db.Connection
	router *mux.Router
}

func testFatalError(e error, message string) {
	if e != nil {
		fmt.Println(message)
		log.Fatal(e.Error())
	}
}

func testQueryError(e error, message string) {
	if e != nil {
		fmt.Println(message, e.Error())
	}
}
