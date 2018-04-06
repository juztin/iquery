package mysql

import (
	"fmt"

	"github.com/juztin/iquery"
	"github.com/juztin/iquery/drivers"

	_ "github.com/go-sql-driver/mysql"
)

const (
	driverName = "mysql"
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	connFmt = "%s:%s@tcp(%s:%d)/%s"
)

type Mysql struct {
	drivers.Driver
}

func (d Mysql) Name() string {
	return driverName
}

func (d Mysql) DriverName() string {
	return driverName
}

func (d Mysql) ConnectionString(s iquery.Server) string {
	return fmt.Sprintf(connFmt, s.Username, s.Password, s.Hostname, s.Port, s.Database)
}

func (d Mysql) SchemasQuery() string {
	return "select SCHEMA_NAME from INFORMATION_SCHEMA.SCHEMATA order by SCHEMA_NAME"
}

func (d Mysql) TablesQuery(schema string) string {
	return fmt.Sprintf("select TABLE_NAME, null as DESCRIPTION from INFORMATION_SCHEMA.TABLES where TABLE_SCHEMA='%s' order by TABLE_NAME", schema)
}

func (d Mysql) ColumnsQuery(schema, table string) string {
	return fmt.Sprintf("select COLUMN_NAME, null as DESCRIPTION, COLUMN_TYPE, COLUMN_DEFAULT, null, case when IS_NULLABLE='YES' then true else false end"+
		" from INFORMATION_SCHEMA.COLUMNS"+
		" where TABLE_SCHEMA = '%s'"+
		" and TABLE_NAME = '%s'"+
		" order by COLUMN_NAME", schema, table)
}

func New() iquery.Driver {
	d := new(Mysql)
	d.Driver = drivers.New(d, nil)
	return d
}
