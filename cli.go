package main

import (
	"fmt"
	"os"
	"strconv"

	ct "github.com/daviddengcn/go-colortext"
	cli "github.com/urfave/cli/v2"
)

type Error string

func (e Error) Error() string { return string(e) }

const argError = Error("Error in argument")

func main() {
	app := cli.NewApp()
	app.Name = "td"
	app.Usage = "Your todos manager"
	app.Version = "1.4.2"
	app.Authors = []*cli.Author{
		{
			Name:  "Gaël Gillard",
			Email: "gillardgael@gmail.com",
		},
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "done, d",
			Usage: "print done todos",
		},
		&cli.BoolFlag{
			Name:  "all, a",
			Usage: "print all todos",
		},
	}
	app.Action = func(c *cli.Context) error {
		var err error
		collection := collection{}

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
	app.Commands = []*cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Initialize a collection of todos",
			Action: func(c *cli.Context) error {
				cwd, err := os.Getwd()
				if err != nil {
					fmt.Printf("%s .\n", err)
					return err
				}

				err = createStoreFileIfNeeded(cwd + "/.todos")
				ct.ChangeColor(ct.Cyan, false, ct.None, false)
				if err != nil {
					fmt.Printf("A \".todos\" file already exist in \"%s\".\n", cwd)
				} else {
					fmt.Printf("A \".todos\" file is now added to \"%s\".\n", cwd)
				}
				ct.ResetColor()
				return nil
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add a new todo",
			Action: func(c *cli.Context) error {

				if c.Args().Len() != 1 {
					fmt.Println()
					ct.ChangeColor(ct.Red, false, ct.None, false)
					fmt.Println("Error")
					ct.ResetColor()
					fmt.Println("You must provide a name to your todo.")
					fmt.Println("Example: td add \"call mum\"")
					fmt.Println()
					return argError
				}

				collection := collection{}
				todo := todo{
					ID:       0,
					Desc:     c.Args().Get(0),
					Status:   "pending",
					Modified: "",
				}
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
			},
		},
		{
			Name:    "modify",
			Aliases: []string{"m"},
			Usage:   "Modify the text of an existing todo",
			Action: func(c *cli.Context) error {

				if c.Args().Len() != 2 {
					fmt.Println()
					ct.ChangeColor(ct.Red, false, ct.None, false)
					fmt.Println("Error")
					ct.ResetColor()
					fmt.Println("You must provide the id and the new text for your todo.")
					fmt.Println("Example: td modify 2 \"call dad\"")
					fmt.Println()
					return argError
				}

				collection := collection{}
				collection.RetrieveTodos()

				args := c.Args()

				id, err := strconv.ParseInt(args.Get(0), 10, 32)
				if err != nil {
					fmt.Println(err)
					return err
				}

				_, err = collection.Modify(id, args.Get(1))
				if err != nil {
					fmt.Println(err)
					return err
				}

				ct.ChangeColor(ct.Cyan, false, ct.None, false)
				fmt.Printf("\"%s\" has now a new description: %s\n", args.Get(0), args.Get(1))
				ct.ResetColor()
				return nil
			},
		},
		{
			Name:    "toggle",
			Aliases: []string{"t"},
			Usage:   "Toggle the status of a todo by giving his id",
			Action: func(c *cli.Context) error {
				var err error

				if c.Args().Len() != 1 {
					fmt.Println()
					ct.ChangeColor(ct.Red, false, ct.None, false)
					fmt.Println("Error")
					ct.ResetColor()
					fmt.Println("You must provide the position of the item you want to change.")
					fmt.Println("Example: td toggle 1")
					fmt.Println()
					return argError
				}

				collection := collection{}
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
			},
		},
		{
			Name:    "clean",
			Aliases: []string{"c"},
			Usage:   "Remove finished todos from the list",
			Action: func(c *cli.Context) error {
				collection := collection{}
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
			},
		},
		{
			Name:    "reorder",
			Aliases: []string{"r"},
			Usage:   "Reset ids of todo (no arguments) or swap the position of two todos",
			Action: func(c *cli.Context) error {
				collection := collection{}
				collection.RetrieveTodos()

				if c.Args().Len() != 1 {
					fmt.Println()
					ct.ChangeColor(ct.Red, false, ct.None, false)
					fmt.Println("Error")
					ct.ResetColor()
					fmt.Println("You must provide two position if you want to swap todos.")
					fmt.Println("Example: td reorder 9 3")
					fmt.Println()
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
			},
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search a string in all todos",
			Action: func(c *cli.Context) error {
				if c.Args().Len() != 1 {
					fmt.Println()
					ct.ChangeColor(ct.Red, false, ct.None, false)
					fmt.Println("Error")
					ct.ResetColor()
					fmt.Println("You must provide a string earch.")
					fmt.Println("Example: td search \"project-1\"")
					fmt.Println()
					return argError
				}

				collection := collection{}
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
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		var err error
		path := getDBPath()

		if path == "" {
			fmt.Println()
			ct.ChangeColor(ct.Red, false, ct.None, false)
			fmt.Println("Error")
			fmt.Println("-----")
			ct.ResetColor()
			fmt.Println("A store for your todos is missing. You have 2 possibilities:")
			fmt.Println("  1. create a \".todos\" file in your local folder.")
			fmt.Println("  2. the environment variable \"TODO_DB_PATH\" could be set.")
			fmt.Println("    (example: \"export TODO_DB_PATH=$HOME/Dropbox/todo.json\" in your .bashrc or .bash_profile)")
			fmt.Println()
		}

		createStoreFileIfNeeded(path)

		return err
	}

	app.Run(os.Args)
}
