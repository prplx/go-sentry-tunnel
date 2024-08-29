package url

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/prplx/go-sentry-tunnel/internal/config"
	"github.com/prplx/go-sentry-tunnel/internal/errors"
)

func ValidateDSN(DSN string, config *config.Config) error {
	for _, dsn := range config.DSN {
		if dsn == DSN {
			return nil
		}
	}

	return errors.ErrorInvalidDSN
}

func BuildSentryUpstreamURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	projectID := strings.TrimPrefix(parsedURL.Path, "/")
	upstreamSentryURL := fmt.Sprintf("https://%s/api/%s/envelope/?sentry_key=%s", parsedURL.Host, projectID, parsedURL.User.String())

	return upstreamSentryURL, nil
}
