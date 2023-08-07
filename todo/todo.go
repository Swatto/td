/*
This package contains todo items and their functions
*/
package todo

import (
	"fmt"
	"regexp"
	"time"

	ct "github.com/daviddengcn/go-colortext"
	"github.com/swatto/td/printer"
)

// Data struct for a single todo
type Todo struct {
	ID       int       `json:"id"`
	Desc     string    `json:"desc"`
	Status   string    `json:"status"`
	Modified time.Time `json:"modified"`
	Period   int       `json:"period"`
	Deadline time.Time `json:"deadline"`
	Created  time.Time `json:"created"`
}

// Prints the todo to screen
func (t *Todo) MakeOutput(useColor bool, isNerd bool) {
	var symbole string
	var color ct.Color

	if t.Status == "done" {
		color = ct.Green
		symbole = printer.Sign(printer.DONE, isNerd)
	} else {
		color = ct.Red
		if t.IsExpired() {
			symbole = printer.Sign(printer.EXPIRED, isNerd)
		} else {
			symbole = printer.Sign(printer.PENDING, isNerd)
		}
	}

	hashtagReg := regexp.MustCompile(`#[^\\s]*`)

	fmt.Printf("%02d | ", t.ID)
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
	fmt.Printf("%-25s", t.Desc[pos:])
	if !t.Deadline.IsZero() {
		if (t.Deadline.Minute() == 59 && t.Deadline.Hour() == 23) ||
			(t.Deadline.Minute() == 0 && t.Deadline.Hour() == 0) {
			fmt.Print(printer.Sign(printer.DEADLINE, isNerd), t.Deadline.Format("Mon, 02 Jan"))
		} else {
			fmt.Print(printer.Sign(printer.DEADLINE, isNerd), t.Deadline.Format("Mon, 02 Jan 15:04"))
		}
	}
	if t.Period != 0 {
		fmt.Print(printer.Sign(printer.PERIOD, isNerd), t.Period)
	}
	fmt.Println(" ")
}

func (t *Todo) IsExpired() bool {
	return !t.Deadline.IsZero() && t.Deadline.Before(time.Now())
}
