FROM golang:alpine AS builder

WORKDIR /opt
COPY . /opt

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o nitrixme

FROM golang:alpine AS release

RUN apk add --no-cache ca-certificates
COPY --from=builder /opt /opt

WORKDIR /opt
ENV GIN_MODE=release
CMD ["/opt/nitrixme"]
