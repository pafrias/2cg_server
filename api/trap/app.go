package trap

import (
	"github.com/gorilla/mux"
	"github.com/pafrias/2cgaming-api/db"
)

// A TCApp contains relevant request logic and model handlers for the trap
// compendium api
type TCApp struct {
	db *db.Connection
}

// SetUp instantiates a new trap compendium db connection and attaches route
// handlers to the main multiplexer
func SetUp(r *mux.Router) {
	var tc *TCApp
	tc.db = db.Open()
	tc.applyRoutes(r)
}

func (app *TCApp) applyRoutes(r *mux.Router) {
	r.HandleFunc("/test", test).Methods("GET")
	r.HandleFunc("/components", app.getComponents()).Methods("GET")
	r.HandleFunc("/components", app.postComponent()).Methods("POST")
	// PATCH
	r.HandleFunc("/components/{type}", app.getComponents()).Methods("GET")
	r.HandleFunc("/upgrades", app.postUpgrade()).Methods("POST")
	r.HandleFunc("/upgrades", app.getUpgrades()).Methods("GET")
}
