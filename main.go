package main

import (
	"fmt"
	"os"

	"github.com/jpillora/opts"
	"github.com/jpillora/pzlr/internal/pzlr"
)

func main() {
	w := pzlr.RunWith{}
	opts.New(&w).Parse()
	if err := pzlr.Run(w); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
