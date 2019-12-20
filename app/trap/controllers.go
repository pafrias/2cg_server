package trap

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"
)

func (s *Service) createComponent(v url.Values) (sql.Result, error) {
	if err := s.DB.Ping(); err != nil {
		return nil, err
	}

	var columns, fields []string
	var values []interface{}

	for column := range v {
		columns = append(columns, column)
		fields = append(fields, "?")
		values = append(values, v.Get(column))
	}

	query := fmt.Sprintf("insert into tc_component (%v) values (%v)", strings.Join(columns, ","), strings.Join(fields, ","))
	return s.DB.Exec(query, values...)
}

func (s *Service) createUpgrade(v url.Values) (sql.Result, error) {
	if err := s.DB.Ping(); err != nil {
		return nil, err
	}

	var columns, fields []string
	var values []interface{}

	for column := range v {
		columns = append(columns, column)
		fields = append(fields, "?")
		values = append(values, v.Get(column))
	}

	query := fmt.Sprintf("insert into tc_upgrade (%v) values (%v)", strings.Join(columns, ","), strings.Join(fields, ","))

	return s.DB.Exec(query, values...)
}

//expand to allow any number of fields
//params will require a join
//type will require a join
func (s *Service) readComponents(ctx context.Context, queryType string) (r *sql.Rows, err error) {
	if err := s.DB.Ping(); err != nil {
		return nil, err
	}

	var query string

	if queryType == "short" {
		query = `
			select id, c.name, ct.name as type
			from tc_component c
				inner join tc_comp_type ct on ct.code = c.type
		`
	} else if queryType == "build" {
		query = `
			select id, c.name, cost, param1 as costp, ct.name as type 
			from tc_component c
				inner join tc_comp_type ct on ct.code = c.type
		`
	} else {
		query = `
			select 
				id, c.name, ct.name as type, text, cost, param1, param2, param3, param4
			from tc_component c
				inner join tc_comp_type ct on ct.code = c.type
		`
	}

	return s.DB.QueryContext(ctx, query)

}

func (s *Service) readUpgrades(ctx context.Context, queryType string) (r *sql.Rows, err error) {
	if err := s.DB.Ping(); err != nil {
		return nil, err
	}
	var query string
	if queryType == "build" {
		query = `
			select
				u.id, u.name, ut.name as type, u.component_id, u.cost, u.max
			from tc_upgrade u
				inner join tc_up_type ut on ut.code = u.type`
	} else {
		query = `
			select
			u.id, u.name, ut.name as type, u.component_id, c.name as component, u.text, u.cost, u.max
			from tc_upgrade u
				left join tc_component c on u.component_id = c.id
				inner join tc_up_type ut on ut.code = u.type`
	}

	return s.DB.QueryContext(ctx, query)
}

func (s *Service) updateTimestamp(table string) (sql.Result, error) {
	result, err := s.DB.Exec("UPDATE update_times SET date = ? WHERE table_name = ?", time.Now(), table)
	fmt.Println(result.LastInsertId())
	fmt.Println(result.RowsAffected())

	if n, _ := result.RowsAffected(); n == 0 {
		result, err = s.DB.Exec("INSERT INTO update_times VALUES (?, ?)", table, time.Now())
	}
	fmt.Println(result.LastInsertId())
	fmt.Println(result.RowsAffected())

	return result, err
}

func (s *Service) readTimestamp(ctx context.Context, table string) (*sql.Row, error) {
	if err := s.DB.Ping(); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT date FROM update_times
		WHERE table_name = '%v';`, table)
	fmt.Println(query)
	return s.DB.QueryRowContext(ctx, query), nil

}
