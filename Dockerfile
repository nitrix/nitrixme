FROM golang:alpine AS builder

WORKDIR /opt
COPY . /opt

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o nitrixme

FROM scratch AS release
COPY --from=builder /opt /opt

WORKDIR /opt
CMD ["/opt/nitrixme"]
