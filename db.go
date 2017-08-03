package main

import (
	"fmt"
	"os"
	"path"
)

const ENVDBPATH = "TODO_DB_PATH"

type DataStore struct {
	Path string
}

func NewDataStore() (*DataStore, error) {
	ds := new(DataStore)
	ds.Path = os.Getenv(ENVDBPATH)

	if ds.Path == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return ds, err
		}
		ds.Path = path.Join(cwd, ".todos")

	} else {
		dir, file := path.Split(ds.Path)
		if file == "" {
			ds.Path = path.Join(dir, ".todos")
		}

		fileInfo, err := os.Stat(ds.Path)

		if os.IsExist(err) {
			if fileInfo.IsDir() {
				ds.Path = path.Join(ds.Path, ".todos")
				os.Setenv(ENVDBPATH, ds.Path)
			}
		}
	}

	return ds, nil
}

func (d *DataStore) Check() error {
	_, err := os.Stat(d.Path)
	if os.IsNotExist(err) {
		return fmt.Errorf("The database file \"%s\" doesn't exists", d.Path)
	}
	return nil
}

func (d *DataStore) Initialize() error {
	var err error

	dir, _ := path.Split(d.Path)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s: One or more directories in this path doesn't exist.", dir)
	}

	_, err = os.Stat(d.Path)
	if os.IsNotExist(err) {
		w, err := os.Create(d.Path)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = w.WriteString("[]")
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("%s: To-do file has been initialized before. ", d.Path)
	}

	return nil
}

//
// var IsNotAFileErr = errors.New("The database path is not a file")
// var LocalDbFileNotFoundErr = errors.New("The local .todos file was not found")
//
// var cachedDBPath = ""
//
// func GetDBPath() string {
// 	if cachedDBPath != "" {
// 		return cachedDBPath
// 	}
//
// 	dbPath, err := getDBPath()
// 	if err != nil {
// 		return ""
// 	}
//
// 	cachedDBPath = dbPath
// 	return dbPath
// }
//
// func getDBPath() (string, error) {
// 	dbPath, err := tryCwdAndParentFolders()
// 	if err != nil {
// 		dbPath, err = tryEnv()
// 	}
// 	return dbPath, err
// }
//
// func tryDir(dir string) (string, error) {
// 	dbPath := path.Join(dir, ".todos")
// 	fi, err := os.Stat(dbPath)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	if fi.IsDir() {
// 		return "", IsNotAFileErr
// 	}
//
// 	return dbPath, nil
// }
//
// func tryCwdAndParentFolders() (string, error) {
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		return "", err
// 	}
//
// 	for {
// 		filePath, err := tryDir(cwd)
// 		if err == nil {
// 			return filePath, err
// 		}
//
// 		if len(cwd) == 1 {
// 			break
// 		}
//
// 		cwd = path.Dir(cwd)
// 	}
//
// 	return "", LocalDbFileNotFoundErr
// }
//
// func tryEnv() (string, error) {
// 	return os.Getenv("TODO_DB_PATH"), nil
// }
