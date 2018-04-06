package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/juztin/iquery"
	"github.com/juztin/iquery/drivers/db2"
	"github.com/juztin/iquery/drivers/mssql"
	"github.com/juztin/iquery/drivers/mysql"
	"github.com/juztin/iquery/drivers/postgres"
	"github.com/juztin/iquery/drivers/sqlite"
	"github.com/juztin/iquery/handlers"
)

type TemplateData struct {
	Servers     iquery.Servers
	Theme       string
	Placeholder string
}

var (
	tmpl        *template.Template
	theme       = "vs"
	placeholder = "SELECT * FROM users;"
	debug       = false
)

func decodeJSON(r *http.Request) (iquery.Query, error) {
	var q iquery.Query
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&q)
	return q, err
}

func index(w http.ResponseWriter, r *http.Request) {
	if debug {
		log.Println("reloading template")
		tmpl = template.Must(template.ParseFiles("./templates/index.tmpl"))
	}

	d := TemplateData{handlers.Servers, theme, placeholder}
	err := tmpl.Execute(w, d)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	iquery.AddDriver(db2.NewAS400(), db2.NewCE(), mssql.New(), mysql.New(), postgres.New(), sqlite.New())
	handlers.Servers = iquery.ServersFromEnvironment()

	// Uncomment when testing to auto-reload template changes
	tmpl = template.Must(template.ParseFiles("./templates/index.tmpl"))

	s := os.Getenv("THEME")
	if s != "" {
		theme = s
	}
	s = os.Getenv("PLACEHOLDER")
	if s != "" {
		placeholder = s
	}

	if os.Getenv("DEBUG") == "true" {
		debug = true
	}

	// Static Content
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Routes
	r.HandleFunc("/databases/", handlers.Databases).Methods("GET")
	r.HandleFunc("/databases/{database}", handlers.Query).Methods("POST").Queries("action", "query")
	r.HandleFunc("/databases/{database}/schemas/", handlers.Schemas).Methods("GET")
	r.HandleFunc("/databases/{database}/schemas/{schema}/tables/", handlers.Tables).Methods("GET")
	r.HandleFunc("/databases/{database}/schemas/{schema}/tables/{table}", handlers.Columns).Methods("GET")
	r.HandleFunc("/", index)

	log.Println("Running...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
