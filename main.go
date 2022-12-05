package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cweill/gotests"
	"github.com/jpillora/opts"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	type config struct {
		ID int `opts:"mode=arg, help=problem <id>entifier must be a number between 1-9999"`
	}
	c := config{}
	opts.New(&c).Parse()

	d := fmt.Sprintf("%04d", c.ID)
	if _, err := os.Stat(d); err == nil {
		return fmt.Errorf("problem %d already exists", c.ID)
	}

	// fmt.Println("paste in the leetcode Go stub here:")
	// function, err := readCode()
	// if err != nil {
	// 	return fmt.Errorf("read code: %w", err)
	// }

	// if err := os.Mkdir(d, 0755); err != nil {
	// 	return fmt.Errorf("problem %s dir could not be created: %w", d, err)
	// }
	// fmt.Printf("created %s/\n", d)

	// contents := fmt.Sprintf("package p%s\n\n%s", d, function)
	// if err := os.WriteFile(fmt.Sprintf("%s/code.go", d), []byte(contents), 0755); err != nil {
	// 	return fmt.Errorf("problem %s code could not be created: %w", d, err)
	// }
	// fmt.Printf("created %s/code.go\n", d)

	gen, err := gotests.GenerateTests(d, &gotests.Options{})
	if err != nil {
		return err
	}
	for _, g := range gen {
		fmt.Printf("generated %s\n", g.Path)
	}
	fmt.Printf("created %s/code_test.go\n", d)
	return nil
}

func readCode() (string, error) {
	code := ""
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		t := scanner.Text()
		if t == "" {
			break
		}
		code += t
		if t == "}" {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("read failed: %w", err)
	}
	if code == "" {
		return "", fmt.Errorf("no code found")
	}
	return code, nil
}
