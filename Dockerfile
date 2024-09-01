FROM golang:alpine as builder

ARG ARCH="arm64"

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

ENV USER=appuser
ENV UID=10001

RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"
    
WORKDIR $GOPATH/src/go-sentry-tunnel

COPY . .

RUN go mod download && go mod verify
RUN GOOS=linux GOARCH=${ARCH} go build -ldflags="-w -s" -o /go/bin/go-sentry-tunnel cmd/api/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/bin/go-sentry-tunnel /go/bin/go-sentry-tunnel

USER appuser:appuser

ENTRYPOINT ["/go/bin/go-sentry-tunnel"]
