package x

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Get(uri string) (io.ReadCloser, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	// hash url with md5
	id := fmt.Sprintf("%x", md5.Sum([]byte(u.String())))
	// cached http get
	return NetCached(u.Hostname(), id, func() (io.ReadCloser, error) {
		resp, err := http.Get(u.String())
		if err != nil {
			return nil, err
		}
		return resp.Body, nil
	})
}

func GetJSON(uri string, data any) error {
	rc, err := Get(uri)
	if err != nil {
		return err
	}
	defer rc.Close()
	if err := json.NewDecoder(rc).Decode(data); err != nil {
		return err
	}
	return nil
}
