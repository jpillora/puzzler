package kernel

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jpillora/puzzler/internal/pzlr/x"
	"golang.org/x/sync/errgroup"
)

type event int

const (
	restart event = iota
)

func Harness() error {
	// must have go installed
	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("go is not installed. you can install it\n" +
			"  manually here: https://golang.org/doc/install\n" +
			"  automatically with: curl https://jpillora.com/dotfiles/bin/install-go | bash")
	}
	// must have a code.go file
	code, err := os.Stat("code.go")
	if err != nil || code.IsDir() {
		return fmt.Errorf("code.go file not found. please create a YYYY/DD/code.go file")
	}
	// match year/day path
	abs, err := filepath.Abs(code.Name())
	if err != nil {
		return err
	}
	dateDir := filepath.Dir(abs)
	day := -1
	if name := filepath.Base(dateDir); regexp.MustCompile(`^\d\d?$`).MatchString(name) {
		day, _ = strconv.Atoi(name)
	} else {
		return fmt.Errorf("parent directory /%s/code.go must be format /DD/code.go representing the AOC day", name)
	}
	yearDir := filepath.Dir(dateDir)
	year := time.Now().Year()
	if name := filepath.Base(yearDir); regexp.MustCompile(`^\d{4}$`).MatchString(name) {
		year, _ = strconv.Atoi(name)
	}
	// optionally provide session cookie
	session := os.Getenv("AOC_SESSION")
	// always fetch README.md
	readme, err := fetchQuestion(year, day, session)
	if err != nil && x.Read("README.md") == "" {
		// only error if no README.md exists
		// if you are offline but have all files, it should still work
		return err
	}
	if err := x.Create("README.md", readme); err != nil {
		return err
	}
	// determine what part user is up to
	part2 := strings.Contains(readme, "Part Two")
	// determine if we should update the existing README
	readmeDisk := x.Read("README.md")
	part2Disk := strings.Contains(readmeDisk, "Part Two")
	if part2 && !part2Disk {
		if strings.Contains(readme, readmeDisk) {
			// strict subset, safe to replace
			if err := os.WriteFile("README.md", []byte(readme), 0600); err != nil {
				return err
			}
			x.Logf("README.md was updated with Part 2")
		} else {
			x.Logf("Part 2 is available, delete your README.md to updated version")
		}
	}
	// now have a README.md file, extract example/s
	for _, n := range []int{1, 2} {
		file := "input-example"
		if n == 2 {
			file += "2"
		}
		file += ".txt"
		if err := x.CreateFunc(file, func() (string, error) {
			codeBlock := "```"
			re := regexp.MustCompile(`(?m)For example.*:\n+` + codeBlock + `\n((.*\n)+?)\n?` + codeBlock + `\n`)
			matches := re.FindAllStringSubmatch(readme, -1)
			if len(matches) == 0 {
				return "", errors.New("no examples found in README.md")
			}
			if n == 2 && len(matches) == 1 {
				return "", nil // skip
			}
			m := matches[n-1]
			return strings.TrimSpace(m[1]), nil
		}); err != nil {
			return err
		}
	}
	// if we have a AOC_SESSION cookie, fetch input
	if session != "" {
		if err := x.CreateFunc("input-user.txt", func() (string, error) {
			return fetchUserInput(year, day, session)
		}); err != nil {
			return err
		}
	}
	//events in file watcher, out proc runner
	events := make(chan event)
	// begin watch files, and begin run-loop
	eg := errgroup.Group{}
	eg.Go(func() error {
		return watch(events)
	})
	eg.Go(func() error {
		return run(year, day, events)
	})
	// TODO:
	// eg.Go listen-keyboard-inputs
	// s -> submit to AOC
	// r -> rerun
	// q -> quit
	return eg.Wait()
}
