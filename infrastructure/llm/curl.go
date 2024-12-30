package llm

import (
	"bytes"
	"net/http"
)

func NewRequest(method string, url string, body []byte, isStream bool) (*http.Request, error) {
	var (
		req *http.Request
		err error
	)
	if len(body) > 0 {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if isStream {
		req.Header.Set("Accept", "text/event-stream")
	}
	req.Header.Set("Cache-Control", "no-cache")

	return req, nil
}
