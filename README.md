# Puzzler

[![CI](https://github.com/jpillora/puzzler/workflows/CI/badge.svg)](https://github.com/jpillora/puzzler/actions?workflow=CI)

Puzzler is a Go (golang) program wrapper to assist in developing solutions to various programming puzzles. An internet connection is only required to fetch the questions. Current, it supports Advent of Code.

## Quick start

Use this template repo https://github.com/jpillora/aoc-in-go

Alternatively, you can setup manually by following the steps below

## Advent of Code

[Advent of Code](https://adventofcode.com) is a yearly series of programming questions based on the [Advent Calendar](https://en.wikipedia.org/wiki/Advent_calendar). For each day leading up to christmas, there is one question released, and from the second it is released, there is a timer running and a leaderboard showing who solved it first.

1. Manually install Go from https://go.dev/dl/ or from brew, etc

1. Make an AOC solutions directory and initialise it as a Go module

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
	$ go run code.go
	Created file README.md
	Created file input-example.txt
	run(part1, input-example) returned in 37µs => 42
	# update code.go to return 43 and...
	file changed code.go
	run(part1, input-example) returned in 34µs => 43
	```

1. You can find your question in `README.md`, iterate on `code.go` until you get the answer


#### AOC Session

**Optionally**, you can set `AOC_SESSION` to your adventofcode.com `session` cookie. That is:

* Login with your browser
* Open developer tools > Application/Storage > Cookies
* Retrieve the contents of `session`
* Export it as `AOC_SESSION`

With your session, `puzzler` will download your user-specifc `input-user.txt` and also update `README.md` with part 2 of the question once you've completed part 1.

Current, your session is NOT used to submit your answer. You still need to login to https://adventofcode.com to submit.

## Leetcode

TODO

* https://leetcode.com
* Problem: Leetcode's programming model is different from AOC. Its not text in, text out. There are structured questions and code samples and other differences that make it tricky.
* Goal: Make the harness API work like:

	```go
	func main() {
		leetcode.Harness(run)
	}
	```

#### TODO

* Submission
	* Requires `AOC_SESSION`
	* Press `s` to submit your most recent pending result the next part
* Press `r` to re-run your program
* Press `c` to cancel any running program