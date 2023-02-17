package test

import "umutsevdi/td/todo"

func ExampleTodo() {
	todo := todo.Todo{
		ID:       0,
		Desc:     "Test td",
		Status:   "pending",
		Modified: "",
	}
	todo.MakeOutput(false)
	// Output: 0 | ✕ Test td
}
