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
   1.5

AUTHORS:
   Gaël Gillard <gillardgael@gmail.com>
   Umut Sevdi <sevdiumut@gmail.com>

COMMANDS:
   init, i     Initialize a collection of todos
   add, a      Add a new todo
   modify, m   Modify the text or any property of an existing todo
   toggle, t   Toggle the status of a todo by giving his id
   delete, d   Remove an existing todo
   reorder, r  Reset ids of todo
   swap, s     swap the position of two todos
   search, S   Search a string in all todos
   filter, f   Search for todos that have upcoming deadlines
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --past         print todos that are past due (default: false)
   --done         print done todos (default: false)
   --all          print all todos (default: false)
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```
