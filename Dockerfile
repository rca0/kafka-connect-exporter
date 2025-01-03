ARG BASE_IMAGE=golang:1.22.1
FROM --platform=$BUILDPLATFORM ${BASE_IMAGE}

ARG TARGETOS TARGETARCH 

ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8
ENV GOPATH /go
ENV CGO_ENABLED=0
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
ENV GO111MODULE=on  

WORKDIR $GOPATH 

ARG GOPROXY
ARG GOPRIVATE 

COPY .version .version 
COPY go.mod go.mod 
COPY go.sum go.sum  
 
RUN go mod download 

COPY main.go main.go 

COPY main.go main.go  
COPY internal/ internal/ 

RUN go build -a -o kafka-connect-exporter main.go 
RUN useradd -u 1001 kafka-connect-exporter 
USER kafka-connect-exporter

CMD ["/go/kafka-connect-exporter"]
