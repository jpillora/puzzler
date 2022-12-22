# Puzzler

[![CI](https://github.com/jpillora/puzzler/workflows/CI/badge.svg)](https://github.com/jpillora/puzzler/actions?workflow=CI)

Puzzler is a Go (golang) program wrapper to assist in developing solutions to various programming puzzles. An internet connection is only required to fetch the questions. Current, it supports:

* Leetcode https://leetcode.com
* 

## Leetcode

TODO HARNESS

## Advent of Code

[Advent of Code](https://adventofcode.com) is a yearly series of programming questions based on the [Advent Calendar](https://en.wikipedia.org/wiki/Advent_calendar). For each day leading up to christmas, there is one question released, and from the second it is released, there is a timer running and a leaderboard showing who solved it first.

1. Install Go 1.18+ from https://go.dev/dl/, or using [my install script](https://github.com/jpillora/dotfiles/blob/main/bin/install-go):

	```sh
	curl https://jpillora.com/dotfiles/bin/install-go | bash
	```

1. Make an Advent of Code directory/repository (choose your own name) and initialise it as a Go module

	```
	mkdir -p my-aoc-solutions
	cd my-aoc-solutions
	go mod init my-aoc-solutions
	```

1. Make a directory and file `YYYY/DD/code.go` where `YYYY`/`DD` is the AOC year/day you'd like to attempt

	```
	mkdir -p 2022/09/
	touch 2022/09/code.go
	```

1. Paste the following into `code.go`

	```go
	package main

	import (
		"github.com/jpillora/puzzler/harness/aoc"
	)

	func main() {
		aoc.Harness(run)
	}

	func run(part2 bool, input string) any {
		if part2 {
			return "not implemented"
		}
		return 42
	}
	```

1. Run `code.go`

	```sh
	go run main.go
	```

1. You should see

	```sh
	Created file README.md
	Created file input-example.txt
	run(part1, input-eg) returned in 37µs => 42
	# update code.go to return 43 and...
	file changed code.go
	run(part1, input-eg) returned in 34µs => 43
	```

	**Optionally** set `AOC_SESSION` to your adventofcode.com `session` cookie and it will also download your specific user input (`input-user.txt`)

1. TODO submission

	* Requires `AOC_SESSION`
	* Press `r` to re-run your program
	* Press `s` to submit your most recent successful result the next part