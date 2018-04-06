package db2

import (
	"fmt"

	"github.com/juztin/iquery"
	"github.com/juztin/iquery/drivers"
)

type DB2_CE struct {
	drivers.Driver
}

func (d DB2_CE) Name() string {
	return "db2-ce"
}

func (d DB2_CE) DriverName() string {
	return driverName
}

func (d DB2_CE) ConnectionString(s iquery.Server) string {
	return fmt.Sprintf(connFmt, s.Database, s.Hostname, s.Port, s.Username, s.Password)
}

func (d DB2_CE) SchemasQuery() string {
	return "select SCHEMANAME from SYSCAT.SCHEMATA;"
}

func (d DB2_CE) TablesQuery(schema string) string {
	return fmt.Sprintf("select trim(NAME) as NAME, '' as TABLE_TEXT"+
		" from SYSIBM.SYSTABLES"+
		" where TYPE='T'"+
		" and CREATOR='%s' order by NAME;", schema)
}

func (d DB2_CE) ColumnsQuery(schema, table string) string {
	return fmt.Sprintf("select"+
		"    c.COLUMN_NAME, '' as COLUMN_HEADING,"+
		"    lower(c.DATA_TYPE) as DATA_TYPE,"+
		"    c.COLUMN_DEFAULT,"+
		"    case when c.CHARACTER_MAXIMUM_LENGTH is not null"+
		"        then c.CHARACTER_MAXIMUM_LENGTH"+
		"        else c.NUMERIC_PRECISION"+
		"    end as LENGTH,"+
		"    case when c.IS_NULLABLE='YES'"+
		"        then 1"+
		"        else 0"+
		"    end as IS_NULLABLE"+
		" from SYSIBM.TABLES t"+
		" join SYSIBM.COLUMNS c"+
		" on t.TABLE_SCHEMA = c.TABLE_SCHEMA"+
		" and t.TABLE_NAME = c.TABLE_NAME"+
		" where t.TABLE_SCHEMA = '%s'"+
		" and t.TABLE_NAME = '%s'"+
		" order by t.TABLE_NAME, c.ORDINAL_POSITION", schema, table)
}
