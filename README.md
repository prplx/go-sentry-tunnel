# Go Sentry Tunnel

A simple Sentry tunnel written in Go. Serves as a proxy for tunneling Sentry browser requests into the Sentry server. Useful when requests to Sentry are blocked by browsers ad blockers. See [official Sentry documentation](https://docs.sentry.io/platforms/javascript/troubleshooting/#using-the-tunnel-option).

## Configuration

The application looks up for the next environment variables:

`DSN` - comma separated list of Sentry DSN's to work with. Required.
`ALLOW_ORIGINS` - comma separated list of origins to bypass CORS browser check. Optional. Default value is `*`.
`PORT` - port this app will be running on. Optional. Default value is `3001`.

## Running with docker

```sh
docker build -t go-sentry-tunnel .
docker run --rm -e 'DSN=https://user@host.ingest.sentry.io/project' -p 3001:3001 go-sentry-tunnel
```

## Running without docker

```sh
git clone https://github.com/prplx/go-sentry-tunnel.git
cd go-sentry-tunnel && go mod install
DSN=https://user@host.ingest.sentry.io/project go run cmd/api/main.go
```
