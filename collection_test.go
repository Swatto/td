package main

import (
	"fmt"
	"testing"
)

var statuses = []string{PENDING, DONE, WIP}

func TestListPendingTodos(t *testing.T) {
	var collection Collection
	var todos []*Todo

	for id := range make([]int, len(statuses)) {
		task := NewTodo()
		task.Id = int64(id)
		task.Status = statuses[id]
		task.Desc = fmt.Sprintf("This is task number %d as %s", task.Id, task.Status)
		todos = append(todos, task)
	}
	collection.Todos = todos
	collection.ListPendingTodos()

	if len(collection.Todos) != 1 {
		t.Error("Expected only one task, got", len(collection.Todos))
	}

	for _, td := range collection.Todos {
		if td.Status != PENDING {
			t.Errorf("Expected status of tasks to be only \"pending\", got at least one as \"%s\"", td.Status)
		}
	}
}

func TestListUndoneTodos(t *testing.T) {
	var collection Collection
	var todos []*Todo

	for id := range make([]int, len(statuses)) {
		task := NewTodo()
		task.Id = int64(id)
		task.Status = statuses[id]
		todos = append(todos, task)
	}
	collection.Todos = todos
	collection.ListUndoneTodos()

	if len(collection.Todos) != 2 {
		t.Error("Expected only one task, got", len(collection.Todos))
	}

	for _, td := range collection.Todos {
		if td.Status != PENDING && td.Status != WIP {
			t.Errorf("Expected status of tasks to be only \"pending\" or \"work in progress\", got at least one as \"%s\"", td.Status)
		}
	}
}

func TestListWorkInProgressTodos(t *testing.T) {
	var collection Collection
	var todos []*Todo

	for id := range make([]int, len(statuses)) {
		task := NewTodo()
		task.Id = int64(id)
		task.Status = statuses[id]
		task.Desc = fmt.Sprintf("This is task number %d as %s", task.Id, task.Status)
		todos = append(todos, task)
	}
	collection.Todos = todos
	collection.ListWorkInProgressTodos()

	if len(collection.Todos) != 1 {
		t.Error("Expected only one task, got", len(collection.Todos))
	}

	for _, td := range collection.Todos {
		if td.Status != WIP {
			t.Errorf("Expected status of tasks to be only \"work in progress\", got at least one as \"%s\"", td.Status)
		}
	}
}

func TestListDoneTodos(t *testing.T) {
	var collection Collection
	var todos []*Todo

	for id := range make([]int, len(statuses)) {
		task := NewTodo()
		task.Id = int64(id)
		task.Status = statuses[id]
		todos = append(todos, task)
	}
	collection.Todos = todos
	collection.ListDoneTodos()

	if len(collection.Todos) != 1 {
		t.Error("Expected only one task, got", len(collection.Todos))
	}

	for _, td := range collection.Todos {
		if td.Status != DONE {
			t.Errorf("Expected status of tasks to be only \"done\", got at least one as \"%s\"", td.Status)
		}
	}
}

func TestToggleStatus(t *testing.T) {
	// PENDING > WIP > DONE > PENDING
	var collection Collection
	var todos []*Todo

	task := NewTodo()
	task.Id = 1
	task.Status = PENDING
	todos = append(todos, task)
	collection.Todos = todos

	collection.Toggle(1)
	if collection.Todos[0].Status != WIP {
		t.Errorf("Expected status to go from \"pending\" to \"work in progress\", changed to \"%s\" instead.", collection.Todos[0].Status)
	}

	collection.Toggle(1)
	if collection.Todos[0].Status != DONE {
		t.Errorf("Expected status to go from \"work in progress\" to \"done\", changed to \"%s\" instead.", collection.Todos[0].Status)
	}

	collection.Toggle(1)
	if collection.Todos[0].Status != PENDING {
		t.Errorf("Expected status to go from \"done\" to \"pending\", changed to \"%s\" instead.", collection.Todos[0].Status)
	}
}

func TestSetStatus(t *testing.T) {
	var collection Collection
	var todos []*Todo

	task := NewTodo()
	task.Id = 1
	task.Status = PENDING
	todos = append(todos, task)
	collection.Todos = todos

	collection.SetStatus(1, DONE)
	if collection.Todos[0].Status != DONE {
		t.Errorf("Expected to set status to \"done\", got \"%s\" instead.", collection.Todos[0].Status)
	}
}

func TestTodoModifyDescription(t *testing.T) {
	var collection Collection
	var todos []*Todo

	oldDesc := []string{
		"Test 1",
		"Test 2",
		"Test 3",
		"Test 4",
	}

	newDesc := []string{
		"New test 1",
		"New test 2",
		"New test 3",
		"New test 4",
	}

	for id := range make([]int, len(oldDesc)) {
		task := NewTodo()
		task.Id = int64(id)
		task.Desc = oldDesc[id]
		todos = append(todos, task)
	}
	collection.Todos = todos

	for id, td := range collection.Todos {
		if td.Desc != oldDesc[id] {
			t.Error("Something is wrong with the test, description should be", oldDesc[id])
			t.FailNow()
		}
		collection.Modify(int64(id), newDesc[id])
		if td.Desc != newDesc[id] {
			t.Error("Something is wrong with the test, description should be", newDesc[id])
			t.FailNow()
		}
	}
}

func TestTodoSwap(t *testing.T) {
	var collection Collection
	var todos []*Todo

	taskDesc := []string{
		"Test 1",
		"Test 2",
		"Test 3",
		"Test 4",
	}

	for id := range make([]int, len(taskDesc)) {
		task := NewTodo()
		task.Id = int64(id + 1)
		task.Desc = taskDesc[id]
		todos = append(todos, task)
	}
	collection.Todos = todos
	collection.Swap(2, 4)

	secondTodo, err := collection.Find(2)
	if err != nil {
		t.Errorf("Expect to find item 2, but this happened: %s", err)
		t.FailNow()
	}
	if secondTodo.Desc != taskDesc[3] {
		t.Errorf("Expected the description from second todo item to have be \"%s\", but it didn't", taskDesc[3])
		t.FailNow()
	}

	lastTodo, err := collection.Find(4)
	if err != nil {
		t.Errorf("Expect to find item 4, but this happened: %s", err)
		t.FailNow()
	}
	if lastTodo.Desc != taskDesc[1] {
		t.Errorf("Expected the description from last todo item to have be \"%s\", but it didn't", taskDesc[1])
	}

}

func TestRemoveAtIndex(t *testing.T) {
	var collection Collection
	var todos []*Todo

	taskDesc := []string{
		"Test 1",
		"Test 2",
		"Test 3",
		"Test 4",
	}

	for id := range make([]int, len(taskDesc)) {
		task := NewTodo()
		task.Id = int64(id + 1)
		task.Desc = taskDesc[id]
		todos = append(todos, task)
	}
	collection.Todos = todos
	collection.RemoveAtIndex(2)
	if len(collection.Todos) != len(todos)-1 {
		t.Errorf("Expected size of current to-dos to one less then original slice, but it was %d and the other %d", len(collection.Todos), len(todos))
	}
}
