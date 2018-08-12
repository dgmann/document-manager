FROM golang:rc-stretch as builder
ENV GO111MODULE=on

RUN apt-get update && \
    apt-get install -q -y wget build-essential pkg-config poppler-utils git g++ libjpeg-dev libpng-dev libtiff-dev libgif-dev --no-install-recommends && \
    wget http://downloads.webmproject.org/releases/webp/libwebp-0.4.2.tar.gz && \
	tar xvzf libwebp-0.4.2.tar.gz && rm libwebp-0.4.2.tar.gz && \
	cd libwebp-0.4.2 && \
	./configure && \
	make && make install && \
	cd && rm -rf libwebp-0.4.2 && \
	wget http://www.imagemagick.org/download/ImageMagick.tar.gz && \
	tar xvzf ImageMagick.tar.gz && rm ImageMagick.tar.gz && \
	cd ImageMagick* && \
	./configure && \
	make && make install && \
	cd && rm -rf ImageMagick* && \
	ldconfig /usr/local/lib && \
	rm -rf /var/lib/apt/lists/*

WORKDIR /src
COPY go.mod go.sum ./
RUN go get

COPY . .
RUN go build -o /processor .

EXPOSE 8181

CMD ["/processor"]