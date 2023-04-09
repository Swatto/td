package db

import (
	"github.com/swatto/td/todo"
	"testing"
	"time"
)

func _TestCollectionIds() []int {
	return []int{0, 1, 2, 4, 5, 6}
}
func _TestMakeCollection() Collection {
	c := Collection{Todos: []*todo.Todo{
		// no deadline, pending
		{
			ID:       0,
			Status:   "pending",
			Created:  time.Now().AddDate(0, -1, 0),
			Modified: time.Now().AddDate(0, -1, 0),
		},
		// no deadline done
		{
			ID:       1,
			Status:   "done",
			Created:  time.Now().AddDate(0, -1, 0),
			Modified: time.Now().AddDate(0, 0, -1),
		},
		// deadline pending not-expired
		{
			ID:       2,
			Status:   "pending",
			Created:  time.Now().AddDate(0, -1, 0),
			Modified: time.Now().AddDate(0, -1, 0),
			Deadline: time.Now().AddDate(0, 1, 0),
		},
		// deadline done not-expired
		{
			ID:       4,
			Status:   "done",
			Created:  time.Now().AddDate(0, -1, 0),
			Modified: time.Now().AddDate(0, -1, 0),
			Deadline: time.Now().AddDate(0, 1, 0),
		},
		// deadline pending expired
		{
			ID:       5,
			Status:   "pending",
			Created:  time.Now(),
			Modified: time.Now(),
			Deadline: time.Now().AddDate(0, -1, 0),
		},
		// deadline done expired
		{
			ID:       6,
			Status:   "done",
			Created:  time.Now().AddDate(0, -1, 0),
			Modified: time.Now(),
			Deadline: time.Now(),
		},
	}}
	c.FetchMap()
	return c
}

func TestGetIndex(t *testing.T) {
	c := _TestMakeCollection()
	for _, i := range _TestCollectionIds() {
		if c.GetIndex(i) == -1 {
			t.Fail()
			t.Log("t should have ", i)
		}
	}
	if c.GetIndex(-1) != -1 {
		t.Fail()
		t.Log("t shouldn't have ", -1)
	}
	if c.GetIndex(int(c.Todos[len(c.Todos)-1].ID)+1) != -1 {
		t.Fail()
		t.Log("t shouldn't have ", len(c.Todos))
	}
}
func TestHas(t *testing.T) {
	c := _TestMakeCollection()
	for _, i := range _TestCollectionIds() {
		if !c.Has(i) {
			t.Fail()
			t.Log("t should have ", i)
		}
	}
	if c.Has(-1) {
		t.Fail()
		t.Log("t shouldn't have ", -1)
	}
	if c.Has(int(c.Todos[len(c.Todos)-1].ID) + 1) {
		t.Fail()
		t.Log("t shouldn't have ", len(c.Todos))
	}
}

func TestFind(t *testing.T) {
	c := _TestMakeCollection()
	for _, i := range _TestCollectionIds() {
		if _, err := c.Find(i); err != nil {
			t.Fail()
			t.Log("t doesn't have ", i)
		}
	}
	if _, err := c.Find(-1); err == nil {
		t.Fail()
		t.Log("t shouldn't have argument ", -1)
	}
	if _, err := c.Find(c.Todos[len(c.Todos)-1].ID + 1); err == nil {
		t.Fail()
		t.Log("t shouldn't have argument ", len(c.Todos))
	}
}

func TestList(t *testing.T) {
	c := _TestMakeCollection()
	c.List(STATUS_DONE, true)
	if len(c.Todos) != 2 {
		t.Fail()
		t.Log("STATUS_DONE should be 2")
	} else {
		for _, v := range c.Todos {
			if v.Status != string(STATUS_DONE) {
				t.Fail()
				t.Log(v.ID, " is not STATUS_DONE")
			}
		}
	}

	c = _TestMakeCollection()
	c.List(STATUS_EXPIRED, true)
	if len(c.Todos) != 3 {
		t.Fail()
		t.Log("STATUS_EXPIRED should be 3 but is ", len(c.Todos))
	}

	c = _TestMakeCollection()
	c.List(STATUS_PENDING, false)
	if len(c.Todos) != 2 {
		t.Fail()
		t.Log("STATUS_PENDING should be 2 but is ", len(c.Todos))
	} else {
		for _, v := range c.Todos {
			if v.Status != string(STATUS_PENDING) {
				t.Fail()
				t.Log(v.ID, " is not STATUS_PENDING")
			}
		}
	}
}

func TestCreateTodo(t *testing.T) {
	c := _TestMakeCollection()
	todo := c.Add("description", time.Time{}, 0)
	if !c.Has(int(todo.ID)) {
		t.Fail()
		t.Log("todo should be created")
	}
}

func TestModify(t *testing.T) {
	c := _TestMakeCollection()
	c.Modify(0, &map[string]string{"desc": "new description"})
	c.Modify(0, &map[string]string{
		"date":       time.Now().Local().String(),
		"period":     "4",
		"random_key": "random_value"})
	if c.Todos[0].Desc != "new description" {
		t.Fail()
		t.Log("description was not updated")
	} else if c.Todos[0].Deadline.After(time.Now().AddDate(0, 0, -1)) {
		t.Fail()
		t.Log("deadline was not updated")
	} else if c.Todos[0].Period != 4 {
		t.Fail()
		t.Log("period was not updated")
	}
}

func TestToggle(t *testing.T) {
	c := _TestMakeCollection()
	c.Todos[0].Period = 1
	todo, err := c.Toggle(0)
	if err != nil {
		t.Fail()
		t.Log("error in function:", err)
	}
	if todo.Status != string(STATUS_PENDING) || todo.Period != 0 {
		t.Fail()
		t.Log("error at period control", todo)
	}
	todo, err = c.Toggle(0)
	if todo.Status != string(STATUS_DONE) {
		t.Fail()
		t.Log("error at pending->done")
	}
}

func TestRemove(t *testing.T) {
	c := _TestMakeCollection()
	err := c.Remove(-1)
	if err == nil {
		t.Fail()
		t.Log(err)
	}
	ids := _TestCollectionIds()
	for i := 0; i < len(ids)/2; i++ {
		j := ids[len(ids)-1-i]
		ids[len(ids)-1-i] = ids[i]
		ids[i] = j
	}
	t.Log(ids)

	for _, i := range ids {
		err = c.Remove(i)
		if err != nil {
			t.Fail()
			t.Log(err)
		}
		if c.Has(i) {
			t.Fail()
			t.Log("collection shouldn't have ", i, " after deletion")
		}
	}
}

func TestReorder(t *testing.T) {
	c := _TestMakeCollection()
	e1 := c.Todos[0]
	e2 := c.Todos[len(c.Todos)-1]
	c.Todos[0], c.Todos[len(c.Todos)-1] = c.Todos[len(c.Todos)-1], c.Todos[0]
	t.Log(c.Todos)
	c.Reorder()
	t.Log(c.Todos)
	if e2.ID > e1.ID {
		t.Fail()
		t.Log("error at reorder")
	}
}

func TestSwap(t *testing.T) {
	c := _TestMakeCollection()
	e1 := c.Todos[0]
	e2 := c.Todos[len(c.Todos)-1]
	err := c.Swap(e1.ID, e2.ID)
	if err != nil {
		t.Fail()
		t.Log("ids should exist")
	}

	if c.GetIndex(e1.ID) != len(c.Todos)-1 || c.GetIndex(e2.ID) != 0 {
		t.Fail()
		t.Log("error in swap")
		t.Log(c.GetIndex(e1.ID), ",", c.GetIndex(e2.ID))
	}

	err = c.Swap(4000, -1000)
	if err == nil {
		t.Fail()
		t.Log("ids shouldn't exist")
	}
}
