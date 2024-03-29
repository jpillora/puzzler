package leetcode

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cweill/gotests/gotests/process"
	"github.com/jpillora/puzzler/internal/pzlr/x"
	gotestsum "gotest.tools/gotestsum/cmd"
)

func Run(id string, flags x.RunFlags) error {
	spec, err := getProblemSpec(id)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://leetcode.com/problems/%s/", spec.Slug())
	fmt.Printf("Found problem #%s %s\n", spec.ID(), url)
	// create problem directory if it doesn't exist
	dir := filepath.Join("leetcode", spec.ID())
	if _, err := os.Stat(dir); errors.Is(err, fs.ErrNotExist) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("problem %s dir could not be created: %w", spec.ID(), err)
		}
		fmt.Printf("Created directory %s/\n", dir)
	}
	// create README.txt if it doesn't exist
	if err := x.CreateFunc(filepath.Join(dir, "README.txt"), func() (string, error) {
		return fetchQuestionText(spec.Slug())
	}); err != nil {
		return err
	}
	// create code.go if it doesn't exist
	codeFile := filepath.Join(dir, "code.go")
	if _, err := os.Stat(codeFile); errors.Is(err, fs.ErrNotExist) {
		stub, err := getProblemCode(spec.Slug())
		if err != nil {
			return err
		}
		comment := "// " + url + "\n"
		code := fmt.Sprintf("package p%s\n\n%s%s", spec.ID(), comment, stub)
		if err := os.WriteFile(codeFile, []byte(code), 0755); err != nil {
			return fmt.Errorf("problem %s code could not be created: %w", spec.ID(), err)
		}
		fmt.Printf("Created stub answer file %s\n", codeFile)
	}
	// create code_test.go if it doesn't exist
	testFile := filepath.Join(dir, "code_test.go")
	if _, err := os.Stat(testFile); errors.Is(err, fs.ErrNotExist) {
		buff := bytes.Buffer{}
		process.Run(&buff, []string{codeFile}, &process.Options{
			AllFuncs: true,
			Subtests: true,
		})
		// modify the code
		// TODO: use a proper AST parser
		code := buff.String()
		code = regexp.MustCompile(`Generated Test.+\n`).ReplaceAllString(code, "")
		if !strings.HasPrefix(code, "package ") {
			return fmt.Errorf("problem %s test could not generated: %s", spec.ID(), code)
		}
		code = regexp.MustCompile(`\bTest_[a-z]`).ReplaceAllStringFunc(code, func(b string) string {
			return "Test" + strings.ToUpper(string(b[5]))
		})
		code = regexp.MustCompile(`\bargs\b`).ReplaceAllString(code, "input")
		code = regexp.MustCompile(`\bwant\b`).ReplaceAllString(code, "output")
		fmt.Printf("Created stub test file %s\n", testFile)
		if err := os.WriteFile(testFile, []byte(code), 0755); err != nil {
			return fmt.Errorf("problem %s code could not be created: %w", spec.ID(), err)
		}
	}
	// open in browser
	if flags.Open {
		if _, err := exec.LookPath("open"); err == nil {
			if err := exec.Command("open", url).Run(); err != nil {
				fmt.Printf("Failed to run 'open %s': %s\n", url, err)
			}
		}
	}
	// start gotestsum watch
	if err := os.Chdir(dir); err != nil {
		return err
	}
	fmt.Print("Start dev. ")
	return gotestsum.Run("leetcode", []string{"--watch", "--format", "testname", "--hide-summary", "skipped,failed,errors"})
}
