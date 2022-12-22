package aoc

import (
	"fmt"
	"log"
	"os"

	"github.com/jpillora/ansi"
)

type RunFn func(part1 bool, input string) any

func Harness(fn RunFn) {
	if err := harness(fn); err != nil {
		log.Fatalf("harness: %s", err)
	}
}

func harness(fn RunFn) error {
	if os.Getenv("AOC_HARNESS") == "1" {
		return user(fn)
	}
	return kernel()
}

func logf(format string, args ...interface{}) {
	fmt.Printf(ansi.Black.String(format+"\n"), args...)
}
