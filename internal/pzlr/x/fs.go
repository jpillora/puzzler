package x

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
)

func Read(file string) string {
	b, _ := os.ReadFile(file)
	return string(b)
}

func MkdirAll(dir string) error {
	s, err := os.Stat(dir)
	if errors.Is(err, fs.ErrNotExist) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		Logf("Created directory %s/", dir)
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
		if contents == "" {
			return nil
		}
		if err := os.WriteFile(path, []byte(contents), 0600); err != nil {
			return err
		}
		Logf("Created file %s", path)
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

type fileHash [md5.Size]byte

type FileCache struct {
	files map[string]fileHash
}

func (c *FileCache) Changed(path string) bool {
	if c.files == nil {
		c.files = make(map[string]fileHash)
	}
	f, err := os.Open(path)
	if err != nil {
		return true // failure to check cache reports as a "change"
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return true
	}
	curr := fileHash{}
	copy(curr[:], h.Sum(nil))
	prev := c.files[path]
	c.files[path] = curr
	changed := curr != prev
	return changed
}
