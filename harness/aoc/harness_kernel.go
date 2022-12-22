package aoc

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/sync/errgroup"
)

type event int

const (
	restart event = iota
)

func kernel() error {
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
	re := fmt.Sprintf(`%c(\d{4})%c(\d{2})%ccode.go`, filepath.Separator, filepath.Separator, filepath.Separator)
	m := regexp.MustCompile(re).FindStringSubmatch(abs)
	if len(m) == 0 {
		return fmt.Errorf("code.go file must be in directory YYYY/DD/code.go")
	}
	year := m[1]
	day := m[2]
	// must have a go.mod file
	valid := exec.Command("go", "mod", "verify").Run() == nil
	if !valid {
		modName := fmt.Sprintf("day%s", day)
		if out, err := exec.Command("go", "mod", "init", modName).CombinedOutput(); err != nil {
			return fmt.Errorf("go mod init failed: %s: %s", err, out)
		}
		logf("created go.mod file")
	}
	//
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

func watch(events chan event) error {
	// code.go exists, check other files too
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()
	// watch directory
	if err = w.Add("."); err != nil {
		return err
	}
	// process events
	restartLast := time.Time{}
	restartPattern := regexp.MustCompile(`(.+\.go|input.+\.txt)`)
	for {
		select {
		case event, ok := <-w.Events:
			if ok && restartPattern.MatchString(event.Name) && time.Since(restartLast) > 250*time.Millisecond {
				logf("file changed %s", event.Name)
				events <- restart
				restartLast = time.Now()
			}
		case err, ok := <-w.Errors:
			if ok {
				return fmt.Errorf("watcher error: %w", err)
			}
		}
	}
}

func run(year, day string, events <-chan event) error {
	// control process start
	start := make(chan bool, 1)
	// "global" process to remote kill
	var proc *exec.Cmd
	go func() {
		for e := range events {
			if e == restart {
				if proc != nil && proc.Process != nil {
					proc.Process.Kill()
					logf("killed code.go process")
				}
				if len(start) == 0 {
					start <- true
				}
			}
		}
	}()
	// go run code.go in a loop
	for {
		proc = exec.Command("go", "run", "code.go")
		proc.Env = append(
			os.Environ(),
			"AOC_HARNESS=1",
			fmt.Sprintf("AOC_YEAR=%s", year),
			fmt.Sprintf("AOC_DAY=%s", day),
		)
		proc.Stdout = os.Stdout
		proc.Stderr = os.Stderr
		if err := proc.Start(); err != nil {
			logf("fail: go run code.go: %s", err)
			continue
		}
		proc.Wait()
		// exited
		num := proc.ProcessState.ExitCode()
		proc = nil
		if num != 0 {
			logf("code.go exited with code %d, waiting for file change...", num)
		}
		// block until restarted
		<-start
	}
}
