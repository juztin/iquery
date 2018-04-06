package postgres

import (
	"fmt"

	"github.com/juztin/iquery"
	"github.com/juztin/iquery/drivers"

	_ "github.com/lib/pq"
)

const (
	driverName = "postgres"
	// "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
	connFmt = "postgres://%s:%s@%s:%d/%s?sslmode=verify-full"
)

type Postgres struct {
	drivers.Driver
}

func (d Postgres) Name() string {
	return driverName
}

func (d Postgres) DriverName() string {
	return driverName
}
func (d Postgres) ConnectionString(s iquery.Server) string {
	return fmt.Sprintf(connFmt, s.Username, s.Password, s.Hostname, s.Port, s.Database)
}

func (d Postgres) SchemasQuery() string {
	panic("not implemented")
}

func (d Postgres) TablesQuery(database string) string {
	panic("not implemented")
}

func (d Postgres) ColumnsQuery(database, table string) string {
	panic("not implemented")
}

func New() iquery.Driver {
	d := new(Postgres)
	d.Driver = drivers.New(d, nil)
	return d
}
