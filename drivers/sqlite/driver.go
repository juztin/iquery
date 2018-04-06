package sqlite

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/juztin/iquery"
	"github.com/juztin/iquery/drivers"

	_ "github.com/mattn/go-sqlite3"
)

const driverName = "sqlite3"

type byName []iquery.Column

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name < s[j].Name }

type Sqlite struct {
	drivers.Driver
}

func (d Sqlite) Name() string {
	return driverName
}

func (d Sqlite) DriverName() string {
	return driverName
}

func (d Sqlite) ConnectionString(s iquery.Server) string {
	// Just use the name of the database as the filename
	return s.Database
}

func (d Sqlite) SchemasQuery() string {
	return `select "main" as name`
}

func (d Sqlite) TablesQuery(database string) string {
	return "select NAME, null as DESCRIPTION from SQLITE_MASTER order by NAME"
}

func (d Sqlite) Columns(s iquery.Server, schema, table string) ([]iquery.Column, error) {
	db, err := sql.Open(d.Name(), d.ConnectionString(s))
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []iquery.Column
	for rows.Next() {
		c := new(iquery.Column)
		var i sql.NullInt64
		err = rows.Scan(&i, &c.Name, &c.Type, &c.IsNullable, &c.Default, &i)
		if err != nil {
			return columns, err
		}
		columns = append(columns, *c)
	}
	sort.Sort(byName(columns))
	return columns, err
}

func New() iquery.Driver {
	d := new(Sqlite)
	d.Driver = drivers.New(d, nil)
	return d
}
