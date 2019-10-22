package main

import (
	"net/http"
)

func main() {
	// create server instance
	s := server{}
	s.createConnection()
	s.routes()
	// open to requests
	http.ListenAndServe(":3001", s.router)

}
