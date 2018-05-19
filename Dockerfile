FROM golang:latest as builder

WORKDIR /go/src/beautifulthings/
COPY ./ .
RUN go get -v ./...
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o serv ./cmd/server
#CMD ["./serv"]

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go/src/beautifulthings/serv .
CMD ["./serv"]