package trap

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pafrias/2cgaming-api/db/models"
	"github.com/pafrias/2cgaming-api/utils"

	"github.com/go-playground/form"
	"github.com/gorilla/mux"
)

/*GetLastUpdate takes in a table name and returns a HandlerFunc which replies with a
timestamp for the last time that table was updated.*/
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

//Consider extending GetComponents to search the query parameters for the fields to return

/*GetComponents fetches data for trigger, targetting, and effect components*/
func (s *Service) GetComponents() http.HandlerFunc {
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

	type Component struct {
		Name string                `form:",required"`
		Type models.JsonNullString `form:",required"`
		Text string                `form:",required"`
		Cost models.JsonNullInt32  `form:",omitempty"`
		P1   models.JsonNullString
		P2   models.JsonNullString
		P4   models.JsonNullString
		P3   models.JsonNullString
	}

	decoder := form.NewDecoder()

	return func(res http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if s.HandleUnprocessableEntity(err, res) {
			return
		}

		var component Component
		err = decoder.Decode(&component, req.Form)
		if s.HandleUnprocessableEntity(err, res) {
			return
		}

		result, err := s.createComponent(req.Form)

		if !s.HandleUnprocessableEntity(err, res) {
			num, _ := result.LastInsertId()
			str := fmt.Sprintf("Success!\nComponent #%v inserted", num)
			s.updateTimestamp("tc_component")
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(str))
		}
	}
}

//Consider extending GetUpgrades to search the query parameters for the fields to return

/*GetUpgrades fetches all trap compendium upgrades*/
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
		err := req.ParseForm()
		if s.HandleUnprocessableEntity(err, res) {
			return
		}

		var upgrade Upgrade
		err = decoder.Decode(&upgrade, req.Form)
		if s.HandleUnprocessableEntity(err, res) {
			return
		}

		result, err := s.createUpgrade(upgrade)

		if !s.HandleUnprocessableEntity(err, res) {
			num, _ := result.LastInsertId()
			str := fmt.Sprintf("Success!\nUpgrade #%v inserted", num)
			s.updateTimestamp("tc_upgrade")
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(str))
		}
	}
}

/*HandleBuildTrap returns a HandlerFunc that generates a random trap with the given cost*/
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

/* s3 -> database

If this project implements Docker, maybe this function could be used later at runtime
in order to load the database dynamically, allowing for horizontal scaling*/

/*LoadUpgrades is for bulk loading of upgrade data.
Use cases include: database failures, horizontal scaling.*/
func (s *Service) LoadUpgrades() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {
		// change to make a request to a aws bucket
		// should not rely on locally stored files.
		jsonFile, err := os.Open("./upgrades.json")
		if s.HandleInternalServerError(err, res) {
			fmt.Println(err)
			return
		}
		defer jsonFile.Close()

		data, _ := ioutil.ReadAll(jsonFile)

		upgrades := []Upgrade{}

		err = json.Unmarshal(data, &upgrades)

		if s.HandleInternalServerError(err, res) {
			fmt.Println(err)
			return
		}

		for _, u := range upgrades {
			_, err = s.createUpgrade(u)
			if s.HandleUnprocessableEntity(err, res) {
				return
			}
		}

	}
}

/* Database -> s3
cacheUpgrades should be on a regular schedule, called every month to create a backup of TC relevant data
to it's respective s3 bucket

maybe this function shouldn't be a route handler, but simply a function called inside OpenService when
intializing the microservice.

func (this *Service) cacheUpgrades() http.HandlerFunc{

	return func (res http.ResponseWriter, req *http.Request) {
		// request all components
		// Marshall data
		// store in S3

		// request all upgrades
		// Marshall data
		// store in S3
	}
}
*/
