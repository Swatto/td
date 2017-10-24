package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Collection struct {
	Todos []*Todo
}

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

func (c *Collection) RemoveAtIndex(item int) {
	s := *c
	s.Todos = append(s.Todos[:item], s.Todos[item+1:]...)
	*c = s
}

func (c *Collection) RetrieveTodos() error {
	file, err := os.OpenFile(GetDBPath(), os.O_RDONLY, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&c.Todos)
	return err
}

func (c *Collection) WriteTodos() error {
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

func (c *Collection) ListPendingTodos() {
	for i := len(c.Todos) - 1; i >= 0; i-- {
		if c.Todos[i].Status != "pending" {
			c.RemoveAtIndex(i)
		}
	}
}

func (c *Collection) ListDoneTodos() {
	for i := len(c.Todos) - 1; i >= 0; i-- {
		if c.Todos[i].Status != "done" {
			c.RemoveAtIndex(i)
		}
	}
}

func (c *Collection) CreateTodo(newTodo *Todo) (int64, error) {
	var highestId int64 = 0
	for _, todo := range c.Todos {
		if todo.Id > highestId {
			highestId = todo.Id
		}
	}

	newTodo.Id = (highestId + 1)
	newTodo.Modified = time.Now().Local().String()
	c.Todos = append(c.Todos, newTodo)

	err := c.WriteTodos()
	return newTodo.Id, err
}

func (c *Collection) Find(id int64) (foundedTodo *Todo, err error) {
	founded := false
	for _, todo := range c.Todos {
		if id == todo.Id {
			foundedTodo = todo
			founded = true
		}
	}
	if !founded {
		err = errors.New("The todo with the id " + strconv.FormatInt(id, 10) + " was not found.")
	}
	return
}

func (c *Collection) Toggle(id int64) (*Todo, error) {
	todo, err := c.Find(id)

	if err != nil {
		return todo, err
	}

	if todo.Status == "done" {
		todo.Status = "pending"
	} else {
		todo.Status = "done"
	}
	todo.Modified = time.Now().Local().String()

	err = c.WriteTodos()
	if err != nil {
		err = errors.New("Todos couldn't be saved")
		return todo, err
	}

	return todo, err
}

func (c *Collection) Modify(id int64, desc string) (*Todo, error) {
	todo, err := c.Find(id)

	if err != nil {
		return todo, err
	}

	todo.Desc = desc
	todo.Modified = time.Now().Local().String()

	err = c.WriteTodos()
	if err != nil {
		err = errors.New("Todos couldn't be saved")
		return todo, err
	}

	return todo, err
}

func (c *Collection) RemoveFinishedTodos() error {
	c.ListPendingTodos()
	err := c.WriteTodos()
	return err
}

func (c *Collection) Reorder() error {
	for i, todo := range c.Todos {
		todo.Id = int64(i + 1)
	}
	err := c.WriteTodos()
	return err
}

func (c *Collection) Swap(idA int64, idB int64) error {
	var positionA int
	var positionB int

	for i, todo := range c.Todos {
		switch todo.Id {
		case idA:
			positionA = i
			todo.Id = idB
		case idB:
			positionB = i
			todo.Id = idA
		}
	}

	c.Todos[positionA], c.Todos[positionB] = c.Todos[positionB], c.Todos[positionA]
	err := c.WriteTodos()
	return err
}

func (c *Collection) Search(sentence string) {
	sentence = regexp.QuoteMeta(sentence)
	re := regexp.MustCompile("(?i)" + sentence)
	for i := len(c.Todos) - 1; i >= 0; i-- {
		if !re.MatchString(c.Todos[i].Desc) {
			c.RemoveAtIndex(i)
		}
	}
}

func (c *Collection) ReorderByIDs(ids []int64) error {
	idsMap := map[int64]int{}
	for index, id := range ids {
		if _, ok := idsMap[id]; ok {
			return fmt.Errorf("The ID %d is already in the list", id)
		}
		idsMap[id] = index
	}

	ordered := make([]*Todo, len(ids))
	rest := []*Todo{}

	for _, todo := range c.Todos {
		if index, ok := idsMap[todo.Id]; ok {
			ordered[index] = todo
			continue
		}
		rest = append(rest, todo)
	}

	newTodos := make([]*Todo, len(c.Todos))
	index := 0
	var idCounter int64 = 1
	for _, todo := range ordered {
		if todo == nil {
			continue
		}
		todo.Id = idCounter
		newTodos[index] = todo
		index++
		idCounter++
	}
	for _, todo := range rest {
		todo.Id = idCounter
		newTodos[index] = todo
		index++
		idCounter++
	}

	c.Todos = newTodos

	return c.WriteTodos()
}
