package main

import (
	"errors"
	"os"
	"path"
)

var errIsNotAFile = errors.New("the database path is not a file")
var errLocalDbFileNotFound = errors.New("the local .todos file was not found")

var cachedDBPath = ""

func getDBPath() string {
	if cachedDBPath != "" {
		return cachedDBPath
	}

	dbPath, err := calculateDBPath()
	if err != nil {
		return ""
	}

	cachedDBPath = dbPath
	return dbPath
}

func calculateDBPath() (string, error) {
	dbPath, err := tryCwdAndParentFolders()
	if err != nil {
		dbPath, err = tryEnv()
	}
	return dbPath, err
}

func tryDir(dir string) (string, error) {
	dbPath := path.Join(dir, ".todos")
	fi, err := os.Stat(dbPath)
	if err != nil {
		return "", err
	}

	if fi.IsDir() {
		return "", errIsNotAFile
	}

	return dbPath, nil
}

func tryCwdAndParentFolders() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		filePath, err := tryDir(cwd)
		if err == nil {
			return filePath, err
		}

		if len(cwd) == 1 {
			break
		}

		cwd = path.Dir(cwd)
	}

	return "", errLocalDbFileNotFound
}

func tryEnv() (string, error) {
	return os.Getenv("TODO_DB_PATH"), nil
}
