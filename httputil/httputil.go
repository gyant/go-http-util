package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type APIClient interface {
	getBaseURL() *url.URL
	getHTTPClient() *http.Client
	getUserAgent() string
}

type HTTPUtil struct{}

func (h *HTTPUtil) NewRequest(c APIClient, method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.getBaseURL().ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.getUserAgent())
	return req, nil
}

func (h *HTTPUtil) Do(c APIClient, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.getHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err
}

func (h *HTTPUtil) Test() {
	fmt.Println("OH HI FROM HTTPUTIL TEST.")
}
