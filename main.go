package main

import (
	"app/api/trap"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	// set up trap compendium database and routes
	trap.SetUp(router.PathPrefix("/api/tc/").Subrouter())
	r.PathPrefix("/").HandlerFunc(serveStaticFiles()).Methods("GET")
	r.Use(logRequests)

	// open to requests
	http.ListenAndServe(":3001", s.router)

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
