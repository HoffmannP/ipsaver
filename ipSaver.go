package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strings"
	"time"
)

const dbFileName = "ipSaver.sqlite"
const listenOn = ":6947"

type Eintrag struct {
	Name string
	Ip   string
	Date time.Time
}

type Row interface {
	Scan(dest ...interface{}) error
}

func SqlEintrag(row Row) (Eintrag, error) {
	var e Eintrag
	var stamp int64
	if err := row.Scan(&e.Name, &e.Ip, &stamp); err != nil {
		if err == sql.ErrNoRows {
			return Eintrag{}, err
		}
	}
	e.Date = time.Unix(stamp, 0)
	return e, nil
}

func (e Eintrag) String() string {
	return fmt.Sprintf("%s\t%s\t%s\n", e.Ip, e.Date, e.Name)
}

type Speicher struct {
	Db *sql.DB
}

func (ips *Speicher) Gespeichert(e Eintrag) bool {
	neuste, err := ips.Neuste(e.Name)
	if err != nil {
		return false
	}
	return neuste.Ip == e.Ip
}

func (ips *Speicher) Speichern(e Eintrag) bool {
	_, err := ips.Db.Exec("INSERT INTO ips VALUES (?, ?, strftime('%s', 'now'))", e.Name, e.Ip)
	if err != nil {
		log.Print("Speicher.Speichern")
		log.Fatal(err)
	}
	return true
}

func (ips *Speicher) Neuste(name string) (Eintrag, error) {
	row := ips.Db.QueryRow("SELECT * FROM ips WHERE name = ? ORDER BY seit DESC LIMIT 1", name)
	return SqlEintrag(row)
}

func (ips *Speicher) count(query string, args ...interface{}) (int, error) {
	row := ips.Db.QueryRow(query, args...)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (ips *Speicher) Verlauf(e Eintrag) []Eintrag {
	var count int
	count, err := ips.count("SELECT COUNT(*) FROM ips WHERE name = ?", e.Name)
	if err != nil {
		log.Print("Speicher.Verlauf(1/3)")
		log.Fatal(err)
	}

	Einträge := make([]Eintrag, count)
	rows, err := ips.Db.Query("SELECT * FROM ips WHERE name = ?", e.Name)
	if err != nil {
		log.Print("Speicher.Verlauf(2/3)")
		log.Fatal(err)
	}

	i := 0
	for rows.Next() {
		Einträge[i], _ = SqlEintrag(rows)
		i++
	}
	if err := rows.Err(); err != nil {
		log.Print("Speicher.Verlauf(3/3)")
		log.Fatal(err)
	}
	return Einträge
}

func (ips *Speicher) Namen() []Eintrag {
	count, err := ips.count("SELECT DISTINCT COUNT(name) FROM ips")
	if err != nil {
		log.Print("Speicher.Namen(1/3)")
		log.Fatal(err)
	}

	rows, err := ips.Db.Query("SELECT DISTINCT name FROM ips")
	if err != nil {
		log.Print("Speicher.Namen(2/3)")
		log.Fatal(err)
	}
	namen := make([]Eintrag, count)
	i := 0
	for rows.Next() {
		namen[i], _ = SqlEintrag(rows)
		i++
	}
	if err := rows.Err(); err != nil {
		log.Print("Speicher.Namen(3/3)")
		log.Fatal(err)
	}
	return namen
}

func (ips *Speicher) Sichern(e Eintrag) int {
	if ips.Gespeichert(e) {
		return 200
	}
	if ips.Speichern(e) {
		return 202
	}
	return 500
}

func (ips *Speicher) Aufruf(w http.ResponseWriter, r *http.Request) {
	ra := r.RemoteAddr
	e := Eintrag{
		Name: r.URL.Path[1:],
		Ip:   ra[:strings.Index(ra, ":")],
	}
	switch r.URL.RawQuery {
	case "show":
		for _, e = range ips.Verlauf(e) {
			w.Write([]byte(e.String()))
		}
	default:
		w.WriteHeader(ips.Sichern(e))
	}
}

func (ips *Speicher) Show(w http.ResponseWriter, r *http.Request) {
	for _, e := range ips.Namen() {
		w.Write([]byte(e.String()))
	}
}

func neuerSpeicher(dbname string) (*Speicher, error) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return nil, err
	}
	db.Exec("CREATE TABLE IF NOT EXISTS ips (name TEXT, ip TEXT, seit INT)")
	return &Speicher{Db: db}, nil
}

func main() {
	is, err := neuerSpeicher(dbFileName)
	if err != nil {
		log.Print("main")
		log.Fatal(err)
	}
	defer is.Db.Close()
	log.Println("Gestartet")
	http.HandleFunc("/", is.Aufruf)
	http.ListenAndServe(listenOn, nil)
}
