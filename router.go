package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pafrias/2cgaming-api/app/trap"

	"github.com/gorilla/mux"
	"github.com/pafrias/2cgaming-api/middleware"
)

func (s *server) createMainRouter() {
	r := mux.NewRouter()

	// create subrouters and delegate paths
	s.applyTrapCompendiumRoutes(r.PathPrefix("/api/tc/").Subrouter())

	// serves ui
	r.PathPrefix("/").HandlerFunc(s.serveStaticFiles()).Methods("GET")

	// middleware
	r.Use(middleware.LogRequests)
	r.Use(middleware.ParseForm)

	s.router = r
}

func (s *server) applyTrapCompendiumRoutes(r *mux.Router) {
	trapAPI := trap.NewHandler(s.db)
	r.HandleFunc("/test", trapAPI.PrintForm()).Methods("POST")

	r.HandleFunc("/components", trapAPI.GetComponents()).Methods("GET")
	r.HandleFunc("/components/{type}", trapAPI.GetComponents()).Methods("GET")
	r.HandleFunc("/components", trapAPI.PostComponent()).Methods("POST")
	// PATCH NEEDED

	r.HandleFunc("/upgrades", trapAPI.GetUpgrades()).Methods("GET")
	r.HandleFunc("/upgrades", trapAPI.PostUpgrade()).Methods("POST")
	// PATCH NEEDED
}

func (s *server) serveStaticFiles() http.HandlerFunc {
	dir := http.Dir("./web")
	fs := http.FileServer(dir)
	return func(res http.ResponseWriter, req *http.Request) {
		path := req.RequestURI
		if s.spaTest(path, "/trapcompendium/") {
			filePath := "./web/trapcompendium/"
			fileName, Type := s.readFileExtension(path)
			filePath += fileName
			fmt.Printf("redirected request from '%v' to '%v'\n", path, filePath)
			file, err := ioutil.ReadFile(filePath)
			if err != nil {
				fmt.Printf("something happened when trying to reach: %v\n", path)
			} else {
				if Type != "" {
					res.Header().Set("Content-Type", "text/"+Type)
				}
				res.Write(file)
			}
		} else {
			fs.ServeHTTP(res, req)
		}
	}
}

// helps routing for single page applications
func (s *server) spaTest(path string, prefix string) bool {
	if strings.HasPrefix(path, prefix) {
		fmt.Printf("prefix '%v' matches path '%v'\n", path, prefix)
		return true //strings.HasSuffix(path, "/") || strings.HasSuffix(path, "/index.html")
	}
	return false
}

func (s *server) readFileExtension(path string) (string, string) {
	pathVals := strings.Split(path, "/")
	switch file := pathVals[len(pathVals)-1]; {
	case file == "index.js" || file == "":
		return "index.html", "html"
	case file == "main.js":
		return "main.js", ""
	case strings.HasSuffix(file, ".css"):
		return "css/" + file, "css"
	default:
		return "media/" + file, ""
	}
}
