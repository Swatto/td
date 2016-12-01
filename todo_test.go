package main

func ExampleTodo() {
	todo := Todo{
		Id:       0,
		Desc:     "Test td",
		Status:   "pending",
		Modified: "",
	}
	todo.MakeOutput(false)
	// Output: 0 | ✕ Test td
}
