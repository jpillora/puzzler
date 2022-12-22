package x

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// type Request struct {
// 	Method string
// 	URL string
// 	Input any
// 	Output any
// }

func GetWith(uri string, headers map[string]string) (io.ReadCloser, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["User-Agent"] = "github.com/jpillora/puzzler"
	// hash url with md5
	h := md5.New()
	h.Write([]byte(u.String()))
	for k, v := range headers {
		h.Write([]byte("|"))
		h.Write([]byte(k))
		h.Write([]byte("|"))
		h.Write([]byte(v))
	}
	id := fmt.Sprintf("%x", h.Sum(nil))
	// cached http get
	return NetCached(u.Hostname(), id, func() (io.ReadCloser, error) {
		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			return nil, err
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		return resp.Body, nil
	})
}

func Get(uri string) (io.ReadCloser, error) {
	return GetWith(uri, nil)
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
