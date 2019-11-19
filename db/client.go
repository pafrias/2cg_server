package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Connection struct {
	*sql.DB
}

/*Open will search the environment for the SQl_USER and SQL_PW variables, using them
to log into the local systems sql server.

Consider extending the functionality to connect to external databases and differerent
db names*/
func Open() Connection {
	user := os.Getenv("SQL_USER")
	password := os.Getenv("SQL_PW")

	connectionString := fmt.Sprintf("%s:%s@/trap_compendium", user, password)

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		fmt.Println("Error connecting to database:")
		log.Fatal(err.Error())
	}

	if err = db.Ping(); err != nil {
		fmt.Println("Error pinging database:")
		log.Fatal(err.Error())
	}

	fmt.Println("Connected to database")
	return Connection{db}
}

func (c *Connection) BuildTables() {
	if err := c.DB.Ping(); err != nil {
		// handle error
	}

	res, err := c.DB.Exec(USERSCHEMA)

	res, err = c.DB.Exec(TCSCHEMA)

	fmt.Printf("Response: %v\nError: %v", res, err)
}
