ARG GO_VERSION
FROM golang:${GO_VERSION} as builder
ARG SERVICE
RUN apt-get update && \
    apt-get install -q -y wget build-essential pkg-config poppler-utils git g++ libjpeg-dev libpng-dev libtiff-dev libgif-dev ghostscript libgs-dev --no-install-recommends && \
	rm -rf /var/lib/apt/lists/*
WORKDIR /src
RUN --mount=type=bind,target=go.mod,source=go.mod \
    --mount=type=bind,target=go.sum,source=go.sum \
    --mount=type=bind,target=cmd,source=cmd \
    --mount=type=bind,target=pkg,source=pkg \
    --mount=type=bind,target=internal,source=internal \
    GOOS=linux CGO_ENABLED=0 go build -o /pdf-processor ./cmd/${SERVICE}


FROM alpine:latest
ARG SERVICE
RUN apk add --no-cache \
    libjpeg-turbo libpng tiff poppler-utils msttcorefonts-installer fontconfig && \
    update-ms-fonts && \
    fc-cache -f

COPY --link --from=builder /pdf-processor /pdf-processor
EXPOSE 8080 9000

CMD ["/pdf-processor"]
