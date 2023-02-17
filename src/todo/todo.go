package todo

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	ct "github.com/daviddengcn/go-colortext"
	"umutsevdi/td/printer"
)

type Todo struct {
	ID       int64     `json:"id"`
	Desc     string    `json:"desc"`
	Status   string    `json:"status"`
	Modified string    `json:"modified"`
	Period   int       `json:"period"`
	Deadline time.Time `json:"deadline"`
	Created  time.Time `json:"created"`
}

func (t *Todo) MakeOutput(useColor bool) {
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
	fmt.Print(t.Desc[pos:])
	if !t.Deadline.IsZero() {
		fmt.Print(printer.DeadlineSign, t.Deadline.Format("Mon, 02 Jan 15:04"))
	}
	if t.Period != 0 {
		fmt.Print(printer.PeriodSign, t.Period)
	}
	fmt.Println(" ")
}
