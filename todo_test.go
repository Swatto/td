package main

func Exampletodo() {
	todo := todo{
		ID:       0,
		Desc:     "Test td",
		Status:   "pending",
		Modified: "",
	}
	todo.MakeOutput(false)
	// Output: 0 | ✕ Test td
}
