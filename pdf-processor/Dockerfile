FROM golang:alpine as builder

RUN apk --no-cache --virtual add openssl imagemagick imagemagick-dev poppler-utils git build-base && \
    wget -O /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && chmod +x /usr/local/bin/dep

WORKDIR /go/src/github.com/dgmann/document-manager/pdf-processor
COPY Gopkg.toml Gopkg.lock ./

RUN dep ensure -vendor-only
COPY . .
RUN go build -o /processor .

EXPOSE 8181

CMD ["/processor"]