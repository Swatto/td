package main

import (
	"errors"
	"strconv"
	"time"
	"umutsevdi/td/todo"

	"github.com/urfave/cli/v2"
)

type Error string

func (e Error) Error() string { return string(e) }

const argError = Error("Error in argument")

func FindAndReplace(args cli.Args, t *todo.Todo) error {
	// description case
	if args.Len() == 2 {
		t.Desc = args.Get(1)
		return nil
	}
	// key value case
	for i := 2; i < args.Len(); i += 2 {
		err := parseEntry(t, args.Get(i-1), args.Get(i))
		if err != nil {
			return err
		}
	}
	return nil
}

func parseEntry(t *todo.Todo, key, value string) error {
	switch key {
	case "--desc", "-d":
		t.Desc = value
	case "--date", "-D":
		d, err := parseDate(value)
		if err != nil {
			return err
		}
		t.Deadline = d
	case "--period", "-p":
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		} else if v < 0 {
			return errors.New("invalid period")
		}
		t.Period = v
	}
	return nil
}

func parsePeriod(v string) (int, error) {
	d, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	} else if d < 0 {
		return 0, errors.New("invalid period")
	}
	return d, nil
}

func parseDate(v string) (time.Time, error) {
	if len(v) == 0 {
		return time.Time{}, nil
	}
	d, err := time.Parse("02/01/2006 15:04", v)
	if err != nil {
		d, err = time.Parse("02/01/2006", v)
		if err != nil {
			d, err = time.Parse("15:04", v)
			if err != nil {
				return time.Time{}, errors.New("invalid date format")
			}
			d = d.AddDate(time.Now().Year(), int(time.Now().Month())-1, time.Now().Day()-1)
		}
	}
	return d, nil
}
