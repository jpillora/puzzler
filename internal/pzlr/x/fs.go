package x

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func MkdirAll(dir string) error {
	s, err := os.Stat(dir)
	if errors.Is(err, fs.ErrNotExist) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		fmt.Printf("Created directory %s/\n", dir)
	} else if err != nil {
		return err
	} else if !s.IsDir() {
		return fmt.Errorf("failed to mkdir: %s: is already a file", dir)
	}
	return nil
}

func CreateFunc(path string, fn func() (string, error)) error {
	s, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		contents, err := fn()
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(contents), 0755); err != nil {
			return err
		}
		fmt.Printf("Created file %s\n", path)
		return nil
	}
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("failed to write: %s: is already a directory", path)
	}
	return nil

}

func Create(file, contents string) error {
	return CreateFunc(file, func() (string, error) { return contents, nil })
}
