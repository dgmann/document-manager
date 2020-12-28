FROM alpine:latest as base

RUN apk add --no-cache \
    libjpeg-turbo libpng tiff poppler-utils

FROM golang:1.15 as builder

RUN apt-get update && \
    apt-get install -q -y wget build-essential pkg-config poppler-utils git g++ libjpeg-dev libpng-dev libtiff-dev libgif-dev ghostscript libgs-dev --no-install-recommends && \
	rm -rf /var/lib/apt/lists/*

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /processor .

FROM base
EXPOSE 8080 9000

COPY --from=builder /processor /processor

CMD ["/processor"]