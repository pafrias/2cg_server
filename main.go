package main

import (
	"app/model"
	"net/http"
)

func main() {
	// create server instance
	s := server{}
	s.db = model.CreateConnection()

	// open to requests
	http.ListenAndServe(":3001", s.routes())

}
