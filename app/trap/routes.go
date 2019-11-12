package trap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/form"
	"github.com/gorilla/mux"
)

var components map[uint16]component
var upgrades map[uint16]upgrade

// Extend to handle a request for any number of fields?
func (s *Server) GetComponents() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		reqType := mux.Vars(req)["type"]
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rows, err := s.readComponents(ctx, reqType)
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

		res.Header().Set("Content-Type", "application/json")
		res.Write(data)
	}
}

// PostComponent does things
func (s *Server) PostComponent() http.HandlerFunc {

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

		result, err := s.createComponent(req.Form)

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

func (s *Server) GetUpgrades() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rows, err := s.readUpgrades(ctx)
		if err != nil {
			// test error
			return
		}
		defer rows.Close()

		var upgrades []upgrade

		for rows.Next() {
			var u upgrade
			err = rows.Scan(&u.ID, &u.Name, &u.Type, &u.ComponentID, &u.Component, &u.Text, &u.Cost, &u.Max)
			if err != nil {
				println(err.Error())
			}
			upgrades = append(upgrades, u)
		}

		data, err := json.Marshal(upgrades)

		res.Header().Set("Content-Type", "application/json")
		res.Write(data)

	}
}

func (s *Server) PostUpgrade() http.HandlerFunc {

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

		result, err := s.createUpgrade(req.Form)

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
