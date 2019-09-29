package main

import (
	"fmt"
	"github.com/dgmann/document-manager/pdf-processor/gopdf"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func main() {
	creator := gopdf.NewPdfCreator()
	port := 9000
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(grpc.MaxRecvMsgSize(1024*1024*300), grpc.MaxSendMsgSize(1024*1024*300))
	processor.RegisterPdfProcessorServer(grpcServer, NewGRPCServer(nil, nil, creator))
	if err := grpcServer.Serve(lis); err != nil {
		log.Error(err)
	}
}
