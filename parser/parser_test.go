package parser

import (
	"testing"
)

func TestPeriod(t *testing.T) {
	if _, err := ParsePeriod("05"); err != nil {
		t.Fail()
		t.Log("number 1 should be parsed")
	}

	if _, err := ParsePeriod("-1"); err == nil {
		t.Fail()
		t.Log("number 2 shouldn't be parsed")
	}

	if _, err := ParsePeriod("2.5"); err == nil {
		t.Fail()
		t.Log("number 3 shouldn't be parsed")
	}
}
func TestParseDate(t *testing.T) {
	if _, err := ParseDate("01/01/1970"); err != nil {
		t.Fail()
		t.Log("date-only should be parsed")
	}
	if _, err := ParseDate("12:34"); err != nil {
		t.Fail()
		t.Log("time-only should be parsed")
	}
	if _, err := ParseDate("05/05/2009 12:34"); err != nil {
		t.Fail()
		t.Log("date-time should be parsed")
	}
}
