FROM golang:1.15 as builder

RUN apt-get update && \
    apt-get install -q -y wget build-essential pkg-config poppler-utils git g++ libjpeg-dev libpng-dev libtiff-dev libgif-dev ghostscript libgs-dev --no-install-recommends && \
	rm -rf /var/lib/apt/lists/*

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /processor .

EXPOSE 8181

CMD ["/processor"]