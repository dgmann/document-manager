DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
protoc --proto_path=$DIR/../api \
  --go_out=$DIR/.. --go_opt=module=github.com/dgmann/document-manager/pdf-processor \
  --go-grpc_out=$DIR/.. --go-grpc_opt=module=github.com/dgmann/document-manager/pdf-processor \
  $DIR/../api/processor.proto
