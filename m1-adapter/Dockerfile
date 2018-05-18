FROM sergeymakinen/oracle-instant-client as builder

RUN apt-get update && apt-get install -y git curl build-essential pkg-config

RUN curl -O https://dl.google.com/go/go1.10.2.linux-amd64.tar.gz && \
    tar -xvf go1.10.2.linux-amd64.tar.gz && \
    mv go /usr/local && \
    export PATH="/usr/local/go/bin:$PATH"

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && chmod +x /usr/local/bin/dep
ENV GOPATH=/go PATH="/usr/local/go/bin:$PATH"
COPY .docker/oci8.pc /usr/lib/pkgconfig
WORKDIR /go/src/github.com/dgmann/document-manager/m1-adapter
COPY Gopkg.toml Gopkg.lock ./

RUN go version
RUN dep ensure -vendor-only
COPY . .
RUN GOOS=linux go build -a -installsuffix cgo -o /m1-adapter .


FROM sergeymakinen/oracle-instant-client
COPY --from=builder /m1-adapter /m1-adapter

ENV DSN="m1/a7haj9@192.168.255.103:1521:M1DB"

CMD ["/m1-adapter"]
