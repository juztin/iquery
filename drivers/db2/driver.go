package db2

import (
	_ "bitbucket.org/phiggins/db2cli"

	"github.com/juztin/iquery"
	"github.com/juztin/iquery/drivers"
)

const (
	driverName = "db2-cli"
	connFmt    = "DATABASE=%s; HOSTNAME=%s; PORT=%d; PROTOCOL=TCPIP; UID=%s; PWD=%s;"
)

func NewAS400() iquery.Driver {
	d := new(DB2_AS400)
	d.Driver = drivers.New(d, nil)
	return d
}

func NewCE() iquery.Driver {
	d := new(DB2_CE)
	d.Driver = drivers.New(d, nil)
	return d
}
