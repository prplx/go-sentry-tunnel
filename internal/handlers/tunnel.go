package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"go-sentry-tunnel/internal/config"
	"go-sentry-tunnel/internal/lib/sl"
	"go-sentry-tunnel/internal/lib/url"
)

type payload struct {
	DSN string `json:"dsn"`
}

type HTTPClient interface {
	PostWithTimeout(URL string, payload []byte, timeout time.Duration) (*http.Response, error)
}

func HandleTunnel(client HTTPClient, logger *slog.Logger, config *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With("op", "handlers/HandleTunnel")
		host := r.Header.Get("Origin")

		for _, h := range config.AllowOrigins {
			if h == host || h == "*" {
				w.Header().Set("Access-Control-Allow-Origin", h)
			}
		}

		envelope, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("could not read request body:", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		piece := strings.Split(string(envelope), "\n")[0]
		var result payload
		decoder := json.NewDecoder(strings.NewReader(piece))
		err = decoder.Decode(&result)
		if err != nil {
			log.Error("could not decode request body:", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if result.DSN == "" {
			log.Error("DSN is empty")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = url.ValidateDSN(result.DSN, config.DSN)
		if err != nil {
			log.Error("invalid DSN provided in the payload:", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		upstreamURL, err := url.BuildSentryUpstreamURL(result.DSN)
		if err != nil {
			log.Error("could not build upstream URL:", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp, err := client.PostWithTimeout(upstreamURL, envelope, config.RequestTimeout)
		if err != nil {
			log.Error("could not send request to upstream:", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error("upstream returned non-200 status code:", slog.Int("status_code", resp.StatusCode))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
