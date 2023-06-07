FROM golang:1.19 AS build-env
RUN apk update && apk add git make
COPY . /go/src/github.com/capitalonline/eks-cloud-controller-manager
RUN cd /go/src/github.com/capitalonline/eks-cloud-controller-manager && go build main.go -o ccm

FROM alpine:3.6
RUN apk update --no-cache && apk add ca-certificates
COPY --from=build-env /ccm /ccm

ENTRYPOINT ["ccm"]
