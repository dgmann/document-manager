package main

import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"

	"github.com/dgmann/document-manager/pdf-processor/pkg/image/imaging"
	"github.com/dgmann/document-manager/pdf-processor/pkg/pdf"
	"github.com/dgmann/document-manager/pdf-processor/pkg/pdf/dual"
	"github.com/dgmann/document-manager/pdf-processor/pkg/pdf/gopdf"
	"github.com/dgmann/document-manager/pdf-processor/pkg/pdf/mupdf"
	"github.com/dgmann/document-manager/pdf-processor/pkg/pdf/pdfcpu"
	"github.com/dgmann/document-manager/pdf-processor/pkg/pdf/poppler"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe(":8080", nil))
	}()
	config := ConfigFromEnv()

	extractors := make(map[string]pdf.ImageConverter)
	extractors["poppler"] = poppler.NewExtractor()
	extractors["pdfcpu"] = pdfcpu.NewExtractor()

	rasterizers := make(map[string]pdf.ImageConverter)
	rasterizers["poppler"] = poppler.NewRasterizer()
	rasterizers["mupdf"] = mupdf.NewProcessor()

	rotator := imaging.NewRotator()
	extractor, ok := extractors[config.Extractor]
	if !ok {
		log.Fatalf("%s is not a valid extractor. Valid values: poppler, pdfcpu", config.Extractor)
	}
	rasterizer, ok := rasterizers[config.Rasterizer]
	if !ok {
		log.Fatalf("%s is not a valid rasterizer. Valid values: poppler, mupdf", config.Extractor)
	}
	converter := dual.NewProcessor(extractor, rasterizer, mupdf.NewProcessor())

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
