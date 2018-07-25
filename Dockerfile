FROM golang:alpine as builder
WORKDIR $GOPATH/src/github.com/kostiamol/webhook

COPY cmd cmd
COPY vendor vendor
RUN CGO_ENABLED=0 GOOS=linux go install -ldflags="-w -s" -v github.com/kostiamol/webhook/cmd/webhook

FROM alpine:latest
COPY --from=builder /go/bin/webhook /bin/webhook