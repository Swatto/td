# td

> Your todo list in your terminal.
>
> ![Screenshot](screenshot.png)

## Usage

### Installation

- From *homebrew*: `brew install td`
- From *binary*: go to the [release page](https://github.com/Swatto/td/releases)
- From *source*: `go get github.com/Swatto/td`

### Information

*td* will look at a `.todos` files to store your todos (like Git does: it will try recursively in each parent folder). This permit to have different list of todos per folder.

If it doesn't find a `.todos`, *td* use an environment variable to store your todos: `TODO_DB_PATH` where you define the path to the JSON file. If the file doesn't exist, the program will create it for you.

### CLI

```
NAME:
   td - Your todos manager

USAGE:
   td [global options] command [command options] [arguments...]

VERSION:
   1.4.1

AUTHOR:
  Gaël Gillard - <gael@gaelgillard.com>

COMMANDS:
   init, i  Initialize a collection of todos
   add, a   Add a new todo
   modify, m   Modify the text of an existing todo
   toggle, t   Toggle the status of a todo by giving his id
   clean Remove finished todos from the list
   reorder, r  Reset ids of todo or swap the position of two todo
   search, s   Search a string in all todos
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --done, -d     print done todos
   --all, -a      print all todos
   --help, -h     show help
   --version, -v  print the version
```
