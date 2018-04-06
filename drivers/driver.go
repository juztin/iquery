package drivers

import (
	"database/sql"
	"fmt"

	"github.com/juztin/iquery"
)

type Driver struct {
	Info
	mapValue MapValueFunc
}

type Info interface {
	Name() string
	DriverName() string
	ConnectionString(s iquery.Server) string
	SchemasQuery() string
	TablesQuery(schema string) string
	ColumnsQuery(schema, table string) string
}

type MapValueFunc func(o interface{}) interface{}
type MapBytesFunc func([]byte) interface{}

func mapBytes(b []byte) interface{} {
	s := string(b)
	return &s
}

func ValueMapper(mapBytesFn MapBytesFunc) MapValueFunc {
	return func(o interface{}) interface{} {
		var value interface{}
		switch o.(type) {
		case nil:
			value = nil
		case []byte:
			value = mapBytesFn(o.([]byte))
		case bool:
			value = o.(bool)
		case string:
			value = o.(string)
		case int:
			value = o.(int)
		case int16:
			value = o.(int16)
		case int32:
			value = o.(int32)
		case int64:
			value = o.(int64)
		case float32:
			value = o.(float32)
		case float64:
			value = o.(float64)
		default:
			s := fmt.Sprintf("%v", o)
			value = &s
		}
		return value
	}
}

func (d Driver) MapRows(r *sql.Rows, columns, limit int) (iquery.Rows, error) {
	var rows iquery.Rows

	// Value/Value-Pointer collections
	values := make([]interface{}, columns)
	pointers := make([]interface{}, columns)
	for i := 0; i < columns; i++ {
		pointers[i] = &values[i]
	}

	// Fetch Rows
	i := 0
	for r.Next() {
		// Stop once we've hit the requested row limit.
		if i == limit {
			break
		}
		i++

		// Scan values into pointer slice.
		err := r.Scan(pointers...)
		if err != nil {
			return rows, err
		}

		// Format the pointer values
		row := []interface{}{}
		for _, o := range values {
			row = append(row, d.mapValue(o))
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func (d Driver) Schemas(s iquery.Server) ([]string, error) {
	db, err := sql.Open(d.DriverName(), d.ConnectionString(s))
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(d.SchemasQuery())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dbs []string
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return dbs, err
		}
		dbs = append(dbs, name)
	}
	return dbs, err
}

func (d Driver) Tables(s iquery.Server, schema string) ([]iquery.Table, error) {
	db, err := sql.Open(d.DriverName(), d.ConnectionString(s))
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(d.TablesQuery(schema))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []iquery.Table
	for rows.Next() {
		var t iquery.Table
		err = rows.Scan(&t.Name, &t.Description)
		if err != nil {
			return tables, err
		}
		tables = append(tables, t)
	}
	return tables, err
}

func (d Driver) Columns(s iquery.Server, schema, table string) ([]iquery.Column, error) {
	db, err := sql.Open(d.DriverName(), d.ConnectionString(s))
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(d.ColumnsQuery(schema, table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []iquery.Column
	for rows.Next() {
		var c iquery.Column
		var length sql.NullInt64
		err = rows.Scan(&c.Name, &c.Description, &c.Type, &c.Default, &length, &c.IsNullable)
		if err != nil {
			return columns, err
		}
		if length.Valid {
			c.Type = fmt.Sprintf("%s(%d)", c.Type, length.Int64)
		}
		columns = append(columns, c)
	}
	err = rows.Err()
	return columns, err
}

func New(i Info, fn MapValueFunc) Driver {
	if fn == nil {
		fn = ValueMapper(mapBytes)
	}
	return Driver{i, fn}
}
