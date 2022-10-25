FROM golang:1.18 as builder

ARG Version
WORKDIR /subsocks
COPY . .

ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w -X "''main.Version='${Version}'"'

FROM alpine:latest

WORKDIR /subsocks
COPY --from=builder /subsocks/subsocks .

ENTRYPOINT ["/subsocks/subsocks"]
