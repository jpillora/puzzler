package aoc

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jpillora/puzzler/internal/pzlr/x"
	gotestsum "gotest.tools/gotestsum/cmd"
)

func Run(day string, flags x.RunFlags) error {
	d, err := strconv.Atoi(day)
	if err != nil {
		return errors.New("expected integer day")
	}
	y := 2022
	dir := filepath.Join("aoc", strconv.Itoa(y), fmt.Sprintf("%02d", d))
	if err := x.MkdirAll(dir); err != nil {
		return err
	}
	if err := x.CreateFunc(filepath.Join(dir, "README.md"), func() (string, error) {
		return fetchQuestion(y, d)
	}); err != nil {
		return err
	}
	pkg := "day" + strconv.Itoa(d)
	ansCode := fmt.Sprintf(ans, pkg)
	if err := x.Create(filepath.Join(dir, "code.go"), ansCode); err != nil {
		return err
	}
	testCode := fmt.Sprintf(test, pkg)
	if err := x.Create(filepath.Join(dir, "code_test.go"), testCode); err != nil {
		return err
	}
	fmt.Printf("Ready: %s\n", dir)
	// start gotestsum watch
	if err := os.Chdir(dir); err != nil {
		return err
	}
	fmt.Print("Start dev. ")
	return gotestsum.Run("aoc", []string{"--watch", "--format", "standard-verbose", "--hide-summary", "skipped,failed,errors"}) //, , "--format", "testname"
}

const ans = `package %s

func run(part1 bool, input string) any {
	return 42
}
`

const test = `package %s

import (
	"fmt"
	"runtime/debug"
	"testing"
)

func TestCode(t *testing.T) {
	print := func() {
		r := recover()
		if r == nil {
			return
		}
		fmt.Printf("%%v\n\n%%s\n", r, string(debug.Stack()))
		t.Fail()
	}
	defer print()
	fmt.Print("example input1 => ")
	fmt.Printf("%%v\n", run(true, inputEg))
	fmt.Print("example input2 => ")
	fmt.Printf("%%v\n", run(false, inputEg))
	fmt.Print("   user input1 => ")
	fmt.Printf("%%v\n", run(true, inputUser))
	fmt.Print("   user input2 => ")
	fmt.Printf("%%v\n", run(false, inputUser))
}

const inputEg = ` + "``" + `

const inputUser = ` + "``" + `
`
