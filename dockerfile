FROM golang:1.11-alpine

EXPOSE 8080

COPY . /go/src/go-proxy-server/

WORKDIR /go/src/go-proxy-server/src/

RUN chmod +x * -R && go get -d -v ./... && go install -v ./...

CMD ["/bin/sh"]