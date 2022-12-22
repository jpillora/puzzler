package kernel

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
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
	// valid directory, fetch README.md
	if err := x.CreateFunc("README.md", func() (string, error) {
		return fetchQuestion(year, day, session)
	}); err != nil {
		return err
	}
	// now have a README.md file, extract example
	if err := x.CreateFunc("input-example.txt", func() (string, error) {
		md, err := os.ReadFile("README.md")
		if err != nil {
			return "", err
		}
		codeBlock := "```"
		re := regexp.MustCompile(`(?m)For example.*:\n+` + codeBlock + `\n((.*\n)+?)\n?` + codeBlock + `\n`)
		example := ""
		if m := re.FindSubmatch(md); len(m) > 0 {
			example = string(m[1])
		}
		return example, nil
	}); err != nil {
		return err
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
	// watch files, and run on change
	eg := errgroup.Group{}
	eg.Go(func() error {
		return watch(events)
	})
	eg.Go(func() error {
		return run(year, day, events)
	})
	return eg.Wait()
}

// must have a go.mod file
// valid := exec.Command("go", "mod", "verify").Run() == nil
// if !valid {
// 	modName := fmt.Sprintf("day%d", day)
// 	cmd := exec.Command("go", "mod", "init", modName)
// 	cmd.Dir = yearDir
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	if hasYearDir {
// 		cmd.Dir = filepath.Dir(yearDir)
// 	}
// 	if out, err := cmd.CombinedOutput(); err != nil {
// 		return fmt.Errorf("go mod init failed: %s: %s", err, out)
// 	}
// 	x.Logf("created go.mod file")
// 	// new go.mod file so lets download jpillora/puzller
// 	cmd = exec.Command("go", "get", "github.com/jpillora/puzzler")
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	if out, err := cmd.CombinedOutput(); err != nil {
// 		return fmt.Errorf("go get puzzler failed: %s: %s", err, out)
// 	}
// 	x.Logf("downloaded github.com/jpillora/puzzler")
// }
