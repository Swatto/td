package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	ct "github.com/daviddengcn/go-colortext"
	"github.com/swatto/td/printer"
)

type todo struct {
	ID       int64  `json:"id"`
	Desc     string `json:"desc"`
	Status   string `json:"status"`
	Modified string `json:"modified"`
}

func (t *todo) MakeOutput(useColor bool) {
	var symbole string
	var color ct.Color

	if t.Status == "done" {
		color = ct.Green
		symbole = printer.OkSign
	} else {
		color = ct.Red
		symbole = printer.KoSign
	}

	hashtagReg := regexp.MustCompile(`#[^\\s]*`)

	spaceCount := 6 - len(strconv.FormatInt(t.ID, 10))

	fmt.Print(strings.Repeat(" ", spaceCount), t.ID, " | ")
	if useColor {
		ct.ChangeColor(color, false, ct.None, false)
	}
	fmt.Print(symbole)
	if useColor {
		ct.ResetColor()
	}
	fmt.Print(" ")
	pos := 0
	for _, token := range hashtagReg.FindAllStringIndex(t.Desc, -1) {
		fmt.Print(t.Desc[pos:token[0]])
		if useColor {
			ct.ChangeColor(ct.Yellow, false, ct.None, false)
		}
		fmt.Print(t.Desc[token[0]:token[1]])
		if useColor {
			ct.ResetColor()
		}
		pos = token[1]
	}
	fmt.Println(t.Desc[pos:])
}
