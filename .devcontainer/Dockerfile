ARG VARIANT=1.21-bullseye
FROM mcr.microsoft.com/devcontainers/go:${VARIANT}

RUN apt update && apt-get install -y -qq \
  # OCR
  libtesseract-dev \
  libleptonica-dev \
  tesseract-ocr-eng \
  tesseract-ocr-deu \
  # PDF Processor
  wget \
  build-essential \
  pkg-config \
  poppler-utils \
  git \
  g++ \
  libjpeg-dev \
  libpng-dev \
  libtiff-dev \
  libgif-dev \
  ghostscript \
  libgs-dev
