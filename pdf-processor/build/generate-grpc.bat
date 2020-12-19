protoc --proto_path=%~dp0..\api^
 --go_out=%~dp0.. --go_opt=module=github.com/dgmann/document-manager/pdf-processor^
 --go-grpc_out=%~dp0.. --go-grpc_opt=module=github.com/dgmann/document-manager/pdf-processor^
 %~dp0..\api\processor.proto
