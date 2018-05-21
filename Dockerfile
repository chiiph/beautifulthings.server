FROM golang:latest as builder

WORKDIR /go/src/beautifulthings/
COPY ./ .
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && \
    chmod +x /usr/local/bin/dep
RUN dep ensure -vendor-only
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o serv ./cmd/server

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go/src/beautifulthings/serv .
EXPOSE 8080/tcp
CMD ["./serv"]
