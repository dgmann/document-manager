package main

import (
	"fmt"
	"github.com/dgmann/document-manager/pdf-processor/gopdf"
	"github.com/dgmann/document-manager/pdf-processor/imagick"
	"github.com/dgmann/document-manager/pdf-processor/mupdf"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe(":8080", nil))
	}()

	imagick.Initialize()
	defer imagick.Terminate()

	rotator := imagick.NewProcessor()
	converter := mupdf.NewProcessor()

	creator := gopdf.NewPdfCreator()
	port := 9000
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to open socket: %v", err)
	}
	log.Info("Starting gRPC Server")
	grpcServer := grpc.NewServer(grpc.MaxRecvMsgSize(1024*1024*300), grpc.MaxSendMsgSize(1024*1024*300))
	processor.RegisterPdfProcessorServer(grpcServer, NewGRPCServer(converter, rotator, creator))
	if err := grpcServer.Serve(lis); err != nil {
		log.Error(err)
	}
}
