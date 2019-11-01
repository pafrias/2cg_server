package trap

import (
	"github.com/gorilla/mux"
)

func (a *TCApp) applyRoutes(r *mux.Router) {
	r.HandleFunc("/components", a.getComponents()).Methods("GET")
	r.HandleFunc("/components", a.postComponent()).Methods("POST")
	// PATCH
	r.HandleFunc("/components/{type}", a.getComponents()).Methods("GET")
	r.HandleFunc("/upgrades", a.postUpgrade()).Methods("POST")
	r.HandleFunc("/upgrades", a.getUpgrades()).Methods("GET")
}
