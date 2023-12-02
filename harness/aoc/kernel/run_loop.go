package kernel

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jpillora/puzzler/internal/pzlr/x"
)

func run(year, day int, part2 bool, events <-chan event) error {
	// control process start
	start := make(chan bool, 1)
	// singleton process to control
	var proc *exec.Cmd
	go func() {
		for e := range events {
			if e == restart {
				if proc != nil && proc.Process != nil {
					proc.Process.Kill()
					x.Logf("killed code.go process")
				}
				if len(start) == 0 {
					start <- true
				}
			}
		}
	}()
	// go run code.go in a loop, pass state through environment variables
	for {
		proc = exec.Command("go", "run", "code.go")
		proc.Env = append(
			os.Environ(),
			"AOC_HARNESS=1",
			fmt.Sprintf("AOC_YEAR=%d", year),
			fmt.Sprintf("AOC_DAY=%d", day),
			fmt.Sprintf("AOC_PART2=%v", part2),
		)
		proc.Stdout = os.Stdout
		proc.Stderr = os.Stderr
		if err := proc.Start(); err != nil {
			x.Logf("fail: go run code.go: %s", err)
			<-start
			continue
		}
		proc.Wait()
		// exited
		num := proc.ProcessState.ExitCode()
		proc = nil
		if num != 0 {
			x.Logf("code.go exited with code %d, waiting for file change...", num)
		}
		// block until restarted
		<-start
	}
}
