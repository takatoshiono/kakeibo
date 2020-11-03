FROM golang:1.15
WORKDIR /go/src/github.com/takatoshiono/kakeibo/backend
COPY ./ .
RUN mkdir -p build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/mf ./cmd/mf/main.go

FROM ubuntu:18.04
WORKDIR /app
COPY --from=0 /go/src/github.com/takatoshiono/kakeibo/backend/build/mf .
USER 10000:10000
