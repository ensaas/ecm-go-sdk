FROM golang:1.13.9-alpine3.11

WORKDIR $GOPATH/src/ecm-sdk-go

COPY  ./ ./

RUN CGO_ENABLED=0 go build -ldflags '-w -s' -o demo ./example/main.go

FROM alpine:3.11.5

WORKDIR /app

COPY --from=0 /go/src/ecm-sdk-go/demo ./

ENTRYPOINT ["./demo"]
