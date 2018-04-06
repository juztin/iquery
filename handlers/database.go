package handlers

import (
	"encoding/json"
	"net/http"
)

var dbs map[string]string

func Databases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if dbs == nil {
		dbs = make(map[string]string)
		for name, server := range Servers {
			//dbs = append(dbs, DB{name, server.Name()})
			dbs[name] = server.Driver.Name()
		}
	}
	//j, err := json.Marshal(dbs)
	j, _ := json.Marshal(dbs)
	w.Write(j)
}

func Schemas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	s, err := server(r)
	if handleError(w, err, "", http.StatusNotFound) {
		return
	}

	schemas, err := s.Driver.Schemas(s)
	if handleError(w, err, "", http.StatusInternalServerError) {
		return
	}
	//j, err := json.Marshal(schemas)
	j, _ := json.Marshal(schemas)
	w.Write(j)
}
