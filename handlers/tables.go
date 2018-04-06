package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func Tables(w http.ResponseWriter, r *http.Request) {
	s, err := server(r)
	if err != nil {
		handleError(w, err, "", http.StatusNotFound)
		return
	}

	vars := mux.Vars(r)
	dbs, err := s.Driver.Tables(s, vars["schema"])
	if handleError(w, err, "", http.StatusInternalServerError) {
		return
	}
	j, _ := json.Marshal(dbs)
	w.Write(j)
}
