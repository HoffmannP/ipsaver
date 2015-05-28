package ipspeicher

import (
	"database/sql"
)

type Speicher struct {
	Db *sql.DB
}

func (ips *Speicher) getNeuste(name string) (Eintrag, error) {
	row := ips.Db.QueryRow("SELECT * FROM ips WHERE name = ? ORDER BY seit DESC LIMIT 1", name)
	return SqlEintrag(row)
}

func (ips *Speicher) istGespeichert(e Eintrag) (bool, error) {
	neuste, err := ips.getNeuste(e.Name)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return neuste.Ip == e.Ip, nil
}

func (ips *Speicher) speichern(e Eintrag) error {
	_, err := ips.Db.Exec("INSERT INTO ips VALUES (?, ?, strftime('%s', 'now'))", e.Name, e.Ip)
	return err
}

func (ips *Speicher) anzahl(query string, args ...interface{}) (int, error) {
	row := ips.Db.QueryRow(query, args...)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (ips *Speicher) liste(count int, query string, args ...interface{}) ([]Eintrag, error) {
	rows, err := ips.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	einträge := make([]Eintrag, count)
	for i := 0; rows.Next(); i++ {
		einträge[i], err = SqlEintrag(rows)
		if err != nil {
			return einträge, err
		}
	}

	if err := rows.Err(); err != nil {
		return einträge, err
	}
	return einträge, nil
}

func (ips *Speicher) Verlauf(e Eintrag) ([]Eintrag, error) {
	count, err := ips.anzahl("SELECT COUNT(*) FROM ips WHERE name = ?", e.Name)
	if err != nil {
		return nil, err
	}

	return ips.liste(
		count,
		"SELECT * FROM ips WHERE name = ? ORDER BY seit DESC",
		e.Name)
}

func (ips *Speicher) Namen() ([]Eintrag, error) {
	count, err := ips.anzahl("SELECT DISTINCT COUNT(name) FROM ips")
	if err != nil {
		return nil, err
	}

	return ips.liste(
		count,
		"SELECT ips1.* FROM ips AS ips1 LEFT JOIN ips as ips2 ON ips1.name = ips2.name AND ips1.seit < ips2.seit WHERE ips2.seit IS NULL")
}

func (ips *Speicher) Sichern(e Eintrag) (bool, error) {
	result, err := ips.istGespeichert(e)
	if err != nil {
		return false, err
	}
	if result {
		return true, nil
	}

	if err = ips.speichern(e); err != nil {
		return false, err
	}
	return false, nil
}

func (ips *Speicher) Schließen() {
	ips.Db.Close()
}

func NewSpeicher(db *sql.DB, err error) (*Speicher, error) {
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS ips (name TEXT, ip TEXT, seit INT)")
	if err != nil {
		return nil, err
	}
	return &Speicher{Db: db}, nil
}
