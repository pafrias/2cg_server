package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pafrias/2cgaming-api/api/trap"
)

func main() {
	router := mux.NewRouter()

	// set up api connections and routes
	trap.SetUp(router.PathPrefix("/api/tc/").Subrouter())

	// serve static files
	// needs a fix for SPAs
	router.PathPrefix("/").HandlerFunc(serveStaticFiles()).Methods("GET")

	// middleware
	router.Use(logRequests)

	// open to requests
	http.ListenAndServe(":3001", router)

}

func serveStaticFiles() http.HandlerFunc {
	dir := http.Dir("./src/app/web")
	fs := http.FileServer(dir)
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("getting file -> ", req.URL)
		fs.ServeHTTP(res, req)
	}
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Println(req.Method, req.URL)
		next.ServeHTTP(res, req)
	})
}
