package main

import (
	ipspeicher "github.com/HoffmannP/ipspeicher"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strings"
)

const dataSourceName = "ipSaver.sqlite"
const driverName = "sqlite3"
const addr = ":6947"

var speicher ipspeicher.Speicher

func handler(w http.ResponseWriter, r *http.Request) {
	ra := r.RemoteAddr
	e := ipspeicher.NewEintrag(
		r.URL.Path[1:],
		ra[:strings.Index(ra, ":")])
	switch r.URL.RawQuery {
	case "show":
		var err error
		var einträge []ipspeicher.Eintrag
		if e.Name == "" {
			einträge, err = speicher.Namen()
		} else {
			einträge, err = speicher.Verlauf(e)
		}
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
		}
		for _, e = range einträge {
			// w.WriteHeader("Content-Type: text/html")
			w.Write([]byte(e.String()))
			w.Write([]byte("\n"))
		}
	default:
		alreadyIn, err := speicher.Sichern(e)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
		}
		if alreadyIn {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(202)
		}

	}
}

func main() {
	speicher, err := ipspeicher.NewSpeicher(driverName, dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer speicher.Schließen()

	log.Println("Gestartet")
	http.HandleFunc("/", handler)
	http.ListenAndServe(addr, nil)
}
