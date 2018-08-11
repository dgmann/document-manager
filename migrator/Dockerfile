FROM golang:rc-stretch as builder
ENV GO111MODULE=on

WORKDIR /src
COPY go.mod go.sum ./
RUN go get

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /validator ./cmd/validator && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /importer ./cmd/importer


FROM openjdk:jre-alpine
RUN wget -O SplitPDF.jar "https://sourceforge.net/projects/splitpdf/files/latest/download"

COPY --from=builder /validator /validator
COPY --from=builder /importer /importer

VOLUME ["/data", "/records", "/splitted"]

CMD ["/validator"]
