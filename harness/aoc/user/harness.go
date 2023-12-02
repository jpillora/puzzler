package user

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

type RunFn func(part1 bool, input string) any

func Harness(fn RunFn) error {
	inputs := 0
	runs := 0
	doPart := os.Getenv("PART")
	doInput := os.Getenv("INPUT")
	for _, part := range []string{"1", "2"} {
		skipPart := doPart != "" && doPart != part
		for _, kind := range []string{"example", "user"} {
			skipInput := doInput != "" && doInput != kind
			file := "input-" + kind
			b1, ok1 := readInput(file)
			b2, ok2 := readInput(file + "2")
			// choose one of the 2 inputs
			input := ""
			if part == "2" && ok2 {
				file += "2"
				input = string(b2)
			} else if ok1 {
				input = string(b1)
			} else {
				continue
			}
			inputs++
			if skipPart || skipInput {
				skip(part, file)
				continue
			}
			ran, success := runPartWith(fn, part, file, input)
			if ran {
				runs++
			}
			if !success {
				break
			}
		}
	}
	if inputs == 0 {
		return errors.New("no input text files found")
	}
	if runs == 0 {
		x.Logf("skipped all parts/inputs")
	}
	return nil
}

func readInput(name string) ([]byte, bool) {
	b, err := os.ReadFile(name + ".txt")
	if err != nil {
		return nil, false
	}
	if len(b) == 0 {
		return nil, false
	}
	return b, true
}

func runPartWith(fn RunFn, part, file, input string) (ran, success bool) {
	ts := time.Now()
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		x.PanicPrint(r)
		result(part, file, ts, false, r)
		success = false
	}()
	value := fn(part == "2", input)
	s, ok := value.(string)
	skip := value == nil || ok && (s == "skip" || s == "not implemented")
	ran = !skip
	if ran {
		result(part, file, ts, true, value)
	}
	success = true
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

func skip(part, file string) {
	fmt.Print(ansi.Black.String("skip(part"))
	fmt.Print(ansi.Cyan.String(part))
	fmt.Print(ansi.Black.String(", "))
	fmt.Print(ansi.Green.String(file))
	fmt.Print(ansi.Black.String(") "))
	fmt.Println()
}

func result(part, file string, ts time.Time, success bool, value any) {
	fmt.Print(ansi.Black.String("run(part"))
	fmt.Print(ansi.Cyan.String(part))
	fmt.Print(ansi.Black.String(", "))
	fmt.Print(ansi.Green.String(file))
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
