package db

import (
	"errors"
	"github.com/swatto/td/parser"
	"github.com/swatto/td/todo"
	"regexp"
	"time"
)

// Group of [todo.Todo]
type Collection struct {
	Todos []*todo.Todo // underlying array
	m     map[int]int  // id - array index map
}

// STATUS is an enum for todo statuses
type STATUS string

const (
	RECENT_DATE    = 10
	STATUS_PENDING = STATUS("pending")
	STATUS_EXPIRED = STATUS("expired")
	STATUS_DONE    = STATUS("done")
	NOT_FOUND      = "todo with the given id was not found."
)

func (c *Collection) deleteByIndex(item int) {
	s := *c
	s.Todos = append(s.Todos[:item], s.Todos[item+1:]...)
	delete(c.m, c.Todos[item].ID)
	*c = s
}

// Fetches the internal map to update the cache.
//
// Use after modifying the array directly without the existing functions.
// Built-in functions such as [*Collection.Add] or [Remove] does not require,
// they update the map internally.
func (c *Collection) FetchMap() {
	m := make(map[int]int, 2*len(c.Todos))
	for k, v := range c.Todos {
		m[v.ID] = k
	}
	c.m = m
}

// Get index of the element in the array. If the element does not exist returns
// -1 instead.
func (c *Collection) GetIndex(id int) int {
	t, ok := c.m[id]
	if !ok {
		return -1
	}
	return t
}

// Returns whether an item with given id exists or not.
func (c *Collection) Has(id int) bool {
	if ix := c.GetIndex(id); ix != -1 {
		return c.Todos[ix] != nil
	}
	return false
}

// Returns a pointer to the element with given id. An error is returned when the
// item is not found.
func (c *Collection) Find(id int) (*todo.Todo, error) {
	if ix := c.GetIndex(id); ix != -1 {
		return c.Todos[ix], nil
	}
	return nil, errors.New(NOT_FOUND)
}

// Filters the collection by the STATUS. Optionally can filter items to last 15 days.
//
//   - STATUS_DONE    : List of items that are done.
//   - STATUS_PENDING : List of items that haven't been completed and their
//
// deadline hasn't arrived or doesn't exist.
//   - STATUS_EXPIRED : List of items that haven't been completed in given time.
func (c *Collection) List(status STATUS, isRecent bool) {
	if status == STATUS_DONE {
		for i := len(c.Todos) - 1; i >= 0; i-- {
			if c.Todos[i].Status != string(STATUS_DONE) {
				c.deleteByIndex(i)
			}
		}
	} else if status == STATUS_EXPIRED {
		for i := len(c.Todos) - 1; i >= 0; i-- {
			if !c.Todos[i].Deadline.IsZero() && c.Todos[i].Deadline.After(time.Now()) {
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
	if isRecent {
		for i := len(c.Todos) - 1; i >= 0; i-- {
			if c.Todos[i].Modified.Before(time.Now().AddDate(0, 0, -RECENT_DATE)) {
				c.deleteByIndex(i)
			}

		}
	}
}

// Creates a new [todo.Todo] with given values and adds to the collection.
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
	newTodo.Modified = time.Now()
	c.Todos = append(c.Todos, newTodo)
	c.m[newTodo.ID] = len(c.Todos) - 1
	return newTodo
}

// Modifies the item with given index based on given map, returns the address of
// the updated todo.
// Modification map may have following values:
//   - desc   : Description
//   - date   : A string that is valid according to [parser.ParseDate]
//   - period : A string that contains the period
//
// Any other string is ignored.
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
	todo.Modified = time.Now()
	return todo, nil
}

// Toggles given todo task as done.
//
// If the given task is:
//   - STATUS_DONE    it becomes STATUS_PENDING.
//   - STATUS_PENDING and it has no period it becomes STATUS_DONE.
//   - STATUS_PENDING and it has period, the period is reduced by one.
func (c *Collection) Toggle(id int) (*todo.Todo, error) {
	todo, err := c.Find(id)
	if err != nil {
		return todo, err
	}
	if todo.Period > 0 {
		todo.Period--
		todo.Modified = time.Now()
		return todo, nil
	}

	if todo.Status == string(STATUS_DONE) {
		todo.Status = string(STATUS_PENDING)
	} else {
		todo.Status = string(STATUS_DONE)
	}
	todo.Modified = time.Now()
	return todo, err
}

// Removes the item with given id from the collection.
func (c *Collection) Remove(id int) error {
	index := c.GetIndex(id)
	if index == -1 {
		return errors.New(NOT_FOUND)
	}
	c.deleteByIndex(index)
	return nil
}

// Updates the ids of all elements in the collection.
func (c *Collection) Reorder() {
	for i, todo := range c.Todos {
		todo.ID = i + 1
	}
	c.FetchMap()
}

// Swaps the positions of the items with given ids.
func (c *Collection) Swap(idA int, idB int) error {
	if !c.Has(idA) || !c.Has(idB) {
		return errors.New(NOT_FOUND)
	}

	indexA := c.GetIndex(idA)
	indexB := c.GetIndex(idB)

	c.m[c.Todos[indexA].ID], c.m[c.Todos[indexB].ID] = indexB, indexA
	c.Todos[indexA], c.Todos[indexB] = c.Todos[indexB], c.Todos[indexA]
	return nil
}

// Searches a text pattern in todos
func (c *Collection) Search(sentence string) {
	sentence = regexp.QuoteMeta(sentence)
	re := regexp.MustCompile("(?i)" + sentence)
	for i := len(c.Todos) - 1; i >= 0; i-- {
		if !re.MatchString(c.Todos[i].Desc) {
			c.deleteByIndex(i)
		}
	}
}

// Returns todos with deadlines before the given date
func (c *Collection) FilterByDate(date time.Time) {
	for i, v := range c.Todos {
		if !v.Deadline.IsZero() && v.Deadline.Before(date) {
			c.deleteByIndex(i)
		}
	}
}
