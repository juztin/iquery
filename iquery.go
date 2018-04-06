package iquery

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Servers map[string]Server

type Server struct {
	Driver   Driver
	Hostname string
	Port     int
	Database string
	Username string
	Password string
}

type Query struct {
	Name  string `json:"database"`
	Query string `json:"query"`
	Limit int    `json:"limit"`
}

type Table struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type Column struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Type        string  `json:"type"`
	Default     *string `json:"default"`
	IsNullable  *bool   `json:"isNullable"`
}

type Rows [][]interface{}

type Result struct {
	Columns []string `json:"columns"`
	Rows    Rows     `json:"rows"`
}

type Driver interface {
	Name() string
	DriverName() string
	ConnectionString(s Server) string
	Schemas(s Server) ([]string, error)
	Tables(s Server, database string) ([]Table, error)
	Columns(s Server, database, table string) ([]Column, error)
	MapRows(r *sql.Rows, columns, limit int) (Rows, error)
}

var drivers map[string]Driver

func AddDriver(d ...Driver) {
	if drivers == nil {
		drivers = make(map[string]Driver)
	}
	for i := range d {
		drivers[d[i].Name()] = d[i]
	}
}

func ServersFromEnvironment() map[string]Server {
	m := make(map[string]Server)
	env := os.Environ()
	for i, s := range env {
		if !strings.HasPrefix(s, "DB_") {
			continue
		}

		s := strings.Split(env[i], "=")
		// Split should result in: `[name, driver, hostname, port, database, username, password]`
		values := strings.Split(s[1], "|")
		if len(values) != 7 {
			log.Fatalln(fmt.Sprintf("Invalid server entry: '%s'", env[i]))
		}
		// Trim values.
		for i := range values[:len(values)-1] {
			values[i] = strings.TrimSpace(values[i])
		}
		// Parse the port.
		port, err := strconv.Atoi(values[3])
		if values[3] != "" && err != nil {
			log.Fatalln(fmt.Sprintf("Invalid port value: '%s' in entry: '%s'", port, env[i]))
		}
		// Add the driver.
		driver, ok := drivers[values[1]]
		if !ok {
			var s []string
			for i := range drivers {
				s = append(s, drivers[i].Name())
			}
			log.Fatalln(fmt.Sprintf("Invalid driver: '%s', must be one of [%s]", env[i], strings.Join(s, ", ")))
		}
		m[values[0]] = Server{driver, values[2], port, values[4], values[5], values[6]}
	}
	return m
}
