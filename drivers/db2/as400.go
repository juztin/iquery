package db2

import (
	"fmt"

	"github.com/juztin/iquery"
	"github.com/juztin/iquery/drivers"
)

type DB2_AS400 struct {
	drivers.Driver
}

func (d DB2_AS400) Name() string {
	return "db2-i"
}

func (d DB2_AS400) DriverName() string {
	return driverName
}

func (d DB2_AS400) ConnectionString(s iquery.Server) string {
	return fmt.Sprintf(connFmt, s.Database, s.Hostname, s.Port, s.Username, s.Password)
}

func (d DB2_AS400) SchemasQuery() string {
	return "select distinct(trim(TABLE_SCHEMA)) as SCHEMA from QSYS2.SYSTABLES order by SCHEMA;"
}

func (d DB2_AS400) TablesQuery(schema string) string {
	return fmt.Sprintf("select trim(TABLE_NAME) as NAME, TABLE_TEXT"+
		" from QSYS2.SYSTABLES"+
		" where TABLE_SCHEMA='%s' and TABLE_TYPE='T'"+
		" order by NAME;", schema)
}

func (d DB2_AS400) ColumnsQuery(schema, table string) string {
	return fmt.Sprintf("select trim(SYSTEM_COLUMN_NAME), trim(COLUMN_HEADING), lower(DATA_TYPE), COLUMN_DEFAULT, LENGTH, case when IS_NULLABLE='yes' then 1 else 0 end"+
		" from QSYS2.SYSCOLUMNS"+
		" where TABLE_SCHEMA='%s' and TABLE_NAME='%s'"+
		" order by COLUMN_NAME;", schema, table)
}
