package url

import (
	"testing"
)

func TestValidateDSN(t *testing.T) {
	type args struct {
		DSN       string
		validDSNs []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid DSN",
			args: args{
				DSN:       "https://example.com",
				validDSNs: []string{"https://example.com", "https://another.com"},
			},
			wantErr: false,
		},
		{
			name: "Invalid DSN",
			args: args{
				DSN:       "https://invalid.com",
				validDSNs: []string{"https://example.com", "https://another.com"},
			},
			wantErr: true,
		},
		{
			name: "Empty DSN",
			args: args{
				DSN:       "",
				validDSNs: []string{"https://example.com", "https://another.com"},
			},
			wantErr: true,
		},
		{
			name: "Empty validDSNs",
			args: args{
				DSN:       "https://example.com",
				validDSNs: []string{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateDSN(tt.args.DSN, tt.args.validDSNs); (err != nil) != tt.wantErr {
				t.Errorf("ValidateDSN() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildSentryUpstreamURL(t *testing.T) {
	type args struct {
		rawURL string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Valid URL with HTTP scheme",
			args: args{
				rawURL: "http://example.com/12345",
			},
			want:    "https://example.com/api/12345/envelope/?sentry_key=",
			wantErr: false,
		},
		{
			name: "Valid URL with HTTPS scheme",
			args: args{
				rawURL: "https://example.com/12345",
			},
			want:    "https://example.com/api/12345/envelope/?sentry_key=",
			wantErr: false,
		},
		{
			name: "Valid URL with user info",
			args: args{
				rawURL: "https://user:pass@example.com/12345",
			},
			want:    "https://example.com/api/12345/envelope/?sentry_key=user:pass",
			wantErr: false,
		},
		{
			name: "Invalid URL",
			args: args{
				rawURL: "://invalid-url",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "URL with query parameters",
			args: args{
				rawURL: "https://example.com/12345?query=123",
			},
			want:    "https://example.com/api/12345/envelope/?sentry_key=",
			wantErr: false,
		},
		{
			name: "URL with fragment",
			args: args{
				rawURL: "https://example.com/12345#fragment",
			},
			want:    "https://example.com/api/12345/envelope/?sentry_key=",
			wantErr: false,
		},
		{
			name: "URL with port",
			args: args{
				rawURL: "https://example.com:8080/12345",
			},
			want:    "https://example.com:8080/api/12345/envelope/?sentry_key=",
			wantErr: false,
		},
		{
			name: "URL with complex path",
			args: args{
				rawURL: "https://example.com/path/to/project/12345",
			},
			want:    "https://example.com/api/path/to/project/12345/envelope/?sentry_key=",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildSentryUpstreamURL(tt.args.rawURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildSentryUpstreamURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BuildSentryUpstreamURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
