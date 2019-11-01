package api

import (
	"app/db"
)

type App struct {
	db *db.Connection
}
