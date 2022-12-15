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
	testCode := fmt.Sprintf(test, pkg, "``")
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

import (
	"errors"
)

func run(input string) (any, error) {
	return nil, errors.New("not implemented")
}`

const test = `package %s

import (
	"testing"
)

const input = %s

func TestEg(t *testing.T) {
	ans, err := run(input)
	if err != nil {
		t.Fatalf("eg => %%s\n", err)
	}
	t.Logf("eg => %%v\n", ans)
}

func TestPart1(t *testing.T) {
	ans, err := run(input)
	if err != nil {
		t.Fatalf("part1 => %%s\n", err)
	}
	t.Logf("part1 => %%v\n", ans)
}

func TestPart2(t *testing.T) {
	ans, err := run(input)
	if err != nil {
		t.Fatalf("part2 => %%s\n", err)
	}
	t.Logf("part2 => %%v\n", ans)
}
`
