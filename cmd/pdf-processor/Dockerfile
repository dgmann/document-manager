ARG GO_VERSION
FROM golang:${GO_VERSION} as builder
ARG SERVICE
RUN apt-get update && \
    apt-get install -q -y wget build-essential pkg-config poppler-utils git g++ libjpeg-dev libpng-dev libtiff-dev libgif-dev ghostscript libgs-dev --no-install-recommends && \
	rm -rf /var/lib/apt/lists/*
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod,from=cache \
    --mount=type=bind,target=.,source=. \
    CGO_ENABLED=0 GOOS=linux CGO_ENABLED=0 go build -o /${SERVICE} ./apps/${SERVICE}


FROM alpine:latest
ARG SERVICE
RUN apk add --no-cache \
    libjpeg-turbo libpng tiff poppler-utils msttcorefonts-installer fontconfig && \
    update-ms-fonts && \
    fc-cache -f

COPY --link --from=builder /${SERVICE} /${SERVICE}
EXPOSE 8080 9000

CMD ["/${SERVICE}"]
