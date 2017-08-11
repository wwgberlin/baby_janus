FROM golang:latest
WORKDIR /go/src/github.com/wwgberlin/baby_janus/api
RUN cd /go/src/github.com/wwgberlin/baby_janus/api
CMD go test ./... && go build && /go/src/github.com/wwgberlin/baby_janus/api/api
EXPOSE 8080