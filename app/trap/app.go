package trap

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-playground/form"
	"github.com/pafrias/2cgaming-api/db"
)

// Server FMI
type Server struct {
	DB *sql.DB
}

// NewServer returns a new instance of the trap API server
func NewServer(db *db.Connection) Server {
	return Server{db.Client}
}

func (s *Server) PrintForm() http.HandlerFunc {

	decoder := form.NewDecoder()

	return func(res http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
		}

		var upgrade upgrade

		if err := decoder.Decode(upgrade, req.Form); err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
		}

		form := req.PostForm
		for field, array := range form {
			fmt.Printf("field: %v\n", field)
			for index, val := range array {
				fmt.Printf("\tindex: [%v], val: %v\n", index, val)
			}
		}
		var response string
		for _, value := range form["type"] {
			response += value
		}

		res.Write([]byte(response))
	}
}
