package handlers

import (
	"log/slog"
	"net/http"
	"reflect"
	"testing"

	"go-sentry-tunnel/internal/config"
)

func TestHandleTunnel(t *testing.T) {
	type args struct {
		l *slog.Logger
		c *config.Config
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HandleTunnel(tt.args.l, tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleTunnel() = %v, want %v", got, tt.want)
			}
		})
	}
}
