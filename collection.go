package main

import (
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strconv"
	"time"
)

const WIP = "wip"
const DONE = "done"
const PENDING = "pending"

type Collection struct {
	Todos []*Todo
}

func NewCollection() (*Collection, error) {
	var collection *Collection
	collection = new(Collection)

	if err := collection.RetrieveTodos(); err != nil {
		return nil, err
	}

	return collection, nil
}

func (c *Collection) RemoveAtIndex(item int) {
	s := *c
	s.Todos = append(s.Todos[:item], s.Todos[item+1:]...)
	*c = s
}

func (c *Collection) RetrieveTodos() error {
	db, err := NewDataStore()
	if err != nil {
		return err
	}
	if err := db.Check(); err != nil {
		return err
	}

	file, err := os.OpenFile(db.Path, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&c.Todos)
	return err
}

func (c *Collection) WriteTodos() error {
	db, err := NewDataStore()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(db.Path, os.O_RDWR|os.O_TRUNC, 0600)
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

func (c *Collection) ListPendingTodos() error {
	for i := len(c.Todos) - 1; i >= 0; i-- {
		if c.Todos[i].Status != PENDING {
			c.RemoveAtIndex(i)
		}
	}
	return nil
}

func (c *Collection) ListUndoneTodos() error {
	for i := len(c.Todos) - 1; i >= 0; i-- {
		if c.Todos[i].Status == DONE {
			c.RemoveAtIndex(i)
		}
	}
	return nil
}

func (c *Collection) ListDoneTodos() error {
	for i := len(c.Todos) - 1; i >= 0; i-- {
		if c.Todos[i].Status != DONE {
			c.RemoveAtIndex(i)
		}
	}
	return nil
}

func (c *Collection) ListWorkInProgressTodos() {
	for i := len(c.Todos) - 1; i >= 0; i-- {
		if c.Todos[i].Status != WIP {
			c.RemoveAtIndex(i)
		}
	}
}

func (c *Collection) CreateTodo(newTodo *Todo) (int64, error) {
	var err error
	var highestId int64 = 0

	for _, todo := range c.Todos {
		if todo.Id > highestId {
			highestId = todo.Id
		}
	}

	newTodo.Id = (highestId + 1)
	newTodo.Modified = time.Now().Local().String()
	c.Todos = append(c.Todos, newTodo)

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

func (c *Collection) SetStatus(id int64, status string) (*Todo, error) {
	todo, err := c.Find(id)

	if err != nil {
		return todo, err
	}

	todo.Status = status
	todo.Modified = time.Now().Local().String()

	return todo, err
}

func (c *Collection) Toggle(id int64) (*Todo, error) {
	var status string

	todo, err := c.Find(id)

	if err != nil {
		return todo, err
	}

	switch todo.Status {
	case PENDING:
		status = WIP
	case WIP:
		status = DONE
	default:
		status = PENDING
	}

	return c.SetStatus(id, status)
}

func (c *Collection) Modify(id int64, desc string) (*Todo, error) {
	todo, err := c.Find(id)

	if err != nil {
		return todo, err
	}

	todo.Desc = desc
	todo.Modified = time.Now().Local().String()

	return todo, err
}

func (c *Collection) RemoveFinishedTodos() error {
	return c.ListPendingTodos()
}

func (c *Collection) Reorder() error {
	for i, todo := range c.Todos {
		todo.Id = int64(i + 1)
	}
	return nil
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
	return nil
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
