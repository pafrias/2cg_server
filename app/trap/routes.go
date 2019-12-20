package trap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pafrias/2cgaming-api/utils"

	"github.com/go-playground/form"
	"github.com/gorilla/mux"
)

//GetLastUpdate is good.
func (s *Service) GetLastUpdate(table string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		row, err := s.readTimestamp(ctx, table)
		if s.HandleInternalServerError(err, res) {
			return
		}
		//fmt.Printf("%+v\n", row)

		var str string
		err = row.Scan(&str)
		if s.HandleInternalServerError(err, res) {
			return
		}
		fmt.Printf("%v\n", str)

		res.WriteHeader(http.StatusOK)
		res.Write([]byte(str))
	}
}

/*GetComponents fetches data for trigger, targetting, and effect components*/
func (s *Service) GetComponents() http.HandlerFunc {
	// Extend to handle a request for any number of fields?
	return func(res http.ResponseWriter, req *http.Request) {
		reqType := mux.Vars(req)["type"]
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rows, err := s.readComponents(ctx, reqType)
		if s.HandleInternalServerError(err, res) {
			return
		}
		defer rows.Close()

		components, err := utils.ScanRowsToArray(rows)
		if s.HandleInternalServerError(err, res) {
			return
		}
		// map ids ?

		data, err := json.Marshal(components)
		if s.HandleInternalServerError(err, res) {
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		res.Write(data)
	}
}

/*PostComponent posts components.
Relies upon auth at the server level to protect routes*/
func (s *Service) PostComponent() http.HandlerFunc {

	decoder := form.NewDecoder()

	return func(res http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
		}

		var component postComponent

		if err := decoder.Decode(&component, req.Form); err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
		}

		result, err := s.createComponent(req.Form)

		if err != nil {
			// handle different db errors ?
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
		} else {
			num, _ := result.LastInsertId()
			str := fmt.Sprintf("Success!\nComponent #%v inserted", num)
			s.updateTimestamp("tc_component")
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(str))
		}
	}
}

/*GetUpgrades gets all trap compendium upgrades*/
func (s *Service) GetUpgrades() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rows, err := s.readUpgrades(ctx, "")
		if s.HandleInternalServerError(err, res) {
			return
		}
		defer rows.Close()

		upgrades, err := utils.ScanRowsToArray(rows)
		if s.HandleInternalServerError(err, res) {
			return
		}

		data, err := json.Marshal(upgrades)

		if s.HandleInternalServerError(err, res) {
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		res.Write(data)
	}
}

/*PostUpgrade posts upgrades.
Relies upon auth at the server level to protect routes*/
func (s *Service) PostUpgrade() http.HandlerFunc {

	decoder := form.NewDecoder()

	return func(res http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
			return
		}

		var upgrade postUpgrade

		if err := decoder.Decode(&upgrade, req.Form); err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
			return
		}

		result, err := s.createUpgrade(req.Form)

		if err != nil {
			fmt.Println(err.Error())
			// handle db connect error
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
		} else {
			num, _ := result.LastInsertId()
			str := fmt.Sprintf("Success!\nUpgrade #%v inserted", num)
			s.updateTimestamp("tc_upgrade")
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(str))
		}
	}
}

//HandleBuildTrap blahl bhlahblahb
func (s *Service) HandleBuildTrap() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {
		budget, err := strconv.ParseInt(mux.Vars(req)["budget"], 0, 0)
		if s.HandleInternalServerError(err, res) {
			return
		} else if budget < 1 || budget > 200 {
			res.WriteHeader(400)
			res.Write([]byte("Cannot build a trap with a budget of 0"))
			return
		}

		ctx := context.TODO()
		components, err := s.readComponents(ctx, "build")
		if s.HandleInternalServerError(err, res) {
			return
		}

		upgrades, err := s.readUpgrades(ctx, "build")
		if s.HandleInternalServerError(err, res) {
			return
		}

		trap, err := buildRandomizedTrap(components, upgrades, int(budget))
		if s.HandleInternalServerError(err, res) {
			return
		}

		data, err := json.Marshal(&trap)
		if s.HandleInternalServerError(err, res) {
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		res.Write(data)
	}
}
