package x

import (
	"fmt"

	"github.com/jpillora/ansi"
)

func Logf(format string, args ...interface{}) {
	fmt.Printf(ansi.Black.String(format+"\n"), args...)
}
