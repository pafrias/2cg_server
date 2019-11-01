package trap

import (
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
		if err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
			return
		}

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

func (app *TCApp) getComponents() http.HandlerFunc {

	// A component is a trigger, target or effect of a trap
	type component struct {
		ID   int           `json:"_id,omitempty"`
		Name string        `json:"name"`
		Type string        `json:"type"`
		Text string        `json:"text,omitempty"`
		Cost sql.NullInt64 `json:"cost,omitempty"`
		P1   string        `json:"param1,omitempty"`
		P2   string        `json:"param2,omitempty"`
		P4   string        `json:"param4,omitempty"`
		P3   string        `json:"param3,omitempty"`
	}

	type shortComponent struct {
		ID   string `json:"_id,omitempty"`
		Name string `json:"name"`
		Type string `json:"type"`
	}

	return func(res http.ResponseWriter, req *http.Request) {

		reqType := mux.Vars(req)["type"]
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rows, err := app.db.GetComponents(ctx, reqType)
		if err != nil {
			// handle error
			return
		}
		defer rows.Close()

		var components []interface{}

		for rows.Next() {
			if reqType == "short" {
				var c shortComponent
				err = rows.Scan(&c.ID, &c.Name, &c.Type)
				if err != nil {
					println(err.Error())
				}
				components = append(components, c)
			} else {
				var c component
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

func test(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	form := req.Form
	for key, val := range form {
		fmt.Printf("key: %v, val: %v\n", key, val)
		for i, str := range val {
			fmt.Printf("\tindex: %v, string: %v\n", i, str)
		}
	}

}
