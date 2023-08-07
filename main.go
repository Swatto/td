package main

import (
	"github.com/swatto/td/cmd"
	"github.com/swatto/td/db"
	cli "github.com/urfave/cli/v2"
	"os"
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
			Name:  "expired",
			Usage: "list tasks that haven't been completed in time",
		},
		&cli.BoolFlag{
			Name:  "done",
			Usage: "list only completed tasks",
		},
		&cli.BoolFlag{
			Name:  "pending",
			Usage: "list only pending tasks",
		},
		&cli.BoolFlag{
			Name:  "all",
			Usage: "list all tasks",
		},
		&cli.IntFlag{
			Name:  "before",
			Usage: "print tasks within the amount of days",
			Value: 14,
		},
		&cli.BoolFlag{
			Name:  "color",
			Usage: "print with color",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  "nerd",
			Usage: "enable nerd fonts, disabled at Windows Builds",
			Value: false,
		},
	}
	app.Action = cmd.TdList
	app.Commands = []*cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Initialize a collection of tasks",
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
		{
			Name:    "reorder",
			Aliases: []string{"r"},
			Usage:   "Reset ids of todo",
			Action:  cmd.TdReorder,
		},
		{
			Name:    "swap",
			Aliases: []string{"s"},
			Usage:   "swap the position of two tasks",
			Action:  cmd.TdSwap,
		},
		{
			Name:    "search",
			Aliases: []string{"S"},
			Usage:   "Search a string in all tasks",
			Action:  cmd.TdSearch,
		},
	}

	app.Before = func(c *cli.Context) error {
		var err error
		path := db.GetDBPath()

		if path == "" {
			cmd.Write(c, cmd.MT_ERROR, "\n", "A store for your todos is missing. You have 2 possibilities:",
				"  1. create a \".todos\" file in your local folder.",
				"  2. the environment variable \"TODO_DB_PATH\" could be set.",
				"    (example: \"export TODO_DB_PATH=$HOME/Dropbox/todo.json\" in your .bashrc or .bash_profile)")
		}
		db.CreateStoreFileIfNeeded(path)
		return err
	}
	app.Run(os.Args)
}
