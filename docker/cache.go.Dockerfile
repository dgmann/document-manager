# syntax=docker/dockerfile:1.7-labs
ARG GO_VERSION
FROM golang:${GO_VERSION} as dev
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
