FROM golang:1.14 as builder

RUN apt-get update && apt-get install -y unzip libaio1 libaio-dev

WORKDIR /tmp
COPY docker/oci8.pc /usr/lib/pkgconfig
RUN wget -O instantclient.zip https://download.oracle.com/otn_software/linux/instantclient/193000/instantclient-basiclite-linux.x64-19.3.0.0.0dbru.zip \
    && unzip instantclient.zip && rm instantclient.zip \
    && wget -O include.zip https://download.oracle.com/otn_software/linux/instantclient/193000/instantclient-sdk-linux.x64-19.3.0.0.0dbru.zip \
    && unzip include.zip && rm include.zip
RUN mkdir -p /usr/local/oracle/lib/oracle && mv instantclient*/*.so* /usr/local/oracle/lib/oracle \
    && mkdir -p /usr/local/oracle/include/oracle && mv instantclient*/sdk/include/* /usr/local/oracle/include/oracle

ENV LD_LIBRARY_PATH=/usr/local/oracle/lib/oracle/

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux go build -a -installsuffix cgo -o /m1-adapter .


FROM ubuntu
RUN apt-get update && apt-get install -y libaio1
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
COPY --from=builder /m1-adapter /m1-adapter
COPY --from=builder /usr/local/oracle/lib/oracle /usr/local/oracle/lib/oracle
ENV LD_LIBRARY_PATH=/usr/local/oracle/lib/oracle/
ENV PORT=80 DB_USERNAME="" DB_PASSWORD="" DB_HOST="" DB_PORT=1521 DB_NAME=M1DB

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/m1-adapter"]
