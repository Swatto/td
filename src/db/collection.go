package db

import (
	"errors"
	"regexp"
	"strconv"
	"time"
	"umutsevdi/td/parser"
	"umutsevdi/td/todo"
)

type Collection struct {
	Todos []*todo.Todo
}

type STATUS string

const (
	STATUS_PENDING = STATUS("pending")
	STATUS_EXPIRED = STATUS("expired")
	STATUS_DONE    = STATUS("done")
)

func (c *Collection) RemoveAtIndex(item int) {
	s := *c
	s.Todos = append(s.Todos[:item], s.Todos[item+1:]...)
	*c = s
}

func (c *Collection) Has(id int) bool {
	if id >= 0 && id <= int(c.Todos[len(c.Todos)-1].ID) {
		for _, todo := range c.Todos {
			if id == int(todo.ID) {
				return true
			}
		}
	}
	return false
}

func (c *Collection) Find(id int64) (*todo.Todo, error) {
	for _, todo := range c.Todos {
		if id == todo.ID {
			return todo, nil
		}
	}
	return nil, errors.New("The todo with the id " + strconv.FormatInt(id, 10) + " was not found.")
}

func (c *Collection) List(status STATUS) {
	// if STATUS_EXPIRED check_date
	// if STATUS_PENDING check_date && check_not_done
	// if STATUS_DONE check_done
	if status == STATUS_DONE {
		for i := len(c.Todos) - 1; i >= 0; i-- {
			if c.Todos[i].Status != string(STATUS_DONE) {
				c.RemoveAtIndex(i)
			}
		}
	} else if status == STATUS_EXPIRED {
		for i := len(c.Todos) - 1; i >= 0; i-- {
			if c.Todos[i].Deadline.IsZero() || c.Todos[i].Deadline.After(time.Now()) {
				c.RemoveAtIndex(i)
			}
		}
	} else {
		for i := len(c.Todos) - 1; i >= 0; i-- {
			if c.Todos[i].Status != string(STATUS_PENDING) ||
				(!c.Todos[i].Deadline.IsZero() && c.Todos[i].Deadline.Before(time.Now())) {
				c.RemoveAtIndex(i)
			}
		}
	}
}

func (c *Collection) CreateTodo(desc string, d time.Time, p int) *todo.Todo {
	newTodo := &todo.Todo{
		ID:       0,
		Desc:     desc,
		Status:   "pending",
		Modified: "",
		Deadline: d,
		Period:   p,
		Created:  time.Now(),
	}
	var highestID int64 = 0
	for _, todo := range c.Todos {
		if todo.ID > highestID {
			highestID = todo.ID
		}
	}
	newTodo.ID = (highestID + 1)
	newTodo.Modified = time.Now().Local().String()
	c.Todos = append(c.Todos, newTodo)
	return newTodo
}

func (c *Collection) ModifyTodo(id int64, m *map[string]string) (*todo.Todo, error) {
	todo, err := c.Find(id)
	if err != nil {
		return nil, err
	}
	if _, ok := (*m)["desc"]; ok {
		todo.Desc = (*m)["desc"]
	}
	if _, ok := (*m)["date"]; ok {
		todo.Deadline, _ = parser.ParseDate((*m)["date"])
	}
	if _, ok := (*m)["period"]; ok {
		todo.Period, _ = parser.ParsePeriod((*m)["period"])
	}
    todo.Modified=time.Now().Local().String()
	return todo, nil
}

func (c *Collection) Toggle(id int64) (*todo.Todo, error) {
	todo, err := c.Find(id)
	if err != nil {
		return todo, err
	}
	if todo.Period > 0 {
		todo.Period--
		todo.Modified = time.Now().Local().String()
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
	return todo, err
}

func (c *Collection) Remove(id int) error {
	if !c.Has(id) {
		return errors.New("The todo with the id " + strconv.Itoa(id) + "was not found.")
	}
	c.RemoveAtIndex(id)
	return nil
}

func (c *Collection) Reorder() {
	for i, todo := range c.Todos {
		todo.ID = int64(i + 1)
	}
}

func (c *Collection) Swap(idA int64, idB int64) error {
	_, err := c.Find(idA)
	if err != nil {
		return err
	}

	_, err = c.Find(idB)
	if err != nil {
		return err
	}

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
