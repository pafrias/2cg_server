package trap

import (
	"app/api"
	"app/db"

	"github.com/gorilla/mux"
)

type TCApp struct {
	api.App
}

func SetUp(r *mux.Router) {
	var tc *TCApp
	tc.db = db.Open()
	tc.applyRoutes(r)
}
