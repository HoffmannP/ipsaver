package ipspeicher

import (
	"database/sql"
	"fmt"
	"time"
)

type Eintrag struct {
	Name string
	Ip   string
	Date time.Time
}

type Row interface {
	Scan(dest ...interface{}) error
}

func NewEintrag(name, ip string) Eintrag {
	return Eintrag{
		Name: name,
		Ip:   ip,
	}
}

func SqlEintrag(row Row) (e Eintrag, err error) {
	var timestamp int64
	if err = row.Scan(&e.Name, &e.Ip, &timestamp); err != nil {
		if err == sql.ErrNoRows {
			return e, err
		}
	}
	e.Date = time.Unix(timestamp, 0)
	return e, nil
}

func (e Eintrag) String() string {
	return fmt.Sprintf(
		"%s\t%s\t%s",
		e.Ip,
		e.Date.Format("2006-01-02 15:04:05"),
		e.Name)
}
