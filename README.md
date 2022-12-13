# Puzzler

[![CI](https://github.com/jpillora/puzzler/workflows/CI/badge.svg)](https://github.com/jpillora/puzzler/actions?workflow=CI)

Puzzler `pzlr(1)` is a command-line tool to **locally** develop solutions to various programming puzzles using Go (golang). You only need to be connected to the internet to fetch the questions. Current, it supports

* Leetcode https://leetcode.com
* Advent of Code https://adventofcode.com

### Install

**Binaries**

<!-- WHEN PUBLIC
[![Releases](https://img.shields.io/github/release/jpillora/pzlr.svg)](https://github.com/jpillora/puzzler/releases)
[![Releases](https://img.shields.io/github/downloads/jpillora/pzlr/total.svg)](https://github.com/jpillora/puzzler/releases) -->

Find [the latest pre-compiled binaries here](https://github.com/jpillora/puzzler/releases/latest)  or download and install it now with:

```sh
# download an install the pzlr command with
curl https://i.jpillora.com/puzzler! | bash
```

**Source**

```sh
go get github.com/jpillora/puzzler/cmd/pzlr@latest
```

### Examples

#### Leetcode

1. Create directory

   ```shell
   mkdir pzlr
   cd pzlr
   ```

1. Open terminal, start leetcode problem `151` with:

   ```shell
   pzlr --open leetcode 151

   Found problem #0151 https://leetcode.com/problems/reverse-words-in-a-string/
   Created directory leetcode/0151/
   Fetching problem code for reverse-words-in-a-string...
   Created stub answer file leetcode/0151/code.go
   Created stub test file leetcode/0151/code_test.go
   Starting gotestsum: Watching 1 directories. Use Ctrl-c to to stop a run or exit.
   ```

1. Open in VS Code (or another editor)

   ```shell
   code leetcode/0151
   ```

1. File `leetcode/0151/code.go` will contain:

   ```go
   package p0151

   func reverseWords(s string) string {
       
   }
   ```

1. File `leetcode/0151/code_test.go` will contain:

   ```go
   package p0151

   import "testing"

   func TestReverseWords(t *testing.T) {
   	type input struct {
   		s string
   	}
   	tests := []struct {
   		name string
   		input input
   		output string
   	}{
   		// TODO: Add test cases.
   	}
   	for _, tt := range tests {
   		tt := tt
   		t.Run(tt.name, func(t *testing.T) {
   			t.Parallel()
   			if got := reverseWords(tt.input.s); got != tt.output {
   				t.Errorf("reverseWords() = %v, output %v", got, tt.output)
   			}
   		})
   	}
   }
   ```

1. You will need to read the question and copy the example input/outputs into the `code_test.go` file

   ```go
		// for example, here are two test cases for problem 151
		{
			input:  input{s: "the sky is blue"},
			output: "blue is sky the",
		},
		{
			input:  input{s: "  hello world  "},
			output: "world hello",
		},
   ```

1. You will need to implement the solution in `code.go`

   ```go
   // for example, here is one solution to problem 151
   func reverseWords(s string) string {
   	spaces := regexp.MustCompile(`\s+`)
   	words := spaces.Split(strings.TrimSpace(s), -1)
   	last := len(words) - 1
   	mid := last / 2
   	for i := range words {
   		if i > mid {
   			break
   		}
   		words[i], words[last-i] = words[last-i], words[i]
   	}
   	return strings.Join(words, " ")
   }
   ```

1. Once tests are passing, you will need to copy your solution into `leetcode.com` and submit there (_TODO submit via CLI_)


#### Advent of Code

TODO

### Future features

* Leetcode
	* Improve test stub file
	* Implement (or borrow) leetcode login code to allow `pzlr leetcode NNNN --submit` (automatically submits if all tests are passing)
* Advent of Code
	* Improve test stub file

### Caveats

* Only supports Go (but could support other languages with a PR)
	* Should be implemented using a `Language` interface which holds the differences
* Go must be installed (if you're brave, you can run `curl https://jpillora.com/dotfiles/bin/`[`install-go`](https://github.com/jpillora/dotfiles/blob/main/bin/install-go) ` | bash`)
* Unit test cases need to be manually filled in (`code_test.go` will contain an empty "test table")
* Answers need to be manually submitted (`code.go` will contain the link to submission page)

### Credits

* https://github.com/cweill/gotests is embedded
* https://github.com/gotestyourself/gotestsum is embedded