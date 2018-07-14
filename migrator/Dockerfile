FROM golang as builder

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && chmod +x /usr/local/bin/dep
WORKDIR /go/src/github.com/dgmann/document-manager/migrator
COPY Gopkg.toml Gopkg.lock ./

RUN dep ensure -vendor-only
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /migrator .


FROM openjdk:jre-alpine
RUN wget -O SplitPDF.jar "https://sourceforge.net/projects/splitpdf/files/latest/download"
COPY --from=builder /migrator /migrator

ENV DIRECTORY=/records DESTINATION=http://api PARSER=generic SENDER=Scan RETRY=5

VOLUME ["/records"]

CMD ["/migrator"]
