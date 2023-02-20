package main

import (
	cli "github.com/urfave/cli/v2"
	"os"
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
			Name:  "done, d",
			Usage: "print done todos",
		},
		&cli.BoolFlag{
			Name:  "all, a",
			Usage: "print all todos",
		},
	}
	app.Action = tdList
	app.Commands = []*cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Initialize a collection of todos",
			Action:  tdInit,
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add a new todo",
			Action:  tdAdd,
		},
		{
			Name:    "modify",
			Aliases: []string{"m"},
			Usage:   "Modify the text of an existing todo",
			Action:  tdModify,
		},
		{
			Name:    "toggle",
			Aliases: []string{"t"},
			Usage:   "Toggle the status of a todo by giving his id",
			Action:  tdToggle,
		},
		{
			Name:    "clean",
			Aliases: []string{"c"},
			Usage:   "Remove finished todos from the list",
			Action:  tdClean,
		},
		{
			Name:    "reorder",
			Aliases: []string{"r"},
			Usage:   "Reset ids of todo (no arguments) or swap the position of two todos",
			Action:  tdReorder,
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search a string in all todos",
			Action:  tdSearch,
		},
		{
			Name:    "filter",
			Aliases: []string{"s"},
			Usage:   "Search for todos that have upcoming deadlines",
			Action:  tdSearchByDate,
		},
	}

	app.Before = func(c *cli.Context) error {
		var err error
		path := db.GetDBPath()

		if path == "" {
			WriteError("A store for your todos is missing. You have 2 possibilities:",
				"  1. create a \".todos\" file in your local folder.",
				"  2. the environment variable \"TODO_DB_PATH\" could be set.",
				"    (example: \"export TODO_DB_PATH=$HOME/Dropbox/todo.json\" in your .bashrc or .bash_profile)")
		}
		db.CreateStoreFileIfNeeded(path)
		return err
	}
	app.Run(os.Args)
}
