package aoc

import (
	"log"
	"os"
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
