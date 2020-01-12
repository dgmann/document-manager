FROM golang:1.13 as dev

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .


FROM dev as builder
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /api .


FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .


FROM scratch
COPY --from=builder /api /api
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
# the tls certificates:
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 80

ENV DB_HOST=localhost DB_NAME=manager PORT=80 RECORD_DIR=/records

VOLUME ["/records", "/archive"]

CMD ["/api"]