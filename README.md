# td

> Your todo list in your terminal.
> ![Screenshot](screenshot.png)

## Usage

### Installation

- From *binary*: go to the [release page](https://github.com/Swatto/td/releases)
- From *source*: `go get github.com/Swatto/td`

### Information

*td* use a environment variable to store your todos `TODO_DB_PATH` where you
define the path to the futur JSON file. If the file doesn't exist, the program
will create it for you.

But if your current directory has a `.todos` file, it will overwrite the environment variable.

### CLI

```
NAME:
   td - Your todos manager

USAGE:
   td [global options] command [command options] [arguments...]

VERSION:
   1.1.1

AUTHOR:
  Gaël Gillard

COMMANDS:
   add, a       Add a new todo
   modify, m    Modify the text of an existing todo
   toggle, t    Toggle the status of a todo by giving his id
   clean        Remove finished todos from the list
   reorder, r   Reset ids of todo
   search, s    Search a string in all todos
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --done, -d           print done todos
   --all, -a            print all todos
   --help, -h           show help
   --version, -v        print the version
```

## License

The MIT License (MIT)

Copyright (c) 2015 Gaël Gillard

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
