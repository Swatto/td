package main

import (
	"fmt"
	ct "github.com/daviddengcn/go-colortext"
	cli "github.com/urfave/cli/v2"
	"os"
	"strconv"
	"strings"
	"time"
	"umutsevdi/td/db"
	"umutsevdi/td/todo"
)

func WriteError(v ...string) {
	ct.ChangeColor(ct.Red, false, ct.None, false)
	fmt.Print("Error: ")
	ct.ResetColor()
	fmt.Println(strings.Join(v, "\n"))
}

func tdList(c *cli.Context) error {
	var err error
	collection := db.Collection{}

	err = collection.RetrieveTodos()
	if err != nil {
		fmt.Println(err)
	} else {
		if !c.IsSet("all") {
			if c.IsSet("done") {
				collection.ListDoneTodos()
			} else {
				collection.ListPendingTodos()
			}
		}

		if len(collection.Todos) > 0 {
			fmt.Println()
			for _, todo := range collection.Todos {
				todo.MakeOutput(true)
			}
			fmt.Println()
		} else {
			ct.ChangeColor(ct.Cyan, false, ct.None, false)
			fmt.Println("There's no todo to show.")
			ct.ResetColor()
		}
	}
	return nil
}

func tdInit(c *cli.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("%s .\n", err)
		return err
	}

	err = db.CreateStoreFileIfNeeded(cwd + "/.todos")
	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	if err != nil {
		fmt.Printf("A \".todos\" file already exist in \"%s\".\n", cwd)
	} else {
		fmt.Printf("A \".todos\" file is now added to \"%s\".\n", cwd)
	}
	ct.ResetColor()
	return nil
}

func tdAdd(c *cli.Context) error {
	if c.Args().Len() == 0 {
		WriteError("You must provide a name to your todo.",
			"Example: td add \"call mum\"")
		return argError
	}

	collection := db.Collection{}
	dt, dtErr := parseDate(c.Args().Get(1))
	if dtErr != nil {
		WriteError("Invalid date time format\nAvailable Formats: \"[dd/MM/yyyy] [hh:mm]\"")
		return argError
	}
	p, _ := strconv.Atoi(c.Args().Get(2))
	todo := todo.Todo{
		ID:       0,
		Desc:     c.Args().Get(0),
		Status:   "pending",
		Modified: "",
		Deadline: dt,
		Period:   p,
		Created:  time.Now(),
	}
	ct.ChangeColor(ct.Blue, false, ct.None, false)
	if !dt.IsZero() {
		fmt.Println("Deadline added")
	}
	if p != 0 {
		fmt.Println("Periodic to do created")
	}
	ct.ResetColor()
	err := collection.RetrieveTodos()
	if err != nil {
		fmt.Println(err)
		return err
	}

	id, err := collection.CreateTodo(&todo)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	fmt.Printf("#%d \"%s\" is now added to your todos.\n", id, c.Args().Get(0))
	ct.ResetColor()
	return nil
}

func tdModify(c *cli.Context) error {
	if c.Args().Len() < 2 {
		WriteError("You must provide the id and the new text for your todo.")
		fmt.Println("Example: td modify 2 \"call dad\"")
		return argError
	}

	collection := db.Collection{}
	collection.RetrieveTodos()

	args := c.Args()

	id, err := strconv.ParseInt(args.Get(0), 10, 32)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	td, err := collection.Find(id)
	if err != nil {
		WriteError(err.Error())
		return err
	}
	err = FindAndReplace(args, td)
	if err != nil {
		WriteError(err.Error())
		return err
	}
	err = collection.WriteTodos()
	if err != nil {
		fmt.Println(err)
		return err
	}

	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	fmt.Println(args.Get(0), "has been updated:")
	ct.ResetColor()
	td.MakeOutput(true)
	return nil
}

func tdToggle(c *cli.Context) error {
	var err error

	if c.Args().Len() != 1 {
		WriteError("You must provide the position of the item you want to change.",
			"Example: td toggle 1")
		return argError
	}

	collection := db.Collection{}
	collection.RetrieveTodos()

	id, err := strconv.ParseInt(c.Args().Get(0), 10, 32)
	if err != nil {
		fmt.Println(err)
		return err
	}

	todo, err := collection.Toggle(id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	fmt.Printf("Your todo is now %s.\n", todo.Status)
	ct.ResetColor()
	return nil
}

func tdClean(c *cli.Context) error {
	collection := db.Collection{}
	collection.RetrieveTodos()

	err := collection.RemoveFinishedTodos()

	if err != nil {
		fmt.Println(err)
		return err
	}

	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	fmt.Println("Your list is now flushed of finished todos.")
	ct.ResetColor()
	return nil
}
func tdReorder(c *cli.Context) error {
	collection := db.Collection{}
	collection.RetrieveTodos()

	if c.Args().Len() != 1 {
		WriteError("You must provide two position if you want to swap todos.",
			"Example: td reorder 9 3")
		return argError
	} else if c.Args().Len() != 2 {
		idA, err := strconv.ParseInt(c.Args().Get(0), 10, 32)
		if err != nil {
			fmt.Println(err)
			return err
		}

		idB, err := strconv.ParseInt(c.Args().Get(1), 10, 32)
		if err != nil {
			fmt.Println(err)
			return err
		}

		_, err = collection.Find(idA)
		if err != nil {
			fmt.Println(err)
			return err
		}

		_, err = collection.Find(idB)
		if err != nil {
			fmt.Println(err)
			return err
		}

		collection.Swap(idA, idB)

		ct.ChangeColor(ct.Cyan, false, ct.None, false)
		fmt.Printf("\"%s\" and \"%s\" has been swapped\n", c.Args().Get(0), c.Args().Get(1))
		ct.ResetColor()
	}

	err := collection.Reorder()

	if err != nil {
		fmt.Println(err)
		return err
	}

	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	fmt.Println("Your list is now reordered.")
	ct.ResetColor()
	return nil
}
func tdSearch(c *cli.Context) error {
	if c.Args().Len() != 1 {
		WriteError("You must provide a string earch.",
			"Example: td search \"project-1\"")
		return argError
	}

	collection := db.Collection{}
	collection.RetrieveTodos()
	collection.Search(c.Args().Get(0))

	if len(collection.Todos) == 0 {
		ct.ChangeColor(ct.Cyan, false, ct.None, false)
		fmt.Printf("Sorry, there's no todos containing \"%s\".\n", c.Args().Get(0))
		ct.ResetColor()
		return argError
	}

	if len(collection.Todos) > 0 {
		fmt.Println()
		for _, todo := range collection.Todos {
			todo.MakeOutput(true)
		}
		fmt.Println()
	} else {
		ct.ChangeColor(ct.Cyan, false, ct.None, false)
		fmt.Println("There's no todo to show.")
		ct.ResetColor()
	}
	return nil
}

func tdSearchByDate(c *cli.Context) error {
	var date time.Time
	if c.Args().Len() < 2 {
		date = time.Now().Add(time.Hour * 24)
	} else {
		var err error
		date, err = parseDate(c.Args().Get(1))
		if err != nil {
			WriteError(err.Error())
			return argError
		}
	}
	collection := db.Collection{}
	collection.RetrieveTodos()
	for _, v := range collection.Todos {
		if !v.Deadline.IsZero() && v.Deadline.Before(date) {
			v.MakeOutput(true)
		}
	}
	return nil
}
