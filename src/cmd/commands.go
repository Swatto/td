package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"umutsevdi/td/db"
	"umutsevdi/td/parser"

	ct "github.com/daviddengcn/go-colortext"
	cli "github.com/urfave/cli/v2"
)

type Error string

const argError = Error("Error in argument")

func (e Error) Error() string { return string(e) }

func WriteError(v ...string) {
	ct.ChangeColor(ct.Red, false, ct.None, false)
	fmt.Print("Error: ")
	ct.ResetColor()
	fmt.Println(strings.Join(v, "\n"))
}

func TdList(c *cli.Context) error {
	collection, err := db.Read()
	if err != nil {
		WriteError(err.Error())
	} else {
		if !c.IsSet("all") {
			if c.IsSet("done") {
				collection.List(db.STATUS_DONE)
			} else if c.IsSet("past") {
				collection.List(db.STATUS_EXPIRED)
			} else {
				collection.List(db.STATUS_PENDING)
			}
		}

		if len(collection.Todos) > 0 {
			for _, todo := range collection.Todos {
				todo.MakeOutput(true)
			}
		} else {
			ct.ChangeColor(ct.Cyan, false, ct.None, false)
			fmt.Println("There's no todo to show.")
			ct.ResetColor()
		}
	}
	return nil
}

func TdInit(c *cli.Context) error {
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

func TdAdd(c *cli.Context) error {
	if c.Args().Len() == 0 {
		WriteError("You must provide a name to your todo.",
			"Example: Td add \"call mum\"")
		return argError
	}
	dt, dtErr := parser.ParseDate(c.Args().Get(1))
	if dtErr != nil {
		WriteError("Invalid date time format\nAvailable Formats: \"[dd/MM/yyyy] [hh:mm]\"")
		return argError
	}
	p, _ := strconv.Atoi(c.Args().Get(2))

	collection, err := db.Read()
	if err != nil {
		WriteError(err.Error())
		return err
	}

	todo := collection.CreateTodo(c.Args().Get(0), dt, p)

	ct.ChangeColor(ct.Blue, false, ct.None, false)
	if !dt.IsZero() {
		fmt.Println("Deadline added")
	}
	if p != 0 {
		fmt.Println("Periodic added")
	}
	ct.ResetColor()
	err = db.Save(collection)
	if err != nil {
		fmt.Println(err)
		return err
	}
	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	fmt.Printf("#%d \"%s\" is now added to your todos.\n", todo.ID, todo.Desc)
	ct.ResetColor()
	return nil
}

func TdModify(c *cli.Context) error {
	if c.Args().Len() < 2 {
		WriteError("You must provide the id and the new text for your todo.")
		fmt.Println("Example: Td modify 2 \"call dad\"")
		return argError
	}
	args := c.Args()

	id, err := strconv.ParseInt(args.Get(0), 10, 32)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	collection, err := db.Read()
	if err != nil {
		WriteError(err.Error())
		return err
	}
	m, err := parser.MapArgs(args)
	if err != nil {
		WriteError(err.Error())
		return err
	}
	Td, err := collection.ModifyTodo(id, m)
	if err != nil {
		WriteError(err.Error())
		return err
	}
	err = db.Save(collection)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	fmt.Println(args.Get(0), "has been updated:")
	ct.ResetColor()
	Td.MakeOutput(true)
	return nil
}

func TdToggle(c *cli.Context) error {
	if c.Args().Len() != 1 {
		WriteError("You must provide the position of the item you want to change.",
			"Example: Td toggle 1")
		return argError
	}

	collection, err := db.Read()
	if err != nil {
		WriteError(err.Error())
		return err
	}

	id, err := strconv.ParseInt(c.Args().Get(0), 10, 32)
	if err != nil {
		WriteError(err.Error())
		return err
	}

	todo, err := collection.Toggle(id)
	if err != nil {
		WriteError(err.Error())
		return err
	}
	err = db.Save(collection)
	if err != nil {
		return err
	}

	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	fmt.Printf("Your todo is now %s.\n", todo.Status)
	ct.ResetColor()
	return nil
}

func TdRemove(c *cli.Context) error {
	var err error

	if c.Args().Len() != 1 {
		WriteError("You must provide the position of the item you want to remove.",
			"Example: Td --remove 1")
		return argError
	}

	id, err := strconv.ParseInt(c.Args().Get(0), 10, 32)
	if err != nil {
		fmt.Println(err)
		return err
	}
	collection, err := db.Read()
	if err != nil {
		WriteError(err.Error())
		return err
	}
	err = collection.Remove(int(id))
	if err != nil {
		WriteError(err.Error())
		return err
	}
	db.Save(collection)

	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	fmt.Printf("Todo at %d is deleted .\n", id)
	ct.ResetColor()
	return nil
}

func TdReorder(c *cli.Context) error {
	collection, err := db.Read()
	if err != nil {
		WriteError(err.Error())
		return err
	}

	if c.Args().Len() != 1 {
		WriteError("You must provide two position if you want to swap todos.",
			"Example: Td reorder 9 3")
		return argError
	} else if c.Args().Len() != 2 {
		idA, err := strconv.ParseInt(c.Args().Get(0), 10, 32)
		if err != nil {
			WriteError(err.Error())
			return err
		}
		idB, err := strconv.ParseInt(c.Args().Get(1), 10, 32)
		if err != nil {
			WriteError(err.Error())
			return err
		}
		err = collection.Swap(idA, idB)
		if err != nil {
			WriteError(err.Error())
			return err
		}
		ct.ChangeColor(ct.Cyan, false, ct.None, false)
		fmt.Printf("\"%s\" and \"%s\" has been swapped\n", c.Args().Get(0), c.Args().Get(1))
		ct.ResetColor()
	}

	collection.Reorder()
	err = db.Save(collection)
	if err != nil {
		WriteError(err.Error())
		return err
	}

	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	fmt.Println("Your list is now reordered.")
	ct.ResetColor()
	return nil
}
func TdSearch(c *cli.Context) error {
	if c.Args().Len() != 1 {
		WriteError("You must provide a string earch.",
			"Example: Td search \"project-1\"")
		return argError
	}

	collection, err := db.Read()
	if err != nil {
		WriteError(err.Error())
		return err
	}
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

func TdSearchByDate(c *cli.Context) error {
	var date time.Time
	if c.Args().Len() < 2 {
		date = time.Now().Add(time.Hour * 24)
	} else {
		var err error
		date, err = parser.ParseDate(c.Args().Get(1))
		if err != nil {
			WriteError(err.Error())
			return argError
		}
	}
	collection, err := db.Read()
	if err != nil {
		WriteError(err.Error())
		return err
	}
	for _, v := range collection.Todos {
		if !v.Deadline.IsZero() && v.Deadline.Before(date) {
			v.MakeOutput(true)
		}
	}
	return nil
}
