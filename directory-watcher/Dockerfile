FROM golang:rc-stretch as builder
ENV GO111MODULE=on

WORKDIR /src
COPY go.mod go.sum ./
RUN go get

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /watcher .


FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .


FROM scratch
COPY --from=builder /watcher /watcher
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
# the tls certificates:
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV DIRECTORY=/records DESTINATION=http://api PARSER=generic SENDER=Scan RETRY=5

VOLUME ["/records"]

ENTRYPOINT ["/watcher"]