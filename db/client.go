package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Connection struct {
	Client *sql.DB
}

func Open() *Connection {
	user := os.Getenv("SQL_USER")
	password := os.Getenv("SQL_PW")
	// needs extendability
	connectionString := fmt.Sprintf("%s:%s@/trap_compendium", user, password)

	conn := new(Connection)

	db, err := sql.Open("mysql", connectionString)
	conn.testFatalError(err, "Error connecting to database")

	err = db.Ping()
	conn.testFatalError(err, "Error pinging database")

	fmt.Println("Connected to database")
	conn.Client = db
	return conn
}

func (c *Connection) testFatalError(err error, message string) {
	if err != nil {
		fmt.Println(message + ":")
		log.Fatal(err.Error())
	}
}
