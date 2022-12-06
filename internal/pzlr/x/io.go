package x

import (
	"bufio"
	"fmt"
	"os"
)

func ReadCode() (string, error) {
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
