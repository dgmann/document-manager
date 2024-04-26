ARG GO_VERSION
FROM golang:${GO_VERSION} as builder
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
RUN --mount=type=cache,target=/go/pkg/mod,from=cache \
    --mount=type=bind,target=.,source=. \
    GOOS=linux go build -a -installsuffix cgo -o /${SERVICE} ./apps/${SERVICE}


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
