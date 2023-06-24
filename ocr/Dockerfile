FROM golang:latest as builder

RUN apt-get update && apt-get install -y -qq \
  libtesseract-dev \
  libleptonica-dev \
  tesseract-ocr-eng \
  tesseract-ocr-deu
# See https://github.com/tesseract-ocr/tessdata for the list of available languages.
# If you want to download these traineddata via `wget`, don't forget to locate
# downloaded traineddata under ${TESSDATA_PREFIX}/tessdata.

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux go build -a -installsuffix cgo -o /ocr-service .


FROM ubuntu:22.04

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

COPY --from=builder /ocr-service /app/ocr-service

EXPOSE 8080

CMD ["/app/ocr-service"]
