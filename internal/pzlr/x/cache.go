package x

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/jpillora/maplock"
)

type Cache struct {
	locks maplock.Maplock
	Base  string
}

var defaultCache = Cache{
	Base: "cache-pzlr",
}

func NewCache(name string) *Cache {
	return defaultCache.Fork(name)
}

func (c *Cache) Fork(name string) *Cache {
	return &Cache{
		Base: c.Base + "-" + name,
	}
}

func (c *Cache) Cached(id string, fn func() (io.ReadCloser, error)) (io.ReadCloser, error) {
	c.locks.Lock(id)
	defer c.locks.Unlock(id)
	file := c.Base + "-" + id + ".bin"
	path := filepath.Join(os.TempDir(), file)
	// already cached to disk within the last day, return file as the read-closer
	s, err := os.Stat(path)
	if err == nil && !s.IsDir() && time.Since(s.ModTime()) < 24*time.Hour {
		f, err := os.Open(path)
		if err == nil {
			return f, nil
		}
	}
	// not cached, return the function (likely a network request)
	rc, err := fn()
	if err != nil {
		return rc, err
	}
	if rc == nil {
		return nil, nil
	}
	// prepare tee file
	tf, err := os.Create(path)
	if err != nil {
		return rc, nil // ignore cache file errors
	}
	// return read closer which also tees of to a tmp file
	return &teeReadCloser{
		teeFile: tf,
		rc:      rc,
	}, nil
}

func Cached(id string, fn func() (io.ReadCloser, error)) (io.ReadCloser, error) {
	return defaultCache.Cached(id, fn)
}

type teeReadCloser struct {
	teeFile *os.File
	rc      io.ReadCloser
}

func (r *teeReadCloser) Read(p []byte) (int, error) {
	if r.rc == nil {
		return 0, io.EOF
	}
	n, err := r.rc.Read(p)
	if n > 0 && !r.closed() {
		r.teeFile.Write(p[:n]) // ignore cache file write errors
	}
	if err == io.EOF {
		r.closeTee()
	}
	return n, err
}

func (r *teeReadCloser) Close() error {
	r.closeTee()
	return r.rc.Close()
}

func (r *teeReadCloser) closeTee() {
	if !r.closed() {
		r.teeFile.Close()
		r.teeFile = nil
	}
}

func (r *teeReadCloser) closed() bool {
	return r.teeFile == nil
}
