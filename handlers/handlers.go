package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/juztin/iquery"
)

var Servers map[string]iquery.Server

type Error struct {
	Error string `json:"error"`
}

func server(r *http.Request) (iquery.Server, error) {
	vars := mux.Vars(r)
	server, ok := Servers[vars["database"]]
	if !ok {
		return server, errors.New(fmt.Sprintf("Database name doesn't exist: '%s'", vars["database"]))
	}
	return server, nil
}

func handleError(w http.ResponseWriter, err error, message string, statusCode int) bool {
	if err == nil {
		return false
	}

	if message == "" {
		message = err.Error()
	}
	j, _ := json.Marshal(Error{message})
	w.WriteHeader(statusCode)
	w.Write(j)
	return true
}
