FROM golang:1.11-alpine

RUN apk add git

COPY . /go/src/go-proxy-server/

WORKDIR /go/src/go-proxy-server/src/

RUN chmod +x * -R && go get -v

CMD ["/bin/sh"]