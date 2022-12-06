package pzlr

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/jpillora/pzlr/internal/pzlr/aoc"
	"github.com/jpillora/pzlr/internal/pzlr/leetcode"
	"github.com/jpillora/pzlr/internal/pzlr/x"
)

type RunWith struct {
	Provider string     `opts:"mode=arg, help=<provider> must be either leetcode or adventofcode"`
	ID       string     `opts:"mode=arg, help=<id> of the provider problem"`
	Flags    x.RunFlags `opts:"mode=embedded"`
}

func Run(w RunWith) error {
	// must have go installed
	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("go is not installed. you can install it\n" +
			"  manually here: https://golang.org/doc/install\n" +
			"  automatically with: curl https://jpillora.com/dotfiles/bin/install-go | bash")
	}
	// must have a go.mod file
	if _, err := os.Stat("go.mod"); errors.Is(err, os.ErrNotExist) {
		if out, err := exec.Command("go", "mod", "init", "pzlr").CombinedOutput(); err != nil {
			return fmt.Errorf("go mod init failed: %s: %s", err, out)
		}
	}
	switch w.Provider {
	case "l", "leetcode":
		return leetcode.Run(w.ID, w.Flags)
	case "a", "aoc", "adventofcode":
		return aoc.Run(w.ID, w.Flags)
	}
	return fmt.Errorf("unknown provider %q", w.Provider)
}
