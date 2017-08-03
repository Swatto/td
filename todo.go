package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/daviddengcn/go-colortext"
	p "github.com/vhugo/td/printer"
)

type Todo struct {
	Id       int64  `json:"id"`
	Desc     string `json:"desc"`
	Status   string `json:"status"`
	Modified string `json:"modified"`
}

func NewTodo() *Todo {
	var todo = new(Todo)
	todo.Status = PENDING
	return todo
}

func (t *Todo) MakeOutput(useColor bool) {
	var symbole string
	var color ct.Color

	switch t.Status {
	case "done":
		color = ct.Green
		symbole = p.OkSign
	case "wip":
		color = ct.Blue
		symbole = p.WpSign
	default:
		color = ct.Red
		symbole = p.KoSign
	}

	hashtagReg := regexp.MustCompile("#[^\\s]*")

	spaceCount := 6 - len(strconv.FormatInt(t.Id, 10))

	fmt.Print(strings.Repeat(" ", spaceCount), t.Id, " | ")
	if useColor == true {
		ct.ChangeColor(color, false, ct.None, false)
	}
	fmt.Print(symbole)
	if useColor == true {
		ct.ResetColor()
	}
	fmt.Print(" ")
	pos := 0
	for _, token := range hashtagReg.FindAllStringIndex(t.Desc, -1) {
		fmt.Print(t.Desc[pos:token[0]])
		if useColor == true {
			ct.ChangeColor(ct.Yellow, false, ct.None, false)
		}
		fmt.Print(t.Desc[token[0]:token[1]])
		if useColor == true {
			ct.ResetColor()
		}
		pos = token[1]
	}
	fmt.Println(t.Desc[pos:])
}
