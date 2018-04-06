package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/juztin/iquery"
)

const defaultLimit = 1000

type RowReader func(http.ResponseWriter, *http.Request, iquery.Driver, *sql.Rows) (interface{}, bool)

func limit(r *http.Request) int {
	s := r.URL.Query().Get("limit")
	if s == "" {
		return defaultLimit
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return defaultLimit
	}
	return n
}

func requestedSQL(r *http.Request) (string, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func query(w http.ResponseWriter, r *http.Request, q string, fn RowReader) (interface{}, bool) {
	// Get Matching Server
	s, err := server(r)
	if handleError(w, err, "", http.StatusNotFound) {
		return nil, false
	}

	driver := s.Driver
	// Open DB Connection
	db, err := sql.Open(driver.DriverName(), driver.ConnectionString(s))
	if handleError(w, err, fmt.Sprintf("Failed to connect to database: %s", err), http.StatusInternalServerError) {
		return nil, false
	}

	// Execute Query
	rows, err := db.Query(q)
	if handleError(w, err, fmt.Sprintf("Query failure: %s", err), http.StatusBadRequest) {
		return nil, false
	}
	defer rows.Close()
	return fn(w, r, driver, rows)
}

func rowReader(w http.ResponseWriter, r *http.Request, d iquery.Driver, rows *sql.Rows) (interface{}, bool) {
	var results []iquery.Result
	for {
		result := iquery.Result{}
		// Get Columns
		var err error
		result.Columns, err = rows.Columns()
		if handleError(w, err, fmt.Sprintf("Failed to retrieve columns: %s", err), http.StatusInternalServerError) {
			return result, false
		}

		// Map Rows
		result.Rows, err = d.MapRows(rows, len(result.Columns), limit(r))
		if handleError(w, err, fmt.Sprintf("Failed read row: %s", err), http.StatusInternalServerError) {
			return result, false
		}

		results = append(results, result)
		if !rows.NextResultSet() {
			break
		}

	}
	if len(results) == 1 {
		return results[0], true
	}
	return results, true
}

func Query(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get Requested SQL
	q, err := requestedSQL(r)
	if handleError(w, err, "Invalid payload", http.StatusBadRequest) {
		return
	}

	// Get DB Results
	result, ok := query(w, r, q, rowReader)
	if !ok {
		return
	}
	j, err := json.Marshal(result)
	w.Write(j)
}
