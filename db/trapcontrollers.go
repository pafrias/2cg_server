package db

import (
	"context"
	"database/sql"
)

func (c *Connection) PostComponent(values []interface{}) (sql.Result, error) {
	if err := c.Client.Ping(); err != nil {
		return nil, err
	}

	query := `
		insert into tc_component
			(name, type, text, cost, param1, param2, param3, param4)
		values
			(?,?,?,?,?,?,?,?)
	`

	return c.Client.Exec(query, values...)
}

func (c *Connection) PostUpgrade(values []interface{}) (sql.Result, error) {
	if err := c.Client.Ping(); err != nil {
		return nil, err
	}

	query := `
		insert into tc_upgrade
			(name, type, text, cost, component_id, max)
		values
			(?,?,?,?,?,?)
	`

	return c.Client.Exec(query, values...)
}

func (c *Connection) GetComponents(ctx context.Context, queryType string) (r *sql.Rows, err error) {
	if err := c.Client.Ping(); err != nil {
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

	return c.Client.QueryContext(ctx, query)

}

func (c *Connection) GetUpgrades(ctx context.Context) (r *sql.Rows, err error) {
	if err := c.Client.Ping(); err != nil {
		return nil, err
	}

	query := `select u.id,u.name,ut.name as type,c.name as component,u.text,u.cost,u.max
		from tc_upgrade u
		left join tc_component c on u.component_id = c.id
		inner join tc_up_type ut on ut.code = u.type`

	return c.Client.QueryContext(ctx, query)
}
