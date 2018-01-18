FROM golang as builder

WORKDIR /go/src/github.com/dgmann/pdf-processor
COPY . .

RUN go-wrapper download && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /processor .

FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .

FROM scratch
COPY --from=builder /processor /processor
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
# the tls certificates:
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 80

CMD ["/processor"]