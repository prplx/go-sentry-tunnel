package url

import (
	"fmt"
	"net/url"
	"strings"

	"errors"
)

func ValidateDSN(DSN string, validDSNs []string) error {
	for _, dsn := range validDSNs {
		if dsn == DSN {
			return nil
		}
	}

	return errors.New("invalid DSN")
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
