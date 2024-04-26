ARG GO_VERSION
FROM golang:${GO_VERSION} as builder
ARG SERVICE

RUN apt-get update && apt-get install -y unzip libaio1 libaio-dev

WORKDIR /tmp
COPY ./apps/${SERVICE}/docker/oci8.pc /usr/lib/pkgconfig
RUN wget -O instantclient.zip https://download.oracle.com/otn_software/linux/instantclient/193000/instantclient-basiclite-linux.x64-19.3.0.0.0dbru.zip \
    && unzip instantclient.zip && rm instantclient.zip \
    && wget -O include.zip https://download.oracle.com/otn_software/linux/instantclient/193000/instantclient-sdk-linux.x64-19.3.0.0.0dbru.zip \
    && unzip include.zip && rm include.zip \
    && mkdir -p /usr/local/oracle/lib/oracle && mv instantclient*/*.so* /usr/local/oracle/lib/oracle \
    && mkdir -p /usr/local/oracle/include/oracle && mv instantclient*/sdk/include/* /usr/local/oracle/include/oracle \
    && rm -r /tmp/*
ENV LD_LIBRARY_PATH=/usr/local/oracle/lib/oracle/

WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod,from=cache \
    --mount=type=bind,target=.,source=. \
    GOOS=linux go build -a -installsuffix cgo -o /${SERVICE} ./apps/${SERVICE}


FROM ubuntu:22.04
ARG SERVICE
RUN apt-get update && apt-get install -y libaio1
COPY --link --from=builder /usr/local/oracle/lib/oracle /usr/local/oracle/lib/oracle
ENV LD_LIBRARY_PATH=/usr/local/oracle/lib/oracle/

COPY --link --from=builder /${SERVICE} /${SERVICE}
COPY --link --chmod=+x ./apps/${SERVICE}/docker/entrypoint.sh /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/${SERVICE}"]
