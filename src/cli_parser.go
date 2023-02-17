package main

import (
	"errors"
	"time"
)

func parseDate(d string) (time.Time, error) {
	dt, err := time.Parse("01/02/2006 15:04", d)
	if err != nil {
		dt, err = time.Parse("01/02/2006", d)
		if err != nil {
			dt, err = time.Parse("15:04", d)
			if err != nil {
				return time.Time{}, errors.New("invalid date format")
			}
			dt.AddDate(time.Now().Day(), int(time.Now().Month()), time.Now().Year())
		}
	}
	return dt, nil
}
