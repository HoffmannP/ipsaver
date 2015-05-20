package ipspeicher

import (
	"testing"
	"time"
)

// fake Scan for Eintrag
func (fakeRow Eintrag) Scan(dest ...interface{}) error {
	*dest[0].(*string) = fakeRow.Name
	*dest[1].(*string) = fakeRow.Ip
	*dest[2].(*int64) = fakeRow.Date.Unix()
	return nil
}

func TestNewEintrag(t *testing.T) {
	exp := Eintrag{
		Name: "Name",
		Ip:   "IP",
		Date: time.Time{},
	}
	is := NewEintrag(exp.Name, exp.Ip)

	if is.Name != exp.Name {
		t.Errorf("is.Name == %s, want %s", is.Name, exp.Name)
	}
	if is.Ip != exp.Ip {
		t.Errorf("is.Ip == %s, want %s", is.Ip, exp.Ip)
	}
	if is.Date != exp.Date {
		t.Errorf("is.Date == %s, want %s", is.Date, exp.Date)
	}
}

func TestSqlEintrag(t *testing.T) {
	exp := Eintrag{
		Name: "Name",
		Ip:   "IP",
		Date: time.Date(2015, 5, 20, 14, 56, 0, 0, time.FixedZone("CEST", +2)),
	}

	is, err := SqlEintrag(exp)
	if err != nil {
		t.Errorf("err == %s, want %s", err, nil)
	}

	if is.Name != exp.Name {
		t.Errorf("is.Name == %s, want %s", is.Name, exp.Name)
	}
	if is.Ip != exp.Ip {
		t.Errorf("is.Ip == %s, want %s", is.Ip, exp.Ip)
	}
	if is.Date.Unix() != exp.Date.Unix() {
		t.Errorf("is.Date == %s, want %s", is.Date, exp.Date)
	}
}

func TestString(t *testing.T) {
	exp := Eintrag{
		Name: "Name",
		Ip:   "IP",
		Date: time.Date(2015, 5, 20, 14, 56, 0, 0, time.FixedZone("CEST", +2)),
	}
	expString := "IP\t2015-05-20 14:56:00\tName"
	isString := exp.String()
	if isString != expString {
		t.Errorf("is.Date == %s, want %s", isString, expString)
	}
}
