package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	p "github.com/Swatto/td/printer"
	"github.com/daviddengcn/go-colortext"
)

type Todo struct {
	Id       int64  `json:"id"`
	Desc     string `json:"desc"`
	Status   string `json:"status"`
	Modified string `json:"modified"`
}

func (t *Todo) MakeOutput() {
	var symbole string
	var color ct.Color

	if t.Status == "done" {
		color = ct.Green
		symbole = p.OkSign
	} else {
		color = ct.Red
		symbole = p.KoSign
	}

	hashtag_reg := regexp.MustCompile("#[^\\s]*")

	space_count := 6 - len(strconv.FormatInt(t.Id, 10))

	fmt.Print(strings.Repeat(" ", space_count), t.Id, " | ")
	ct.ChangeColor(color, false, ct.None, false)
	fmt.Print(symbole)
	ct.ResetColor()
	fmt.Print(" ")
	pos := 0
	for _, token := range hashtag_reg.FindAllStringIndex(t.Desc, -1) {
		fmt.Print(t.Desc[pos:token[0]])
		ct.ChangeColor(ct.Yellow, false, ct.None, false)
		fmt.Print(t.Desc[token[0]:token[1]])
		ct.ResetColor()
		pos = token[1]
	}
	fmt.Println(t.Desc[pos:])
}
