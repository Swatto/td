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
