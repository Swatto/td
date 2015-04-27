package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/daviddengcn/go-colortext"
)

func main() {
	app := cli.NewApp()
	app.Name = "td"
	app.Usage = "Your todos manager"
	app.Version = "1.1.1"
	app.Author = "Gaël Gillard"
	app.Email = ""
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "done, d",
			Usage: "print done todos",
		},
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "print all todos",
		},
	}
	app.Action = func(c *cli.Context) {
		var err error
		collection := Collection{}

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
					todo.MakeOutput()
				}
				fmt.Println()
			} else {
				ct.ChangeColor(ct.Cyan, false, ct.None, false)
				fmt.Println("There's no todo to show.")
				ct.ResetColor()
			}
		}
	}
	app.Commands = []cli.Command{
		{
			Name:      "add",
			ShortName: "a",
			Usage:     "Add a new todo",
			Action: func(c *cli.Context) {

				if len(c.Args()) != 1 {
					fmt.Println()
					ct.ChangeColor(ct.Red, false, ct.None, false)
					fmt.Println("Error")
					ct.ResetColor()
					fmt.Println("You must provide a name to your todo.")
					fmt.Println("Example: td add \"call mum\"")
					fmt.Println()
					return
				}

				collection := Collection{}
				todo := Todo{
					Id:       0,
					Desc:     c.Args()[0],
					Status:   "pending",
					Modified: "",
				}
				err := collection.RetrieveTodos()
				if err != nil {
					fmt.Println(err)
				}

				id, err := collection.CreateTodo(&todo)
				if err != nil {
					fmt.Println(err)
				}

				ct.ChangeColor(ct.Cyan, false, ct.None, false)
				fmt.Printf("#%d \"%s\" is now added to your todos.\n", id, c.Args()[0])
				ct.ResetColor()
			},
		},
		{
			Name:      "modify",
			ShortName: "m",
			Usage:     "Modify the text of an existing todo",
			Action: func(c *cli.Context) {

				if len(c.Args()) != 2 {
					fmt.Println()
					ct.ChangeColor(ct.Red, false, ct.None, false)
					fmt.Println("Error")
					ct.ResetColor()
					fmt.Println("You must provide the id and the new text for your todo.")
					fmt.Println("Example: td modify 2 \"call dad\"")
					fmt.Println()
					return
				}

				collection := Collection{}
				collection.RetrieveTodos()

				id, err := strconv.ParseInt(c.Args()[0], 10, 32)
				if err != nil {
					fmt.Println(err)
					return
				}

				_, err = collection.Modify(id, c.Args()[1])
				if err != nil {
					fmt.Println(err)
					return
				}

				ct.ChangeColor(ct.Cyan, false, ct.None, false)
				fmt.Printf("\"%s\" has now a new description: %s", c.Args()[0], c.Args()[1])
				ct.ResetColor()
			},
		},
		{
			Name:      "toggle",
			ShortName: "t",
			Usage:     "Toggle the status of a todo by giving his id",
			Action: func(c *cli.Context) {
				var err error

				if len(c.Args()) != 1 {
					fmt.Println()
					ct.ChangeColor(ct.Cyan, false, ct.None, false)
					fmt.Println("Error")
					ct.ResetColor()
					fmt.Println("You must the id of the item you want to change todo.")
					fmt.Println("Example: td toggle 1")
					fmt.Println()
					return
				}

				collection := Collection{}
				collection.RetrieveTodos()

				id, err := strconv.ParseInt(c.Args()[0], 10, 32)
				if err != nil {
					fmt.Println(err)
					return
				}

				todo, err := collection.Toggle(id)
				if err != nil {
					fmt.Println(err)
					return
				}

				ct.ChangeColor(ct.Cyan, false, ct.None, false)
				fmt.Printf("Your todo is now %s.\n", todo.Status)
				ct.ResetColor()
			},
		},
		{
			Name:  "clean",
			Usage: "Remove finished todos from the list",
			Action: func(c *cli.Context) {
				collection := Collection{}
				collection.RetrieveTodos()

				err := collection.RemoveFinishedTodos()

				if err != nil {
					fmt.Println(err)
					return
				} else {
					ct.ChangeColor(ct.Cyan, false, ct.None, false)
					fmt.Println("Your list is now flushed of finished todos.")
					ct.ResetColor()
				}
			},
		},
		{
			Name:      "reorder",
			ShortName: "r",
			Usage:     "Reset ids of todo",
			Action: func(c *cli.Context) {
				collection := Collection{}
				collection.RetrieveTodos()

				err := collection.Reorder()

				if err != nil {
					fmt.Println(err)
					return
				}

				ct.ChangeColor(ct.Cyan, false, ct.None, false)
				fmt.Println("Your list is now reordered.")
				ct.ResetColor()
			},
		},
		{
			Name:      "search",
			ShortName: "s",
			Usage:     "Search a string in all todos",
			Action: func(c *cli.Context) {
				if len(c.Args()) != 1 {
					fmt.Println()
					ct.ChangeColor(ct.Red, false, ct.None, false)
					fmt.Println("Error")
					ct.ResetColor()
					fmt.Println("You must provide a string earch.")
					fmt.Println("Example: td search \"project-1\"")
					fmt.Println()
					return
				}

				collection := Collection{}
				collection.RetrieveTodos()
				collection.Search(c.Args()[0])

				if len(collection.Todos) == 0 {
					ct.ChangeColor(ct.Cyan, false, ct.None, false)
					fmt.Printf("Sorry, there's no todos containing \"%s\".\n", c.Args()[0])
					ct.ResetColor()
					return
				}

				if len(collection.Todos) > 0 {
					fmt.Println()
					for _, todo := range collection.Todos {
						todo.MakeOutput()
					}
					fmt.Println()
				} else {
					ct.ChangeColor(ct.Cyan, false, ct.None, false)
					fmt.Println("There's no todo to show.")
					ct.ResetColor()
				}
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		var err error
		path := GetDBPath()

		if path == "" {
			fmt.Println()
			ct.ChangeColor(ct.Red, false, ct.None, false)
			fmt.Println("Error")
			ct.ResetColor()
			fmt.Println("The environment variable \"TODO_DB_PATH\" must be set.")
			fmt.Println("Example: \"export TODO_DB_PATH=$HOME/Dropbox/todo.json\" in your .bash_profile")
			fmt.Println()
		}

		CreateStoreFileIfNeeded(path)

		return err
	}

	app.Run(os.Args)
}
