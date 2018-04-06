package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func Columns(w http.ResponseWriter, r *http.Request) {
	s, err := server(r)
	if handleError(w, err, "", http.StatusNotFound) {
		return
	}

	vars := mux.Vars(r)
	tables, err := s.Driver.Columns(s, vars["schema"], vars["table"])
	if handleError(w, err, "", http.StatusInternalServerError) {
		return
	}
	j, _ := json.Marshal(tables)
	w.Write(j)
}
