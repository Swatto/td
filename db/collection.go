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
	m     map[int]int
}

type STATUS string

const (
	STATUS_PENDING = STATUS("pending")
	STATUS_EXPIRED = STATUS("expired")
	STATUS_DONE    = STATUS("done")
)

func (c *Collection) FetchMap() {
	m := make(map[int]int, 2*len(c.Todos))
	for k, v := range c.Todos {
		m[v.ID] = k
	}
	c.m = m
}

func (c *Collection) deleteByIndex(item int) {
	s := *c
	s.Todos = append(s.Todos[:item], s.Todos[item+1:]...)
	delete(c.m, c.Todos[item].ID)
	*c = s
}

func (c *Collection) GetIndex(id int) int {
	t, ok := c.m[id]
	if !ok {
		return -1
	}
	return t
}

func (c *Collection) Has(id int) bool {
	if ix := c.GetIndex(id); ix != -1 {
		return c.Todos[ix] != nil
	}
	return false
}

func (c *Collection) Find(id int) (*todo.Todo, error) {
	if ix := c.GetIndex(id); ix != -1 {
		return c.Todos[ix], nil
	}
	return nil, errors.New("The todo with the id " + strconv.FormatInt(int64(id), 10) + " was not found.")
}

func (c *Collection) List(status STATUS) {
	if status == STATUS_DONE {
		for i := len(c.Todos) - 1; i >= 0; i-- {
			if c.Todos[i].Status != string(STATUS_DONE) {
				c.deleteByIndex(i)
			}
		}
	} else if status == STATUS_EXPIRED {
		for i := len(c.Todos) - 1; i >= 0; i-- {
			if c.Todos[i].Deadline.IsZero() || c.Todos[i].Deadline.After(time.Now()) {
				c.deleteByIndex(i)
			}
		}
	} else {
		for i := len(c.Todos) - 1; i >= 0; i-- {
			if c.Todos[i].Status != string(STATUS_PENDING) ||
				(!c.Todos[i].Deadline.IsZero() && c.Todos[i].Deadline.Before(time.Now())) {
				c.deleteByIndex(i)
			}
		}
	}
}

func (c *Collection) Add(desc string, d time.Time, p int) *todo.Todo {
	newTodo := &todo.Todo{
		Desc:     desc,
		Status:   string(STATUS_PENDING),
		Deadline: d,
		Period:   p,
		Created:  time.Now(),
	}
	var highestID int = 0
	for _, todo := range c.Todos {
		if todo.ID > highestID {
			highestID = todo.ID
		}
	}
	newTodo.ID = (highestID + 1)
	newTodo.Modified = time.Now().Local().String()
	c.Todos = append(c.Todos, newTodo)
	c.m[newTodo.ID] = len(c.Todos) - 1
	return newTodo
}

func (c *Collection) Modify(id int, m *map[string]string) (*todo.Todo, error) {
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
	todo.Modified = time.Now().Local().String()
	return todo, nil
}

func (c *Collection) Toggle(id int) (*todo.Todo, error) {
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

	if todo.Status == string(STATUS_DONE) {
		todo.Status = string(STATUS_PENDING)
	} else {
		todo.Status = string(STATUS_DONE)
	}
	todo.Modified = time.Now().Local().String()
	return todo, err
}

func (c *Collection) Remove(id int) error {
	index := c.GetIndex(id)
	if index == -1 {
		return errors.New("The todo with the id " + strconv.Itoa(id) + "was not found.")
	}
	c.deleteByIndex(index)
	return nil
}

func (c *Collection) Reorder() {
	for i, todo := range c.Todos {
		todo.ID = i + 1
	}
	c.FetchMap()
}

func (c *Collection) Swap(idA int, idB int) error {
	if !c.Has(idA) || !c.Has(idB) {
		return errors.New("No such todo")
	}

	indexA := c.GetIndex(idA)
	indexB := c.GetIndex(idB)

	c.m[c.Todos[indexA].ID], c.m[c.Todos[indexB].ID] = indexB, indexA
	c.Todos[indexA], c.Todos[indexB] = c.Todos[indexB], c.Todos[indexA]
	return nil
}

func (c *Collection) Search(sentence string) {
	sentence = regexp.QuoteMeta(sentence)
	re := regexp.MustCompile("(?i)" + sentence)
	for i := len(c.Todos) - 1; i >= 0; i-- {
		if !re.MatchString(c.Todos[i].Desc) {
			c.deleteByIndex(i)
		}
	}
}

func (c *Collection) FilterByDate(date time.Time) {
	for i, v := range c.Todos {
		if !v.Deadline.IsZero() && v.Deadline.Before(date) {
			c.deleteByIndex(i)
		}
	}

}
