package main

import (
	"database/sql"
	ipsaver "github.com/HoffmannP/ipsaver"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strings"
)

const dataSourceName = "ipSaver.sqlite"
const driverName = "sqlite3"
const addr = ":6947"

var speicher *ipsaver.Speicher
var call = 0

func handler(w http.ResponseWriter, r *http.Request) {
	ra := r.RemoteAddr
	e := ipsaver.NewEintrag(
		r.URL.Path[1:],
		ra[:strings.Index(ra, ":")])
	switch r.URL.RawQuery {
	case "show":
		var err error
		var einträge []ipsaver.Eintrag
		if e.Name == "" {
			einträge = []ipsaver.Eintrag{}
			err = nil
			einträge, err = speicher.Namen()
		} else {
			einträge = []ipsaver.Eintrag{}
			err = nil
			einträge, err = speicher.Verlauf(e)
		}
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
		}
		for _, e = range einträge {
			w.Write([]byte(e.String()))
			w.Write([]byte("\n"))
		}
	default:
		alreadyIn, err := speicher.Sichern(e)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		if alreadyIn {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(202)
		}

	}
}

func main() {
	var err error
	speicher, err = ipsaver.NewSpeicher(sql.Open(driverName, dataSourceName))
	if err != nil {
		log.Fatal(err)
	}
	defer speicher.Schließen()
	log.Println("Gestartet")
	http.HandleFunc("/", handler)
	http.ListenAndServe(addr, nil)
}
