 # syntax=docker/dockerfile:1.10
ARG GO_VERSION=1.23
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS builder
ARG SERVICE
WORKDIR /src
RUN --mount=type=bind,target=go.mod,source=go.mod \
    --mount=type=bind,target=go.sum,source=go.sum \ 
    go mod download
RUN --mount=type=bind,target=go.mod,source=go.mod \
    --mount=type=bind,target=go.sum,source=go.sum \
    --mount=type=bind,target=cmd,source=cmd \
    --mount=type=bind,target=pkg,source=pkg \
    --mount=type=bind,target=internal,source=internal \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /${SERVICE} ./cmd/${SERVICE}


FROM alpine:latest AS alpine
RUN apk --no-cache add tzdata zip ca-certificates && \
    zip -r -0 /zoneinfo.zip /usr/share/zoneinfo

FROM --platform=$BUILDPLATFORM scratch AS raw
ARG SERVICE
ARG SERVICE_NAME=${SERVICE}
ARG SERVICE_VERSION=latest
ARG TARGETOS
ARG TARGETARCH
COPY --from=builder "/${SERVICE}" "/${SERVICE_NAME}_${SERVICE_VERSION}_${TARGETOS}_${TARGETARCH}"

FROM --platform=$BUILDPLATFORM scratch
ARG SERVICE
ENV ZONEINFO=/zoneinfo.zip
COPY --link --from=alpine /zoneinfo.zip /
# the tls certificates:
COPY --link --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --link --from=builder /${SERVICE} /service

ENV HTTP_PORT=8080
EXPOSE 8080

CMD ["/service"]
