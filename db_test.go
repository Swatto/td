package main

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestGetDbPath(t *testing.T) {
	err := _CreateFakeDb()
	if err != nil {
		t.Error(err)
	}
	path := GetDBPath()
	fmt.Println(path)
	err = _DeleteFakeDb()
	if err != nil {
		t.Error(err)
	}
}

func _CreateFakeDb() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	cwd = cwd + "/.todos"

	fi, err := os.Stat(cwd)
	if (err != nil && os.IsNotExist(err)) || fi.Size() == 0 {
		w, err := os.Create(cwd)
		_, err = w.WriteString("[]")
		defer w.Close()
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func _DeleteFakeDb() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	cwd = cwd + "/.todos"

	stats, err := os.Stat(cwd)
	if err == nil && !stats.IsDir() {
		err := os.Remove(cwd)
		return err
	}

	if err != nil {
		return err
	}

	return errors.New("bad return")
}
