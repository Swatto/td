package db

import (
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strconv"
	"time"
	"umutsevdi/td/todo"
)

type Collection struct {
	Todos []*todo.Todo
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

func (c *Collection) CreateTodo(newTodo *todo.Todo) (int64, error) {
	var highestID int64 = 0
	for _, todo := range c.Todos {
		if todo.ID > highestID {
			highestID = todo.ID
		}
	}

	newTodo.ID = (highestID + 1)
	newTodo.Modified = time.Now().Local().String()
	c.Todos = append(c.Todos, newTodo)

	err := c.WriteTodos()
	return newTodo.ID, err
}

func (c *Collection) Find(id int64) (*todo.Todo, error) {
	for _, todo := range c.Todos {
		if id == todo.ID {
			return todo, nil
		}
	}
	return nil, errors.New("The todo with the id " + strconv.FormatInt(id, 10) + " was not found.")
}

func (c *Collection) Toggle(id int64) (*todo.Todo, error) {
	todo, err := c.Find(id)

	if err != nil {
		return todo, err
	}
	if todo.Period > 0 {
		todo.Period--
		todo.Modified = time.Now().Local().String()
		//		if !todo.Deadline.IsZero() {
		//		todo.Deadline = todo.Deadline.Add(todo.Deadline.Sub(todo.Created))
		//	}
		c.WriteTodos()
		if err != nil {
			err = errors.New("todos couldn't be saved")
			return todo, err
		}

		return todo, nil
	}

	if todo.Status == "done" {
		todo.Status = "pending"
	} else {
		todo.Status = "done"
	}
	todo.Modified = time.Now().Local().String()

	err = c.WriteTodos()
	if err != nil {
		err = errors.New("todos couldn't be saved")
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
		todo.ID = int64(i + 1)
	}
	err := c.WriteTodos()
	return err
}

func (c *Collection) Swap(idA int64, idB int64) error {
	var positionA int
	var positionB int

	for i, todo := range c.Todos {
		switch todo.ID {
		case idA:
			positionA = i
			todo.ID = idB
		case idB:
			positionB = i
			todo.ID = idA
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

func (c *Collection) FilterByDate(date time.Time) {
	for i, v := range c.Todos {
		if !v.Deadline.IsZero() && v.Deadline.Before(date) {
			c.RemoveAtIndex(i)
		}
	}

}
