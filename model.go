package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func (s *server) createConnection() {
	user := os.Getenv("SQL_USER")
	password := os.Getenv("SQL_PW")
	connectionString := fmt.Sprintf("%s:%s@/trap_compendium", user, password)

	db, err := sql.Open("mysql", connectionString)
	s.testFatalError(err, "Error connecting to database")

	err = db.Ping()
	s.testFatalError(err, "Error pinging to database")

	s.db = db
	fmt.Println("Connected to database")
}

// A Component is a trigger, target or effect of a trap
type Component struct {
	ID     int    `json:"id" bson:"id"`
	Name   string `json:"name" bson:"name"`
	Type   int    `json:"type" bson:"type"`
	Text   string `json:"text" bson:"text"`
	Cost   int    `json:"cost,omitempty" bson:"cost,omitempty"`
	Param1 string `json:"param1,omitempty" bson:"param1,omitempty"`
	Param2 string `json:"param2,omitempty" bson:"param2,omitempty"`
	Param3 string `json:"param3,omitempty" bson:"param3,omitempty"`
	Param4 string `json:"param4,omitempty" bson:"param4,omitempty"`
}

// UpgradeType some some some.
type UpgradeType struct {
	Code int
	Name string
}
