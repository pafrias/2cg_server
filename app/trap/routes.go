package trap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/form"
	"github.com/gorilla/mux"

	"github.com/pafrias/2cgaming-api/db/models"
)

// Extend to handle a request for any number of fields?
func (h *Handler) GetComponents() http.HandlerFunc {

	type shortComponent struct {
		ID   string `json:"_id,omitempty"`
		Name string `json:"name"`
		Type string `json:"type"`
	}

	return func(res http.ResponseWriter, req *http.Request) {

		reqType := mux.Vars(req)["type"]
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rows, err := h.DB.GetComponents(ctx, reqType)
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
				var c models.Component
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
func (h *Handler) PostComponent() http.HandlerFunc {

	decoder := form.NewDecoder()

	return func(res http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
		}

		var component models.Component

		if err := decoder.Decode(&component, req.Form); err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
		}

		result, err := h.DB.PostComponent(req.Form)

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

func (h *Handler) GetUpgrades() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rows, err := h.DB.GetUpgrades(ctx)
		if err != nil {
			// test error
			return
		}
		defer rows.Close()

		var upgrades []models.Upgrade

		for rows.Next() {
			var u models.Upgrade
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

func (h *Handler) PostUpgrade() http.HandlerFunc {

	decoder := form.NewDecoder()

	return func(res http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
			return
		}

		var upgrade models.Upgrade

		if err := decoder.Decode(upgrade, req.Form); err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte(err.Error()))
			return
		}

		result, err := h.DB.PostUpgrade(req.Form)

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
