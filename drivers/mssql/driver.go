package mssql

import (
	"encoding/hex"
	"fmt"
	"unicode"

	"github.com/juztin/iquery"
	"github.com/juztin/iquery/drivers"
	"github.com/juztin/iquery/vendor/github.com/satori/go.uuid"

	_ "github.com/denisenkom/go-mssqldb"
)

const (
	driverName = "mssql"
	connFmt    = "server=%s;port=%d;user id=%s;password=%s;database=%s"
)

type MSSQL struct {
	drivers.Driver
}

func (d MSSQL) Name() string {
	return driverName
}

func (d MSSQL) DriverName() string {
	return driverName
}

func (d MSSQL) ConnectionString(s iquery.Server) string {
	return fmt.Sprintf(connFmt, s.Hostname, s.Port, s.Username, s.Password, s.Database)
}

func (d MSSQL) SchemasQuery() string {
	return "select NAME from SYS.DATABASES order by NAME"
}

func (d MSSQL) TablesQuery(database string) string {
	return fmt.Sprintf("select TABLE_NAME, null as DESCRIPTION from %s.INFORMATION_SCHEMA.TABLES where TABLE_TYPE='base table' order by TABLE_NAME", database)
}

func (d MSSQL) ColumnsQuery(database, table string) string {
	return fmt.Sprintf("select COLUMN_NAME, null as DESCRIPTION, DATA_TYPE, COLUMN_DEFAULT, CHARACTER_MAXIMUM_LENGTH, cast(case when IS_NULLABLE='yes' then 1 else 0 end as bit)"+
		" from %s.INFORMATION_SCHEMA.COLUMNS"+
		" where table_name='%s'"+
		" order by COLUMN_NAME", database, table)
}

/*
SELECT * FROM sys.types order
=============================

name                 system_type_id user_type_id schema_id principal_id max_length precision scale collation_name    is_nullable is_user_defined is_assembly_type default_object_id rule_object_id is_table_type
-------------------- -------------- ------------ --------- ------------ ---------- --------- ----- ----------------- ----------- --------------- ---------------- ----------------- -------------- -------------
image                34             34           4         NULL         16         0         0     NULL              1           0               0                0                 0              0
text                 35             35           4         NULL         16         0         0     Persian_100_CI_AI 1           0               0                0                 0              0
uniqueidentifier     36             36           4         NULL         16         0         0     NULL              1           0               0                0                 0              0
date                 40             40           4         NULL         3          10        0     NULL              1           0               0                0                 0              0
time                 41             41           4         NULL         5          16        7     NULL              1           0               0                0                 0              0
datetime2            42             42           4         NULL         8          27        7     NULL              1           0               0                0                 0              0
datetimeoffset       43             43           4         NULL         10         34        7     NULL              1           0               0                0                 0              0
tinyint              48             48           4         NULL         1          3         0     NULL              1           0               0                0                 0              0
smallint             52             52           4         NULL         2          5         0     NULL              1           0               0                0                 0              0
int                  56             56           4         NULL         4          10        0     NULL              1           0               0                0                 0              0
smalldatetime        58             58           4         NULL         4          16        0     NULL              1           0               0                0                 0              0
real                 59             59           4         NULL         4          24        0     NULL              1           0               0                0                 0              0
money                60             60           4         NULL         8          19        4     NULL              1           0               0                0                 0              0
datetime             61             61           4         NULL         8          23        3     NULL              1           0               0                0                 0              0
float                62             62           4         NULL         8          53        0     NULL              1           0               0                0                 0              0
sql_variant          98             98           4         NULL         8016       0         0     NULL              1           0               0                0                 0              0
ntext                99             99           4         NULL         16         0         0     Persian_100_CI_AI 1           0               0                0                 0              0
bit                  104            104          4         NULL         1          1         0     NULL              1           0               0                0                 0              0
decimal              106            106          4         NULL         17         38        38    NULL              1           0               0                0                 0              0
numeric              108            108          4         NULL         17         38        38    NULL              1           0               0                0                 0              0
smallmoney           122            122          4         NULL         4          10        4     NULL              1           0               0                0                 0              0
bigint               127            127          4         NULL         8          19        0     NULL              1           0               0                0                 0              0
hierarchyid          240            128          4         NULL         892        0         0     NULL              1           0               1                0                 0              0
geometry             240            129          4         NULL         -1         0         0     NULL              1           0               1                0                 0              0
geography            240            130          4         NULL         -1         0         0     NULL              1           0               1                0                 0              0
varbinary            165            165          4         NULL         8000       0         0     NULL              1           0               0                0                 0              0
varchar              167            167          4         NULL         8000       0         0     Persian_100_CI_AI 1           0               0                0                 0              0
binary               173            173          4         NULL         8000       0         0     NULL              1           0               0                0                 0              0
char                 175            175          4         NULL         8000       0         0     Persian_100_CI_AI 1           0               0                0                 0              0
timestamp            189            189          4         NULL         8          0         0     NULL              0           0               0                0                 0              0
nvarchar             231            231          4         NULL         8000       0         0     Persian_100_CI_AI 1           0               0                0                 0              0
nchar                239            239          4         NULL         8000       0         0     Persian_100_CI_AI 1           0               0                0                 0              0
xml                  241            241          4         NULL         -1         0         0     NULL              1           0               0                0                 0              0
sysname              231            256          4         NULL         256        0         0     Persian_100_CI_AI 0           0               0                0                 0              0
CalculatedCreditInfo 243            257          9         NULL         -1         0         0     NULL              0           1               0                0                 0              1
udt_QoutaDetail      243            258          21        NULL         -1         0         0     NULL              0           1               0                0                 0              1
BeforeUpdate         243            259          22        NULL         -1         0         0     NULL              0           1               0                0                 0              1
udt_StoreInventory   243            260          26        NULL         -1         0         0     NULL              0           1               0                0                 0              1
udt_WKFHistory       243            261          32        NULL         -1         0         0     NULL              0           1               0                0                 0              1
IDTable              243            262          1         NULL         -1         0         0     NULL
*/
func mapBytes(b []byte) interface{} {
	// TODO: Not sure if this is the best way to solve this.
	//       _(Have tried getting at `rows.ColumnTypes()`, but don't get any useful information.)_

	switch len(b) {
	/*case 6, 7, 8:
	if b[len(b)-5] == 46 { // [money]? _(checking that the 4ᵗʰ character from the right is: '.')_
		return string(b)
	}*/
	case 16: // [uniqueidentifier]
		u, err := uuid.FromBytes(b)
		if err == nil {
			s := u.String()
			return &s
		}
	}

	if isASCII(b) {
		s := string(b)
		return &s
	}

	// As a last resort, encode the bytes to hex.
	return hex.EncodeToString(b)
}

func isASCII(b []byte) bool {
	for i := range b {
		if b[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func New() iquery.Driver {
	mapper := drivers.ValueMapper(mapBytes)
	d := new(MSSQL)
	d.Driver = drivers.New(d, mapper)
	return d
}
