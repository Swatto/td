package todo

import "time"

func ExampleTodo() {
	todo := Todo{
		ID:       0,
		Desc:     "Test td",
		Status:   "pending",
		Modified: time.Time{},
	}
	todo.MakeOutput(false)
	// Output: 00 | ✕ Test td
}
