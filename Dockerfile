FROM golang:alpine as builder
WORKDIR $GOPATH/src/github.com/kostiamol/k8s-webhook

COPY cmd cmd
COPY vendor vendor
RUN CGO_ENABLED=0 GOOS=linux go install -ldflags="-w -s" -v github.com/kostiamol/k8s-webhook/cmd/k8s-webhook

FROM alpine:latest
COPY --from=builder /go/bin/k8s-webhook /bin/k8s-webhook