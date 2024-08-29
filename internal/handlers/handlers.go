package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/prplx/go-sentry-tunnel/internal/config"
	"github.com/prplx/go-sentry-tunnel/internal/lib/url"
)

type payload struct {
	DSN string `json:"dsn"`
}

func HandleTunnel(w http.ResponseWriter, r *http.Request, c *config.Config) {
	host := r.Header.Get("Origin")
	for _, h := range c.AllowOrigins {
		if h == host || h == "*" {
			w.Header().Set("Access-Control-Allow-Origin", h)
		}
	}

	envelope, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorToResponse(w)
		return
	}
	defer r.Body.Close()

	piece := strings.Split(string(envelope), "\n")[0]
	var result payload
	decoder := json.NewDecoder(strings.NewReader(piece))
	err = decoder.Decode(&result)
	if err != nil {
		writeErrorToResponse(w)
		return
	}
	if result.DSN == "" {
		writeErrorToResponse(w)
		return
	}

	err = url.ValidateDSN(result.DSN, c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	upstreamURL, err := url.BuildSentryUpstreamURL(result.DSN)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := postWithTimeout(upstreamURL, envelope, c.RequestTimeout)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func writeErrorToResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func postWithTimeout(URL string, payload []byte, timeout time.Duration) (*http.Response, error) {
	req, err := http.NewRequest("POST", URL, bytes.NewReader((payload)))
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: timeout,
	}

	return client.Do(req)
}
