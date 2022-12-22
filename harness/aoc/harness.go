package aoc

import (
	"log"
	"os"

	"github.com/jpillora/puzzler/harness/aoc/kernel"
	"github.com/jpillora/puzzler/harness/aoc/user"
)

func Harness(fn user.RunFn) {
	if err := harness(fn); err != nil {
		log.Fatalf("harness: %s", err)
	}
}

func harness(fn user.RunFn) error {
	if os.Getenv("AOC_HARNESS") == "1" {
		return user.Harness(fn)
	}
	return kernel.Harness()
}
