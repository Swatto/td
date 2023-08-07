/*
Contains functions that are called from the CLI
@author umutsevdi - 2023-02-25
*/
package cmd

import (
	"fmt"
	"github.com/swatto/td/db"
	"github.com/swatto/td/parser"
	"os"
	"strconv"
	"strings"

	ct "github.com/daviddengcn/go-colortext"
	cli "github.com/urfave/cli/v2"
)

type Error string
type MessageType int

const (
	argError  = Error("Error in argument")
	MT_INFO   = MessageType(0)
	MT_UPDATE = MessageType(1)
	MT_ERROR  = MessageType(2)
)

func (e Error) Error() string { return string(e) }

// Prints the generic error message to the screen
//
// All given strings are joined with " " by default. Replacement character can
// be changed by given as the second parameter. All following strings are joined using
// given character. Possible replacement characters are:
//   - ""
//   - "\n"
//   - "\t"
//   - ", "
//   - ","
func Write(c *cli.Context, m MessageType, v ...string) {
	join := " "
	if v[0] == "" || v[0] == "\n" || v[0] == "\t" || v[0] == ", " || v[0] == "," {
		join = v[0]
	}
	if c.Bool("color") {
		switch m {
		case MT_UPDATE:
			ct.ChangeColor(ct.Blue, false, ct.None, false)
		case MT_INFO:
			ct.ChangeColor(ct.Cyan, false, ct.None, false)
		case MT_ERROR:
			ct.ChangeColor(ct.Red, false, ct.None, false)
			fmt.Print("Error: ")
			ct.ResetColor()
		}
	}
	fmt.Println(strings.Join(v, join))
	if c.Bool("color") {
		ct.ResetColor()
	}
}

// Lists all todo items
func TdList(c *cli.Context) error {
	collection, err := db.Read()
	if err != nil {
		Write(c, MT_ERROR, err.Error())
	} else {
		var status db.STATUS = db.STATUS_ANY
		if c.IsSet("done") {
			status = db.STATUS_DONE
		} else if c.IsSet("expired") {
			status = db.STATUS_EXPIRED
		} else if c.IsSet("pending") {
			status = db.STATUS_PENDING
		}
		collection.List(status)
		if !c.IsSet("all") && c.Int("before") > 0 {
			collection.FilterRecent(c.Int("before"), status)
		}

		if len(collection.Todos) > 0 {
			for _, todo := range collection.Todos {
				todo.MakeOutput(c.Bool("color"), c.Bool("nerd"))
			}
		} else {
			Write(c, MT_INFO, "There is no todo to show.")
		}
	}
	return nil
}

// Initialize a collection of todos
func TdInit(c *cli.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("%s .\n", err)
		return err
	}

	err = db.CreateStoreFileIfNeeded(cwd + "/.todos")
	ct.ChangeColor(ct.Cyan, false, ct.None, false)
	if err != nil {
		Write(c, MT_INFO, "A \".todos\" file already exist in ", cwd)
	} else {
		Write(c, MT_UPDATE, "A \".todos\" file is now added to ", cwd)
	}
	ct.ResetColor()
	return nil
}

// Add a new todo
func TdAdd(c *cli.Context) error {
	if c.Args().Len() == 0 {
		Write(c, MT_ERROR, "\n", "You must provide a name to your todo.",
			"Example: Td add \"call mum\"")
		return argError
	}
	dt, dtErr := parser.ParseDate(c.Args().Get(1))
	if dtErr != nil {
		Write(c, MT_ERROR, "\n", "Invalid date time format",
			"Available Formats: \"[dd/MM/yyyy] [hh:mm]\"")
		return argError
	}
	p, _ := strconv.Atoi(c.Args().Get(2))

	collection, err := db.Read()
	if err != nil {
		Write(c, MT_ERROR, err.Error())
		return err
	}

	todo := collection.Add(c.Args().Get(0), dt, p)

	if !dt.IsZero() {
		Write(c, MT_INFO, "Deadline added")
	}
	if p != 0 {
		Write(c, MT_INFO, "Periodic added")
	}
	err = db.Save(collection)
	if err != nil {
		fmt.Println(err)
		return err
	}
	Write(c, MT_UPDATE, "", "\"", todo.Desc, "\" is now added to your todos at ", strconv.Itoa(todo.ID), ".")
	return nil
}

// Modify the text or any property of an existing todo
//
// If there is only one argument, it will be interpreted as [todo.Todo.Desc]
// If there are more than one arguments they will be mapped based on following options:
//   - \-d | desc   : Description
//   - \-D | date   : A string that is valid according to [parser.ParseDate]
//   - \-p | period : A string that contains the period
func TdModify(c *cli.Context) error {
	if c.Args().Len() < 2 {
		Write(c, MT_ERROR, "\n", "You must provide the id and the new text for your todo.",
			"Example: Td modify 2 \"call dad\"")
		return argError
	}
	args := c.Args()

	id, err := strconv.Atoi(args.Get(0))
	if err != nil {
		Write(c, MT_ERROR, err.Error())
		return err
	}
	collection, err := db.Read()
	if err != nil {
		Write(c, MT_ERROR, err.Error())
		return err
	}
	m, err := parser.MapArgs(args)
	if err != nil {
		Write(c, MT_ERROR, err.Error())
		return err
	}
	Td, err := collection.Modify(id, m)
	if err != nil {
		Write(c, MT_ERROR, err.Error())
		return err
	}
	err = db.Save(collection)
	if err != nil {
		fmt.Println(err)
		return err
	}

	Write(c, MT_UPDATE, args.Get(0), " has been updated:")
	Td.MakeOutput(c.Bool("color"), c.Bool("nerd"))
	return nil
}

// Toggle the status of a todo by giving his id
func TdToggle(c *cli.Context) error {
	if c.Args().Len() != 1 {
		Write(c, MT_ERROR, "\n", "You must provide the position of the item you want to change.",
			"Example: Td toggle 1")
		return argError
	}

	collection, err := db.Read()
	if err != nil {
		Write(c, MT_ERROR, err.Error())
		return err
	}

	id, err := strconv.Atoi(c.Args().Get(0))
	if err != nil {
		Write(c, MT_ERROR, err.Error())
		return err
	}

	todo, err := collection.Toggle(id)
	if err != nil {
		Write(c, MT_ERROR, err.Error())
		return err
	}
	err = db.Save(collection)
	if err != nil {
		return err
	}

	Write(c, MT_UPDATE, "Your todo is now ", todo.Status)
	return nil
}

// Remove an existing todo
func TdRemove(c *cli.Context) error {
	var err error

	if c.Args().Len() != 1 {
		Write(c, MT_ERROR, "\n", "You must provide the position of the item you want to remove.",
			"Example: Td remove 1")
		return argError
	}

	id, err := strconv.Atoi(c.Args().Get(0))
	if err != nil {
		fmt.Println(err)
		return err
	}
	collection, err := db.Read()
	if err != nil {
		Write(c, MT_ERROR, err.Error())
		return err
	}
	err = collection.Remove(int(id))
	if err != nil {
		Write(c, MT_ERROR, err.Error())
		return err
	}
	db.Save(collection)

	Write(c, MT_UPDATE, "Todo at", strconv.Itoa(id), "is deleted.")
	return nil
}

// Reset ids of todo
func TdReorder(c *cli.Context) error {
	collection, err := db.Read()
	if err != nil {
		Write(c, MT_ERROR, err.Error())
		return err
	}
	err = db.Save(collection.Reorder())
	if err != nil {
		Write(c, MT_ERROR, err.Error())
		return err
	}

	Write(c, MT_UPDATE, "Your list is now reordered.")
	return nil
}

// Swap the position of two todos
func TdSwap(c *cli.Context) error {
	collection, err := db.Read()
	if err != nil {
		Write(c, MT_ERROR, err.Error())
		return err
	}

	if c.Args().Len() != 2 {
		Write(c, MT_ERROR, "\n", "You must provide two position if you want to swap todos.",
			"Example: Td swap 9 3")
		return argError
	}
	idA, err := strconv.Atoi(c.Args().Get(0))
	if err != nil {
		Write(c, MT_ERROR, err.Error())
		return err
	}
	idB, err := strconv.Atoi(c.Args().Get(1))
	if err != nil {
		Write(c, MT_ERROR, err.Error())
		return err
	}
	err = collection.Swap(idA, idB)
	if err != nil {
		Write(c, MT_ERROR, err.Error())
		return err
	}
	ct.ChangeColor(ct.Cyan, false, ct.None, false)

	Write(c, MT_UPDATE, c.Args().Get(0), "and", c.Args().Get(1), "has been swapped.")
	ct.ResetColor()

	err = db.Save(collection)
	if err != nil {
		Write(c, MT_ERROR, err.Error())
		return err
	}
	return nil
}

// Search a string in all todos
func TdSearch(c *cli.Context) error {
	if c.Args().Len() != 1 {
		Write(c, MT_ERROR, "\n", "You must provide a string earch.",
			"Example: Td search \"project-1\"")
		return argError
	}

	collection, err := db.Read()
	if err != nil {
		Write(c, MT_ERROR, err.Error())
		return err
	}
	collection.Search(c.Args().Get(0))

	if len(collection.Todos) == 0 {
		Write(c, MT_INFO, "Sorry, there's no todos containing", c.Args().Get(0))
		return argError
	}

	if len(collection.Todos) > 0 {
		for _, todo := range collection.Todos {
			todo.MakeOutput(c.Bool("color"), c.Bool("nerd"))
		}
	} else {
		Write(c, MT_INFO, "There's no todo to show.")
	}
	return nil
}
