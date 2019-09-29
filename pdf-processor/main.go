package main

import (
	"fmt"
	"github.com/dgmann/document-manager/pdf-processor/imagick"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	"github.com/dgmann/document-manager/pdf-processor/poppler"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	rotator := imagick.NewProcessor()
	converter := poppler.NewProcessor()
	port := 9000
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(grpc.MaxRecvMsgSize(1024*1024*300), grpc.MaxSendMsgSize(1024*1024*300))
	processor.RegisterPdfProcessorServer(grpcServer, NewGRPCServer(converter, rotator))
	if err := grpcServer.Serve(lis); err != nil {
		log.Error(err)
	}
}
