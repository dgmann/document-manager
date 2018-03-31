FROM golang as builder

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 && chmod +x /usr/local/bin/dep
WORKDIR /go/src/github.com/dgmann/document-manager-api
COPY Gopkg.toml Gopkg.lock ./

RUN dep ensure -vendor-only
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /api .


FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .


FROM scratch
COPY --from=builder /api /api
COPY --from=builder /go/src/github.com/dgmann/document-manager-api/migrations /migrations
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
# the tls certificates:
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 80

ENV DB_HOST=localhost DB_USER=postgres DB_PASSWORD=postgres DB_PORT=5432 DB_NAME=manager PORT=80

CMD ["/api"]