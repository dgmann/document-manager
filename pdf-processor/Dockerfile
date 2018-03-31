FROM golang:alpine as builder

WORKDIR /go/src/github.com/dgmann/pdf-processor
COPY . .

RUN apk --no-cache --virtual add openssl imagemagick-dev poppler-utils git build-base
RUN wget -O /bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && chmod +x /bin/dep && \
    dep ensure && go build -o /processor .

EXPOSE 80

CMD ["/processor"]