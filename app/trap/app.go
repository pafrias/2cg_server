package trap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pafrias/2cgaming-api/db"

	"github.com/go-playground/form"
)

//App is the
type App struct {
	*db.Connection
	components
	upgrades
}

// NewHandler returns a new instance of the trap API server
func NewHandler(db *db.Connection) App {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	if err := db.PingContext(cxt); err != nil {
		//handle it
	}
	cancel()

	return App{db, components{}, upgrades{}}
}

func (a *App) PrintForm() http.HandlerFunc {

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

func (a *App) Test() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Println(a.components)
		fmt.Println(a.upgrades)

	}
}
