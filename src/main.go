package main

import (
	cli "github.com/urfave/cli/v2"
	"os"
	"umutsevdi/td/cmd"
	"umutsevdi/td/db"
)

func main() {
	app := cli.NewApp()
	app.Name = "td"
	app.Usage = "Your todos manager"
	app.Version = "1.4.2~fork"
	app.Authors = []*cli.Author{
		{
			Name:  "Gaël Gillard",
			Email: "gillardgael@gmail.com",
		},
		{
			Name:  "Umut Sevdi",
			Email: "sevdiumut@gmail.com",
		},
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "past, p",
			Usage: "print todos that are past due",
		},
		&cli.BoolFlag{
			Name:  "done, d",
			Usage: "print done todos",
		},
		&cli.BoolFlag{
			Name:  "all, a",
			Usage: "print all todos",
		},
	}
	app.Action = cmd.TdList
	app.Commands = []*cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Initialize a collection of todos",
			Action:  cmd.TdInit,
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add a new todo",
			Action:  cmd.TdAdd,
		},
		{
			Name:    "modify",
			Aliases: []string{"m"},
			Usage:   "Modify the text or any property of an existing todo",
			Action:  cmd.TdModify,
		},
		{
			Name:    "toggle",
			Aliases: []string{"t"},
			Usage:   "Toggle the status of a todo by giving his id",
			Action:  cmd.TdToggle,
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "Remove an existing todo",
			Action:  cmd.TdRemove,
		},
		//{
		//	Name:    "clean",
		//	Aliases: []string{"c"},
		//	Usage:   "Remove finished todos from the list",
		//	Action:  cmd.TdClean,
		//},
		{
			Name:    "reorder",
			Aliases: []string{"r"},
			Usage:   "Reset ids of todo (no arguments) or swap the position of two todos",
			Action:  cmd.TdReorder,
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search a string in all todos",
			Action:  cmd.TdSearch,
		},
		{
			Name:    "filter",
			Aliases: []string{"f"},
			Usage:   "Search for todos that have upcoming deadlines",
			Action:  cmd.TdSearchByDate,
		},
	}

	app.Before = func(c *cli.Context) error {
		var err error
		path := db.GetDBPath()

		if path == "" {
			cmd.WriteError("A store for your todos is missing. You have 2 possibilities:",
				"  1. create a \".todos\" file in your local folder.",
				"  2. the environment variable \"TODO_DB_PATH\" could be set.",
				"    (example: \"export TODO_DB_PATH=$HOME/Dropbox/todo.json\" in your .bashrc or .bash_profile)")
		}
		db.CreateStoreFileIfNeeded(path)
		return err
	}
	app.Run(os.Args)
}