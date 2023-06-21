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


FROM debian:bookworm

RUN apt-get update && apt-get install -y -qq \
  libtesseract-dev \
  libleptonica-dev \
  tesseract-ocr-eng \
  tesseract-ocr-deu

COPY --from=builder /ocr-service /app/ocr-service

CMD ["/app/ocr-service"]
