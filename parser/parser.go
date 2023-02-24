package parser

import (
	"errors"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

func MapArgs(args cli.Args) (*map[string]string, error) {
	m := make(map[string]string, 3)
	// description case
	if args.Len() == 2 {
		m["desc"] = args.Get(1)
		return &m, nil
	}
	for i := 2; i < args.Len(); i += 2 {
		k := args.Get(i - 1)
		v := args.Get(i)
		switch k {
		case "--desc", "-d":
			m["desc"] = args.Get(i)
		case "--date", "-D":
			_, err := ParseDate(v)
			if err != nil {
				return nil, err
			}
			m["date"] = v
		case "--period", "-p":
			vi, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			} else if vi < 0 {
				return nil, errors.New("invalid period")
			}
			m["period"] = v
		}
	}
	return &m, nil
}

func ParsePeriod(v string) (int, error) {
	d, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	} else if d < 0 {
		return 0, errors.New("invalid period")
	}
	return d, nil
}

func ParseDate(v string) (time.Time, error) {
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
