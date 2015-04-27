package main

import (
	"errors"
	"os"
	"path"
)

var IsNotAFileErr = errors.New("The database path is not a file")
var LocalDbFileNotFoundErr = errors.New("The local .todos file was not found")

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
		return "", IsNotAFileErr
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

	return "", LocalDbFileNotFoundErr
}

func tryEnv() (string, error) {
	return os.Getenv("TODO_DB_PATH"), nil
}
