package db

import (
	"encoding/json"
	"errors"
	"os"
)

func CreateStoreFileIfNeeded(path string) error {
	fi, err := os.Stat(path)
	if (err != nil && os.IsNotExist(err)) || fi.Size() == 0 {
		w, _ := os.Create(path)
		_, err = w.WriteString("[]")
		defer w.Close()
		return err
	}

	if err != nil {
		return err
	}

	if fi.Size() != 0 {
		return errors.New("StoreAlreadyExist")
	}

	return nil
}

func Read() (*Collection, error) {
	c := Collection{}
	file, err := os.OpenFile(GetDBPath(), os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&c.Todos)
	return &c, err
}

func Save(c *Collection) error {
	file, err := os.OpenFile(GetDBPath(), os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	data, err := json.MarshalIndent(&c.Todos, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}
