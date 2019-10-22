package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (s *server) serveStaticFiles() http.HandlerFunc {
	dir := http.Dir("./src/app/web")
	fs := http.FileServer(dir)
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("getting file -> ", req.URL)
		fs.ServeHTTP(res, req)
	}
}

func (s *server) postComponent() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		req.ParseForm()
		fmt.Printf("%v, %T\n", req.Form, req.Form)
		fmt.Printf("%v, %T\n", req.Form["name"], req.Form["name"])
		res.Write([]byte("yo"))
	}
}

func (s *server) getComponent() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("yo"))
	}
}

func (s *server) getUpgradeTypes() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		err := s.db.Ping()
		s.testFatalError(err, "Error connecting to database")

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		rows, err := s.db.QueryContext(ctx, "select * from tc_up_type")
		s.testQueryError(err, "Error querying upgrade type table")
		defer rows.Close()

		var upgradeTypes []UpgradeType

		for rows.Next() {
			var u UpgradeType
			rows.Scan(&u.Code, &u.Name)
			fmt.Printf("%+v\n", u)
			upgradeTypes = append(upgradeTypes, u)
		}
		data, _ := json.Marshal(upgradeTypes)

		res.Write(data)
	}
}
