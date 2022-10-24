FROM golang:1.18 as builder

ARG VERSION
WORKDIR /subsocks
COPY . .

ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct
RUN go build -ldflags "-X main.Version=${VERSION}"

FROM alpine:latest

WORKDIR /subsocks
COPY --from=builder /subsocks/subsocks .

ENTRYPOINT ["/subsocks/subsocks"]
