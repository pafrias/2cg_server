package trap

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
)

func (s *Server) createComponent(v url.Values) (sql.Result, error) {
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

func (s *Server) createUpgrade(v url.Values) (sql.Result, error) {
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

func (s *Server) readComponents(ctx context.Context, queryType string) (r *sql.Rows, err error) {
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

func (s *Server) readUpgrades(ctx context.Context) (r *sql.Rows, err error) {
	if err := s.DB.Ping(); err != nil {
		return nil, err
	}

	query := `select u.id,u.name,ut.name as type,u.component_id, c.name as component,u.text,u.cost,u.max
		from tc_upgrade u
		left join tc_component c on u.component_id = c.id
		inner join tc_up_type ut on ut.code = u.type`

	return s.DB.QueryContext(ctx, query)
}
