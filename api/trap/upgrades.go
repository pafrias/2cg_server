package trap

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (app *TCApp) postUpgrade() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		values, err := parseUpgradeValues(req)
		if err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
			return
		}

		result, err := app.db.PostUpgrade(values)

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

func (app *TCApp) getUpgrades() http.HandlerFunc {

	// An Upgrade is a trigger, target or effect of a trap
	type upgrade struct {
		ID        int            `json:"_id,omitempty"`
		Name      string         `json:"name"`
		Type      string         `json:"type"`
		Component sql.NullString `json:"component,omitempty"`
		Text      string         `json:"text"`
		Cost      int            `json:"cost"`
		Max       int            `json:"max"`
	}

	return func(res http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rows, err := app.db.GetUpgrades(ctx)
		if err != nil {
			// test error
			return
		}
		defer rows.Close()

		var upgrades []upgrade

		for rows.Next() {
			var u upgrade
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
