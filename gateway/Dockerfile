FROM golang:latest
WORKDIR /go/src/github.com/wwgberlin/baby_janus/gateway
RUN cd /go/src/github.com/wwgberlin/baby_janus/gateway
CMD go test ./... && go run /go/src/github.com/wwgberlin/baby_janus/gateway/cmd/gateway/main.go
EXPOSE 8080