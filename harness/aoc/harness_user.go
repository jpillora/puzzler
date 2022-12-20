package aoc

import (
	"errors"
	"fmt"
	"log"
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
		if err := runWith(fn, id, input); err != nil {
			return err
		}
		ran++
	}
	if ran == 0 {
		return errors.New("no input files found (expected input-*.txt or input_*.txt)")
	}
	return nil
}

func runWith(fn RunFn, id, input string) error {
	// run part1
	if err := runPartWith(fn, id, false, input); err != nil {
		return err
	}
	// run part2
	if err := runPartWith(fn, id, true, input); err != nil {
		return err
	}
	return nil
}

func runPartWith(fn RunFn, id string, part2 bool, input string) error {
	ts := time.Now()
	p := "1"
	if part2 {
		p = "2"
	}
	defer x.PanicPrint()
	value := fn(part2, input)
	s, ok := value.(string)
	skip := ok && s == "not implemented"
	if skip {
		return nil
	}
	out := ansi.Bright.String(fmt.Sprintf("%v", value))
	if strings.Contains(out, "\n") {
		out = "\n" + out
	}
	log.Printf("run(part%s, input-%s) returned in %s => %s", ansi.Cyan.String(p), ansi.Green.String(id), ansi.Cyan.String(since(ts)), out)
	return nil
}

func since(ts time.Time) string {
	d := time.Since(ts)
	s := d.String()
	re := regexp.MustCompile(`\.\d+`)
	return re.ReplaceAllString(s, "")
}
