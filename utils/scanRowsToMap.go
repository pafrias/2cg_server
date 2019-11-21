package utils

import (
	"database/sql"
	"fmt"
	"strconv"
)

func ScanRowsToArray(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, _ := rows.ColumnTypes()
	length := len(columns)
	result := []map[string]interface{}{}

	for rows.Next() {
		receiver := makeResultReceiver(length)
		if err := rows.Scan(receiver...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})

		for i := 0; i < length; i++ {
			col := columns[i]
			key := col.Name()
			val := *(receiver[i]).(*interface{})
			if val == nil {
				// fmt.Println(key, " value is null")
				// row[key] = nil
				continue
			}

			switch col.DatabaseTypeName() {
			case "INT":
				str := string(val.([]uint8))
				parsedVal, err := strconv.ParseInt(str, 0, 32)
				if err == nil {
					row[key] = int(parsedVal)
				} else {
					fmt.Println("error parsing ", key, " value")
					fmt.Println(err.Error())
				}
			case "VARCHAR":
				parsed := string(val.([]uint8))
				if !isEmpty(parsed) {
					row[key] = parsed
				}
			default:
				fmt.Printf("unsupport data type '%s' now\n", col.DatabaseTypeName())
				// TODO remember add other data type
			}
		}
		result = append(result, row)
	}
	return result, nil
}

func makeResultReceiver(length int) []interface{} {
	result := make([]interface{}, 0, length)
	for i := 0; i < length; i++ {
		var pointer interface{}
		pointer = struct{}{}
		result = append(result, &pointer)
	}
	return result
}

func isNull(val interface{}) bool {
	if val == nil {
		return true
	}
	return false
}

func isEmpty(val string) bool {
	if val == "" {
		return true
	}
	return false
}
