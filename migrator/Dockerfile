FROM golang:rc-stretch as builder
ENV GO111MODULE=on

WORKDIR /src
COPY go.mod go.sum ./
RUN go get

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /validator ./cmd/validator && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /importer ./cmd/importer


FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .


FROM scratch
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
# the tls certificates:
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /validator /validator
COPY --from=builder /importer /importer

VOLUME ["/data", "/records", "/splitted"]

CMD ["/importer"]
