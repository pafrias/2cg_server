package trap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pafrias/2cgaming-api/utils"

	"github.com/go-playground/form"
	"github.com/gorilla/mux"
)

/*GetComponents fetches data for trigger, targetting, and effect components
Extend to handle a request for any number of fields?*/
func (a *App) GetComponents() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {
		reqType := mux.Vars(req)["type"]
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rows, err := a.readComponents(ctx, reqType)
		if a.Test500Error(err, res) {
			return
		}
		defer rows.Close()

		components, err := utils.ScanRowsToArray(rows)
		if a.Test500Error(err, res) {
			return
		}
		// map ids ?

		data, err := json.Marshal(components)
		if a.Test500Error(err, res) {
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		res.Write(data)
	}
}

/*PostComponent posts components.
Relies upon auth at the server level to protect routes*/
func (a *App) PostComponent() http.HandlerFunc {

	decoder := form.NewDecoder()

	return func(res http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
		}

		var component component

		if err := decoder.Decode(&component, req.Form); err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
		}

		result, err := a.createComponent(req.Form)

		if err != nil {
			// handle different db errors ?
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
		} else {
			num, _ := result.LastInsertId()
			str := fmt.Sprintf("Success!\nUpdate #%v inserted", num)
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(str))
		}
	}
}

/*GetUpgrades gets all trap compendium upgrades*/
func (a *App) GetUpgrades() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rows, err := a.readUpgrades(ctx, "")
		if a.Test500Error(err, res) {
			return
		}
		defer rows.Close()

		upgrades, err := utils.ScanRowsToArray(rows)
		if a.Test500Error(err, res) {
			return
		}

		data, err := json.Marshal(upgrades)

		if a.Test500Error(err, res) {
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		res.Write(data)
	}
}

/*PostUpgrade posts upgrades.
Relies upon auth at the server level to protect routes*/
func (a *App) PostUpgrade() http.HandlerFunc {

	decoder := form.NewDecoder()

	return func(res http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
			return
		}

		var upgrade upgrade

		if err := decoder.Decode(upgrade, req.Form); err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
			return
		}

		result, err := a.createUpgrade(req.Form)

		if err != nil {
			fmt.Println(err.Error())
			// handle db connect error
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
		} else {
			num, _ := result.LastInsertId()
			str := fmt.Sprintf("Success!\nUpdate #%v inserted", num)
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(str))
		}
	}
}

//HandleBuildTrap blahl bhlahblahb
func (a *App) HandleBuildTrap() http.HandlerFunc {

	trapBuilder := createBuilderFunc()

	return func(res http.ResponseWriter, req *http.Request) {
		ctx := context.TODO()
		rows, err := a.readComponents(ctx, "build")
		if a.Test500Error(err, res) {
			return
		}

		components, err := utils.ScanRowsToArray(rows)
		if a.Test500Error(err, res) {
			return
		}

		rows, err = a.readUpgrades(ctx, "build")
		if a.Test500Error(err, res) {
			return
		}

		upgrades, err := utils.ScanRowsToArray(rows)
		if a.Test500Error(err, res) {
			return
		}

		trap, err := trapBuilder(components, upgrades, 70)
		if a.Test500Error(err, res) {
			return
		}

		data, err := json.Marshal(trap)
		if a.Test500Error(err, res) {
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		res.Write(data)
		fmt.Println("Trap Complete")
	}
}
