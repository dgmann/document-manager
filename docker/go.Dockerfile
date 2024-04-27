 # syntax=docker/dockerfile:1
ARG GO_VERSION
FROM golang:${GO_VERSION} as builder
ARG SERVICE
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod,from=cache \
    --mount=type=bind,target=.,source=. \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /${SERVICE} ./apps/${SERVICE}


FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates && \
    zip -r -0 /zoneinfo.zip /usr/share/zoneinfo


FROM scratch
ARG SERVICE
ENV ZONEINFO /zoneinfo.zip
COPY --link --from=alpine /zoneinfo.zip /
# the tls certificates:
COPY --link --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --link --from=builder /${SERVICE} /${SERVICE}

EXPOSE 8080

CMD ["/${SERVICE}"]
