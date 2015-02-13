package main

import (
	"errors"
	"os"
	"path"
)

var IsNotAFileErr = errors.New("The database path is not a file")

var cachedDBPath = ""

func GetDBPath() string {
	if cachedDBPath != "" {
		return cachedDBPath
	}

	dbPath, err := getDBPath()
	if err != nil {
		return ""
	}

	cachedDBPath = dbPath
	return dbPath
}

func getDBPath() (string, error) {
	dbPath, err := tryCwd()
	if err != nil {
		dbPath, err = tryEnv()
	}
	return dbPath, err
}

func tryCwd() (string, error) {
	cw, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dbPath := path.Join(cw, ".todos")
	fi, err := os.Stat(dbPath)
	if err != nil {
		return "", err
	}

	if fi.IsDir() {
		return "", IsNotAFileErr
	}

	return dbPath, nil
}

func tryEnv() (string, error) {
	return os.Getenv("TODO_DB_PATH"), nil
}
