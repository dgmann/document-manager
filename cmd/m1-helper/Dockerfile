ARG GO_VERSION
ARG OUTPUT="/out"
FROM golang:${GO_VERSION} AS builder
ARG SERVICE
ARG OUTPUT
RUN apt-get update && apt-get install -y zip
WORKDIR /src
RUN --mount=type=bind,target=go.mod,source=go.mod \
    --mount=type=bind,target=go.sum,source=go.sum \
    --mount=type=bind,target=cmd,source=cmd \
    --mount=type=bind,target=pkg,source=pkg \
    --mount=type=bind,target=internal,source=internal <<EOF
set -e

mkdir "$OUTPUT"
cd cmd/${SERVICE}
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H=windowsgui" -o "$OUTPUT"/m1-helper.exe .
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o "$OUTPUT"/m1-helper_console.exe .
zip --junk-paths "$OUTPUT"/m1-helper-64.zip "$OUTPUT"/*.exe
rm "$OUTPUT"/*.exe
GOOS=windows GOARCH=386 go build -ldflags="-s -w -H=windowsgui" -o "$OUTPUT"/m1-helper.exe .
GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o "$OUTPUT"/m1-helper_console.exe .
zip --junk-paths "$OUTPUT"/m1-helper-32.zip "$OUTPUT"/*.exe
rm "$OUTPUT"/*.exe
EOF


FROM scratch
ARG OUTPUT
COPY --from=builder "${OUTPUT}"/* .
