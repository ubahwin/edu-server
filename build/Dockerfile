FROM golang:1.23.3-alpine3.20 AS builder

ENV SERVICE=edu-server

ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /go/src/$SERVICE

RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest

RUN apk --update --no-cache add ca-certificates gcc libtool make musl-dev protoc git

COPY . /go/src/$SERVICE/
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go mod download
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod CGO_ENABLED=0 go build -gcflags "all=-N -l" -o $SERVICE cmd/app/*.go

FROM scratch

ENV SERVICE=edu-server
COPY --from=builder /go/bin/dlv /dlv
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/$SERVICE/$SERVICE /$SERVICE
EXPOSE 8080 2345 7070
ENTRYPOINT ["/dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./edu-server"]
CMD []
