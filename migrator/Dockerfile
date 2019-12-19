FROM golang:latest as builder
ARG GITHUB_TOKEN

RUN git config --global url."https:///$GITHUB_TOKEN:x-oauth-basic@github.com/".insteadOf "https://github.com/"

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /migrator .


FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .


FROM scratch
WORKDIR /
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
# the tls certificates:
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /migrator /
COPY --from=builder /src/web /web

VOLUME ["/data", "/records"]

ENV DB_HOST=localhost \
    DB_USER=sa \
    DB_PASSWORD=Password% \
    DB_NAME=PdfDatabase \
    DB_INSTANCE="" \
    RECORD_DIR=/records \
    DATA_DIR=/data \
    RETRY_COUNTER=3 \
    API_URL=http://api/api

CMD ["/migrator"]
