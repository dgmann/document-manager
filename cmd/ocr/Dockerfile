ARG GO_VERSION=1.22
FROM golang:${GO_VERSION} AS builder
ARG SERVICE

RUN apt-get update && apt-get install -y -qq \
  libtesseract-dev \
  libleptonica-dev \
  tesseract-ocr-eng \
  tesseract-ocr-deu
# See https://github.com/tesseract-ocr/tessdata for the list of available languages.
# If you want to download these traineddata via `wget`, don't forget to locate
# downloaded traineddata under ${TESSDATA_PREFIX}/tessdata.

WORKDIR /src
RUN --mount=type=bind,target=go.mod,source=go.mod \
  --mount=type=bind,target=go.sum,source=go.sum \
  --mount=type=bind,target=cmd,source=cmd \
  --mount=type=bind,target=pkg,source=pkg \
  --mount=type=bind,target=internal,source=internal \
    GOOS=linux go build -a -installsuffix cgo -o /${SERVICE} ./cmd/${SERVICE}


FROM ubuntu:22.04
ARG SERVICE
RUN apt-get update && apt-get install -y -qq lsb-release wget gnupg && \
    echo "deb https://notesalexp.org/tesseract-ocr5/$(lsb_release -cs)/ $(lsb_release -cs) main" \
    | tee /etc/apt/sources.list.d/notesalexp.list > /dev/null && \
    wget -O - https://notesalexp.org/debian/alexp_key.asc | apt-key add - && \
    apt-get update && \
    apt-get install -y -qq \
    libtesseract-dev \
    libleptonica-dev \
    tesseract-ocr-eng-best \
    tesseract-ocr-deu-best

COPY --link --from=builder /${SERVICE} /${SERVICE}
EXPOSE 8080

CMD ["/${SERVICE}"]
