package main

import (
  "fmt"
  "os"
  "time"
  "strconv"
  "errors"
  "github.com/codegangsta/cli"
  "github.com/fatih/color"
)

func main() {
  app := cli.NewApp()
  app.Name = "td"
  app.Usage = "Your todos manager"
  app.Version = "0.1.0"
  app.Author = "Gaël Gillard"
  app.Email = "not_here@gmy_privacy.com"
  app.Flags = []cli.Flag {
    cli.BoolFlag{
      Name: "done, d",
      Usage: "print done todos",
    },
    cli.BoolFlag{
      Name: "all, a",
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
        if c.IsSet("done"){
          collection.ListDoneTodos()
        } else {
          collection.ListPendingTodos()
        }
      }

      fmt.Println()
      for _, todo := range collection.Todos {
        todo.MakeOutput()
      }     
           
      fmt.Println()
    }
  }

  app.Commands = []cli.Command{
    {
      Name: "create",
      ShortName: "add",
      Usage: "Create a new todo",
      Action: func(c *cli.Context) {

        if len(c.Args()) != 1 {
          fmt.Println()
          color.Red("Error")
          fmt.Println("You must provide a name to your todo.")
          fmt.Println("Example: todo create \"call mum\"")
          fmt.Println()
          return
        }

        collection := Collection{}
        todo := Todo{
          Id: 0,
          Desc: c.Args()[0],
          Status: "pending",
          Modified: (time.Now().Local()).Format(time.Stamp),
        }
        err := collection.RetrieveTodos()
        if err != nil {
          fmt.Println(err)
        }

        err = collection.CreateTodo(&todo)
        if err != nil {
          fmt.Println(err)
        }

        color.Cyan("\"%s\" is now added to your todos.", c.Args()[0])
      },
    },
    {
      Name: "toggle",
      ShortName: "t",
      Usage: "Toggle the status of a todo by giving his id",
      Action: func(c *cli.Context) {
        collection := Collection{}
        collection.RetrieveTodos()

        id, err := strconv.ParseInt(c.Args()[0], 10, 32)
        if err != nil {
          fmt.Println(err)
          return
        }

        todo, err_update := collection.Toggle(id)
        if err_update != nil {
          fmt.Println(err)
          return
        }
        todo.MakeOutput()
      },
    },
    {
      Name: "clean",
      Usage: "Remove finished todos from the list",
      Action: func(c *cli.Context) {
        collection := Collection{}
        collection.RetrieveTodos()

        err := collection.RemoveFinishedTodos()

        if err != nil {
          fmt.Println(err)
          return
        } else {
          color.Cyan("Your list is now flushed of finished todos.")
        }
      },
    },
  }

  app.Before = func(c *cli.Context) error {
    var err error

    path := os.Getenv("TODO_DB_PATH")
    if path == "" {
      err = errors.New("The environment variable \"TODO_DB_PATH\" need to be set.")
      err = errors.New("")
      fmt.Println()
      color.Red("Error")
      fmt.Println("The environment variable \"TODO_DB_PATH\" need to be set.")
      fmt.Println("Example: \"export TODO_DB_PATH=$HOME/Dropbox/todo.json\" in your .bash_profile")
      fmt.Println()
    }

    return err
  }

  app.Run(os.Args)
}