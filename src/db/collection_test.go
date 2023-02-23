package db

import (
	"testing"
	"time"
	"umutsevdi/td/todo"
)

func _TestCollectionIds() []int {
	return []int{0, 1, 2, 4, 5, 6}
}
func _TestMakeCollection() Collection {
	return Collection{Todos: []*todo.Todo{
		// no deadline, pending
		{
			ID:       0,
			Status:   "pending",
			Created:  time.Now().AddDate(0, -1, 0),
			Modified: time.Now().AddDate(0, -1, 0).Local().String(),
		},
		// no deadline done
		{
			ID:       1,
			Status:   "done",
			Created:  time.Now().AddDate(0, -1, 0),
			Modified: time.Now().AddDate(0, -1, 0).Local().String(),
		},
		// deadline pending not-expired
		{
			ID:       2,
			Status:   "pending",
			Created:  time.Now().AddDate(0, -1, 0),
			Modified: time.Now().AddDate(0, -1, 0).Local().String(),
			Deadline: time.Now().AddDate(0, 1, 0),
		},
		// deadline done not-expired
		{
			ID:       4,
			Status:   "done",
			Created:  time.Now().AddDate(0, -1, 0),
			Modified: time.Now().AddDate(0, -1, 0).Local().String(),
			Deadline: time.Now().AddDate(0, 1, 0),
		},
		// deadline pending expired
		{
			ID:       5,
			Status:   "pending",
			Created:  time.Now(),
			Modified: time.Now().Local().String(),
			Deadline: time.Now().AddDate(0, -1, 0),
		},
		// deadline done expired
		{
			ID:       6,
			Status:   "done",
			Created:  time.Now().AddDate(0, -1, 0),
			Modified: time.Now().AddDate(0, -1, 0).Local().String(),
			Deadline: time.Now().AddDate(0, -1, 0),
		},
	}}
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
		if _, err := c.Find(int64(i)); err != nil {
			t.Fail()
			t.Log("t doesn't have ", i)
		}
	}
	if _, err := c.Find(int64(-1)); err == nil {
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
	c.List(STATUS_DONE)
	if len(c.Todos) != 3 {
		t.Fail()
		t.Log("STATUS_DONE should be 3")
	} else {
		for _, v := range c.Todos {
			if v.Status != string(STATUS_DONE) {
				t.Fail()
				t.Log(v.ID, " is not STATUS_DONE")
			}
		}
	}

	c = _TestMakeCollection()
	c.List(STATUS_EXPIRED)
	if len(c.Todos) != 2 {
		t.Fail()
		t.Log("STATUS_EXPIRED should be 2")
	}

	c = _TestMakeCollection()
	c.List(STATUS_PENDING)
	if len(c.Todos) != 2 {
		t.Fail()
		t.Log("STATUS_PENDING should be 2")
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
	todo := c.CreateTodo("description", time.Time{}, 0)
	if !c.Has(int(todo.ID)) {
		t.Fail()
		t.Log("todo should be created")
	}
}

func TestModifyTodo(t *testing.T) {
	c := _TestMakeCollection()
	c.ModifyTodo(0, &map[string]string{"desc": "new description"})
	c.ModifyTodo(0, &map[string]string{
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
