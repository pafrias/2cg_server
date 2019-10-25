package main

import (
	"app/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (s *server) postComponent() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		values := parseComponentValues(req)

		result, err := s.db.PostComponent(values)

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

func (s *server) postUpgrade() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		values := parseUpgradeValues(req)

		result, err := s.db.PostUpgrade(values)

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

func (s *server) getComponents() http.HandlerFunc {

	type component struct {
		ID     int    `json:"_id,omitempty"`
		Name   string `json:"name"`
		Type   string `json:"type"`
		Text   string `json:"text"`
		Cost   int    `json:"cost,omitempty"`
		Param1 string `json:"param1,omitempty"`
		Param2 string `json:"param2,omitempty"`
		Param3 string `json:"param3,omitempty"`
		Param4 string `json:"param4,omitempty"`
	}

	return func(res http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rows, err := s.db.GetComponents(ctx)
		defer rows.Close()
		if err != nil {
			s.testQueryError(err, "something")
			return
		}

		var components []model.Component

		for rows.Next() {
			var c model.Component
			rows.Scan(&c)
			components = append(components, c)
			// fmt.Printf("%+v\n", c)
		}
		data, _ := json.Marshal(components)

		res.Write(data)
	}
}

func (s *server) serveStaticFiles() http.HandlerFunc {
	dir := http.Dir("./src/app/web")
	fs := http.FileServer(dir)
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("getting file -> ", req.URL)
		fs.ServeHTTP(res, req)
	}
}

func parseComponentValues(req *http.Request) (results []interface{}) {
	req.ParseForm()
	form := req.Form
	results = append(results, form.Get("name"), form.Get("type"), form.Get("text"), form.Get("cost"), form.Get("p1"), form.Get("p2"), form.Get("p3"), form.Get("p4"))
	return
}

func parseUpgradeValues(req *http.Request) (results []interface{}) {
	req.ParseForm()
	form := req.Form
	results = append(results, form.Get("name"), form.Get("type"), form.Get("text"), form.Get("cost"), form.Get("component_id"), form.Get("max"))
	return
}

// func scanToComponent(rows *sql.Rows) (results []interface{}) {

// 	for rows.Next() {
// 		c := Component{}
// 		rows.Scan(&c)
// 		results = append(results, c)
// 	}

// 	return results
// }
