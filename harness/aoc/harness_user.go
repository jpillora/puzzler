package aoc

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jpillora/ansi"
	"github.com/jpillora/puzzler/internal/pzlr/x"
)

func user(fn RunFn) error {
	files, err := os.ReadDir(".")
	if err != nil {
		return err
	}
	inputRe := regexp.MustCompile(`^input[-_]([\w_]+)\.txt?$`)
	ran := 0
	for _, file := range files {
		name := file.Name()
		if file.IsDir() {
			continue // ignore dirs
		}
		m := inputRe.FindStringSubmatch(name)
		if len(m) == 0 {
			continue
		}
		id := m[1]
		b, err := os.ReadFile(name)
		if err != nil {
			return err
		}
		if len(b) == 0 {
			continue // ignore empty files
		}
		input := string(b)
		next := runWith(fn, id, input)
		ran++
		if !next {
			break
		}
	}
	if ran == 0 {
		return errors.New("no input files found (expected input-*.txt or input_*.txt)")
	}
	return nil
}

func runWith(fn RunFn, id, input string) bool {
	// run part1
	return runPartWith(fn, id, false, input) &&
		// run part2
		runPartWith(fn, id, true, input)
}

func runPartWith(fn RunFn, id string, part2 bool, input string) (next bool) {
	ts := time.Now()
	p := "1"
	if part2 {
		p = "2"
	}
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		x.PanicPrint(r)
		result(p, id, ts, false, r)
		next = false
	}()
	value := fn(part2, input)
	s, ok := value.(string)
	skip := value == nil || ok && (s == "skip" || s == "not implemented")
	print := !skip
	if print {
		result(p, id, ts, true, value)
	}
	next = true
	return
}

func since(ts time.Time) string {
	d := time.Since(ts)
	s := d.String()
	re := regexp.MustCompile(`\.\d+`)
	return re.ReplaceAllString(s, "")
}

func output(v any) string {
	out := ansi.Bright.String(fmt.Sprintf("%v", v))
	if strings.Contains(out, "\n") {
		out = "\n" + out
	}
	return out
}

func result(p, id string, ts time.Time, success bool, value any) {
	fmt.Print(ansi.Black.String("run(part"))
	fmt.Print(ansi.Cyan.String(p))
	fmt.Print(ansi.Black.String(", input-"))
	fmt.Print(ansi.Green.String(id))
	fmt.Print(ansi.Black.String(") "))
	if success {
		fmt.Print(ansi.Green.String("returned"))
	} else {
		fmt.Print(ansi.Red.String("panicked"))
	}
	fmt.Print(ansi.Black.String(" in "))
	fmt.Print(ansi.Cyan.String(since(ts)))
	fmt.Print(ansi.Black.String(" => "))
	fmt.Print(output(value))
	fmt.Println()
}
