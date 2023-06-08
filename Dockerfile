FROM golang:1.20 AS build-env
COPY . /go/src/app/
RUN  go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.io,direct
RUN cd /go/src/app
WORKDIR /go/src/app
RUN go build -o  ccm ./cmd/main.go
#CMD sleep 3000
ENTRYPOINT ["./ccm"]