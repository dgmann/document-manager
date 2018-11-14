FROM sergeymakinen/oracle-instant-client as builder

RUN apt-get update && apt-get install -y git curl build-essential pkg-config

RUN curl -O https://dl.google.com/go/go1.11.2.linux-amd64.tar.gz && \
    tar -xvf go1.11.2.linux-amd64.tar.gz && \
    mv go /usr/local && \
    export PATH="/usr/local/go/bin:$PATH"

ENV GOPATH=/go PATH="/usr/local/go/bin:$PATH"
COPY .docker/oci8.pc /usr/lib/pkgconfig

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux go build -a -installsuffix cgo -o /m1-adapter .


FROM sergeymakinen/oracle-instant-client
COPY .docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
COPY --from=builder /m1-adapter /m1-adapter

ENV PORT=80 DB_USERNAME="" DB_PASSWORD="" DB_HOST="" DB_PORT=1521 DB_NAME=M1DB

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/m1-adapter"]
