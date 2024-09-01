package http

import (
	"bytes"
	"net/http"
	"time"
)

type HTTPClient struct {
}

func (c *HTTPClient) PostWithTimeout(URL string, payload []byte, timeout time.Duration) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewReader((payload)))
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: timeout,
	}

	return client.Do(req)
}
