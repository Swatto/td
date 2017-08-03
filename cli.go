package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/daviddengcn/go-colortext"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "td"
	app.Usage = "Your todos manager"
	app.Version = "1.4.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Gaël Gillard",
			Email: "gael@gaelgillard.com",
		},
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "done, d",
			Usage: "print done todos",
		},
		cli.BoolFlag{
			Name:  "wip, w",
			Usage: "print working in progress todos",
		},
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "print all todos",
		},
	}
	app.Action = func(c *cli.Context) error {

		collection, err := NewCollection()
		if err != nil {
			return exitError(err)
		}

		if !c.IsSet("all") {
			switch {
			case c.IsSet("done"):
				collection.ListDoneTodos()

			case c.IsSet("wip"):
				collection.ListWorkInProgressTodos()

			default:
				collection.ListUndoneTodos()
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
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:      "init",
			ShortName: "i",
			Usage:     "Initialize a collection of todos. If not path defined, it will create a file named .todos in the current directory.",
			UsageText: "td init",
			Action: func(c *cli.Context) error {

				db, err := NewDataStore()
				if err != nil {
					return exitError(err)
				}

				if err := db.Initialize(); err != nil {
					return exitError(err)
				}

				ct.ChangeColor(ct.Cyan, false, ct.None, false)
				fmt.Printf("Initialized empty to-do file as \"%s\".\n", db.Path)
				ct.ResetColor()
				return nil
			},
		},
		{
			Name:      "add",
			ShortName: "a",
			Usage:     "Add a new todo",
			UsageText: "td add \"call mum\"",
			Action: func(c *cli.Context) error {

				if len(c.Args()) != 1 {
					return exitError(
						fmt.Errorf("You must provide a name to your todo.\nUsage: %s", c.Command.UsageText))
				}

				collection, err := NewCollection()
				if err != nil {
					return exitError(err)
				}

				todo := NewTodo()
				todo.Desc = c.Args()[0]

				id, err := collection.CreateTodo(todo)
				if err != nil {
					return exitError(err)
				}

				if err := collection.WriteTodos(); err != nil {
					return exitError(err)
				}

				ct.ChangeColor(ct.Cyan, false, ct.None, false)
				fmt.Printf("#%d \"%s\" is now added to your todos.\n", id, c.Args()[0])
				ct.ResetColor()
				return nil
			},
		},
		{
			Name:      "modify",
			ShortName: "m",
			Usage:     "Modify the text of an existing todo",
			UsageText: "td modify 2 \"call dad\"",
			Action: func(c *cli.Context) error {

				if len(c.Args()) != 2 {
					return exitError(
						fmt.Errorf("You must provide the id and the new text for your todo.\nUsage: %s", c.Command.UsageText))
				}

				collection, err := NewCollection()
				if err != nil {
					return exitError(err)
				}

				id, err := strconv.ParseInt(c.Args()[0], 10, 32)
				if err != nil {
					return exitError(err)
				}

				_, err = collection.Modify(id, c.Args()[1])
				if err != nil {
					return exitError(err)
				}

				if err := collection.WriteTodos(); err != nil {
					return exitError(err)
				}

				ct.ChangeColor(ct.Cyan, false, ct.None, false)
				fmt.Printf("\"%s\" has now a new description: %s\n", c.Args()[0], c.Args()[1])
				ct.ResetColor()
				return nil
			},
		},
		{
			Name:      "toggle",
			ShortName: "t",
			Usage:     "Toggle the status of a todo by giving his id",
			UsageText: "td toggle 1",
			Action: func(c *cli.Context) error {

				if len(c.Args()) != 1 {
					return exitError(
						fmt.Errorf("You must provide the position of the item you want to change.\nUsage: %s", c.Command.UsageText))
				}

				collection, err := NewCollection()
				if err != nil {
					return exitError(err)
				}

				id, err := strconv.ParseInt(c.Args()[0], 10, 32)
				if err != nil {
					return exitError(err)
				}

				todo, err := collection.Toggle(id)
				if err != nil {
					return exitError(err)
				}

				var status string

				switch todo.Status {
				case "wip":
					status = "marked as work in progress"
				default:
					status = todo.Status
				}

				ct.ChangeColor(ct.Cyan, false, ct.None, false)
				fmt.Printf("Your todo is now %s.\n", status)
				ct.ResetColor()
				return nil
			},
		},
		{
			Name:      "wip",
			ShortName: "w",
			Usage:     "Change the status of a todo to \"Work In Progress\" by giving its id",
			UsageText: "td wip 1",
			Action: func(c *cli.Context) error {

				if len(c.Args()) != 1 {
					return exitError(
						fmt.Errorf("You must provide the position of the item you want to change.\nUsage: %s", c.Command.UsageText))
				}

				collection, err := NewCollection()
				if err != nil {
					return exitError(err)
				}

				id, err := strconv.ParseInt(c.Args()[0], 10, 32)
				if err != nil {
					return exitError(err)
				}

				todo, err := collection.SetStatus(id, WIP)
				if err != nil {
					return exitError(err)
				}

				if err := collection.WriteTodos(); err != nil {
					return exitError(err)
				}

				var status string
				switch todo.Status {
				case WIP:
					status = "marked as work in progress"
				default:
					status = todo.Status
				}

				ct.ChangeColor(ct.Cyan, false, ct.None, false)
				fmt.Printf("Your todo is now %s.\n", status)
				ct.ResetColor()
				return nil
			},
		},
		{
			Name:      "clean",
			ShortName: "c",
			Usage:     "Remove finished todos from the list",
			Action: func(c *cli.Context) error {

				collection, err := NewCollection()
				if err != nil {
					return exitError(err)
				}

				if err := collection.RemoveFinishedTodos(); err != nil {
					return exitError(err)
				}

				if err := collection.WriteTodos(); err != nil {
					return exitError(err)
				}

				ct.ChangeColor(ct.Cyan, false, ct.None, false)
				fmt.Println("Your list is now flushed of finished todos.")
				ct.ResetColor()
				return nil
			},
		},
		{
			Name:      "reorder",
			ShortName: "r",
			Usage:     "Reset ids of todo",
			Action: func(c *cli.Context) error {

				collection, err := NewCollection()
				if err != nil {
					return exitError(err)
				}

				if err := collection.Reorder(); err != nil {
					return exitError(err)
				}

				if err := collection.WriteTodos(); err != nil {
					return exitError(err)
				}

				ct.ChangeColor(ct.Cyan, false, ct.None, false)
				fmt.Println("Your list is now reordered.")
				ct.ResetColor()
				return nil
			},
		},
		{
			Name:      "swap",
			ShortName: "sw",
			Usage:     "Swap the position of two todos",
			UsageText: "td swap 9 3",
			Action: func(c *cli.Context) error {

				if len(c.Args()) != 1 {
					return exitError(
						fmt.Errorf("You must provide two position if you want to swap todos.\nUsage: %s", c.Command.UsageText))
				}

				collection, err := NewCollection()
				if err != nil {
					return exitError(err)
				}

				idA, err := strconv.ParseInt(c.Args()[0], 10, 32)
				if err != nil {
					return exitError(err)
				}

				idB, err := strconv.ParseInt(c.Args()[1], 10, 32)
				if err != nil {
					return exitError(err)
				}

				_, err = collection.Find(idA)
				if err != nil {
					return exitError(err)
				}

				_, err = collection.Find(idB)
				if err != nil {
					return exitError(err)
				}

				if err := collection.Swap(idA, idB); err != nil {
					return exitError(err)
				}

				if err := collection.Reorder(); err != nil {
					return exitError(err)
				}

				if err := collection.WriteTodos(); err != nil {
					return exitError(err)
				}

				ct.ChangeColor(ct.Cyan, false, ct.None, false)
				fmt.Printf("\"%s\" and \"%s\" has been swapped\n", c.Args()[0], c.Args()[1])
				ct.ResetColor()

				return nil
			},
		},
		{
			Name:      "search",
			ShortName: "s",
			Usage:     "Search a string in all todos",
			UsageText: "td search \"project-1\"",
			Action: func(c *cli.Context) error {

				if len(c.Args()) != 1 {
					return exitError(
						fmt.Errorf("You must provide a string search.\nUsage: %s", c.Command.UsageText))
				}

				collection, err := NewCollection()
				if err != nil {
					return exitError(err)
				}

				collection.Search(c.Args()[0])

				if len(collection.Todos) == 0 {
					ct.ChangeColor(ct.Cyan, false, ct.None, false)
					fmt.Printf("Sorry, there's no todos containing \"%s\".\n", c.Args()[0])
					ct.ResetColor()
					return nil
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
			},
		},
	}

	app.After = func(c *cli.Context) error {
		db, _ := NewDataStore()
		ct.ChangeColor(ct.Magenta, false, ct.None, false)
		// fmt.Println(fmt.Sprintf("%-25s %s", `to-do file:`, db.Path))
		fmt.Println(db.Path)
		ct.ResetColor()
		return nil
	}

	app.Before = func(c *cli.Context) error {

		if len(c.Args()) == 1 {
			exceptions := []string{"init", "i", "help", "h"}
			for _, x := range exceptions {
				if c.Args()[0] == x {
					return nil
				}
			}
		}

		db, err := NewDataStore()
		if err != nil {
			return err
		}

		if err := db.Check(); err != nil {
			errDS := fmt.Errorf(`
===============================================================================

ERROR:

  File to store your todos could be found. Your current file location is:
  %s

  Run 'td init' to start a new to-do list or set/update the environment
  variable named 'TODO_DB_PATH' with the correct location of your to-dos file.

  Example 'export TODO_DB_PATH=$HOME/Dropbox/todo.json'

  If 'TODO_DB_PATH' is blank, it will reference to a file named '.todos' in the
  current working folder, and if there's no file, it will create one.

===============================================================================

				 `, db.Path)

			return cli.NewExitError(errDS, 1)
		}

		return nil
	}

	app.Run(os.Args)
}

func exitError(message error) *cli.ExitError {
	ct.ChangeColor(ct.Red, false, ct.None, false)
	return cli.NewExitError(message, 1)
}
