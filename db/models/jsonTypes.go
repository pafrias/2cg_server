package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

//These types have cleaner marshalling behavior for null values

type JsonNullString struct {
	sql.NullString `json:",omitempty"`
}

func (v JsonNullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

type JsonNullInt32 struct {
	sql.NullInt32
}

func (v JsonNullInt32) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int32)
	} else {
		return json.Marshal(nil)
	}
}

type JsonNullInt64 struct {
	sql.NullInt64
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	fmt.Println("attempting to marshal int64")
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}
