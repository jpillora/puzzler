package x

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Get(url string) (io.ReadCloser, error) {
	// hash url with md5
	id := fmt.Sprintf("%x", md5.Sum([]byte(url)))
	// cached http get
	return Cached(id, func() (io.ReadCloser, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		return resp.Body, nil
	})
}

func GetJSON(url string, data any) error {
	rc, err := Get(url)
	if err != nil {
		return err
	}
	defer rc.Close()
	if err := json.NewDecoder(rc).Decode(data); err != nil {
		return err
	}
	return nil
}
