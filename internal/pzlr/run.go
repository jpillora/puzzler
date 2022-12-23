package pzlr

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jpillora/puzzler/internal/pzlr/leetcode"
	"github.com/jpillora/puzzler/internal/pzlr/x"
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
	valid := exec.Command("go", "mod", "verify").Run() == nil
	if !valid {
		modName := "pzlr"
		cmd := exec.Command("go", "mod", "init", modName)
		cmd.Dir = "."
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("go mod init failed: %s: %s", err, out)
		}
		x.Logf("created go.mod file")
		// new go.mod file so lets download jpillora/puzller
		cmd = exec.Command("go", "get", "github.com/jpillora/puzzler")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("go get puzzler failed: %s: %s", err, out)
		}
		x.Logf("downloaded github.com/jpillora/puzzler")
	}
	switch w.Provider {
	case "l", "leetcode":
		return leetcode.Run(w.ID, w.Flags)
	}
	return fmt.Errorf("unknown provider %q", w.Provider)
}
