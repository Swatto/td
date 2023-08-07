//go:build windows
// +build windows

package printer

type Printer uint8

const (
	DONE     = Printer(0)
	PENDING  = Printer(1)
	EXPIRED  = Printer(2)
	DEADLINE = Printer(3)
	PERIOD   = Printer(4)
)

var signs = []string{"V", "X", "X", " [] ", " x "}

func Sign(status Printer, isNerd bool) string {
	if int(status) > len(signs) {
		return ""
	}
	return signs[uint8(status)]

}
