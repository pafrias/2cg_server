package trap

import (
	"app/model"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (app *TCApp) postComponent() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		values, err := parseComponentValues(req)
		// handle parsing error

		result, err := app.db.PostComponent(values)

		if err != nil {
			fmt.Println(err.Error())
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
		} else {
			num, _ := result.RowsAffected()
			str := fmt.Sprintf("Success!\n%v rows inserted", num)
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(str))
		}
	}
}

func (app *TCApp) postUpgrade() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		values, err := parseUpgradeValues(req)
		// handle parsing error

		result, err := app.db.PostUpgrade(values)

		if err != nil {
			fmt.Println(err.Error())
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
		} else {
			num, _ := result.RowsAffected()
			str := fmt.Sprintf("Success!\n%v rows inserted", num)
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(str))
		}
	}
}

func (app *TCApp) getComponents() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		reqType := mux.Vars(req)["type"]
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rows, err := app.db.GetComponents(ctx, reqType)
		if err != nil {
			app.testQueryError(err, "something")
			return
		}
		defer rows.Close()

		var components []interface{}

		for rows.Next() {
			if reqType == "short" {
				var c model.ShortComponent
				err = rows.Scan(&c.ID, &c.Name, &c.Type)
				if err != nil {
					println(err.Error())
				}
				components = append(components, c)
			} else {
				var c model.Component
				err = rows.Scan(&c.ID, &c.Name, &c.Type, &c.Text, &c.Cost, &c.P1, &c.P2, &c.P3, &c.P4)
				if err != nil {
					println(err.Error())
				}
				components = append(components, c)
			}
		}

		data, _ := json.Marshal(components)

		res.WriteHeader(200)
		res.Write(data)
	}
}

func (app *TCApp) getUpgrades() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rows, err := app.db.GetUpgrades(ctx)
		if err != nil {
			app.testQueryError(err, "something")
			return
		}
		defer rows.Close()

		var upgrades []model.Upgrade

		for rows.Next() {
			var u model.Upgrade
			err = rows.Scan(&u.ID, &u.Name, &u.Type, &u.Component, &u.Text, &u.Cost, &u.Max)
			if err != nil {
				println(err.Error())
			}
			upgrades = append(upgrades, u)
		}

		data, err := json.Marshal(upgrades)

		res.WriteHeader(200)
		res.Write(data)

	}
}

// this logic may be better located in the "model" package
func parseComponentValues(req *http.Request) ([]interface{}, error) {
	var results []interface{}
	req.ParseForm()
	form := req.Form
	var cost sql.NullInt64

	if form.Get("cost") != "" {
		x, err := strconv.ParseInt(form.Get("cost"), 10, 64)
		if err != nil {
			return results, err
		}
		cost.Int64 = x
		cost.Valid = true
	}

	results = append(results, form.Get("name"), form.Get("type"), form.Get("text"), cost, form.Get("param1"), form.Get("param2"), form.Get("param3"), form.Get("param4"))
	return results, nil
}

func parseUpgradeValues(req *http.Request) ([]interface{}, error) {
	var results []interface{}
	req.ParseForm()
	form := req.Form
	var componentID sql.NullInt64

	if form.Get("component_id") != "" {
		x, err := strconv.ParseInt(form.Get("component_id"), 10, 64)
		if err != nil {
			return results, err
		}
		componentID.Int64 = x
		componentID.Valid = true
	}
	results = append(results, form.Get("name"), form.Get("type"), form.Get("text"), form.Get("cost"), componentID, form.Get("max"))
	return results, nil
}

// func scanToComponent(rows *sql.Rows) (results []interface{}) {

// 	for rows.Next() {
// 		c := Component{}
// 		rows.Scan(&c)
// 		results = append(results, c)
// 	}

// 	return results
// }
