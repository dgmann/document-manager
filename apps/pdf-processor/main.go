package main

import (
	"context"
	"fmt"
	"github.com/dgmann/document-manager/pkg/opentelemetry"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/dgmann/document-manager/pdf-processor/pkg/pdf/unipdf"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"

	"github.com/dgmann/document-manager/pdf-processor/pkg/image/imaging"
	"github.com/dgmann/document-manager/pdf-processor/pkg/pdf"
	"github.com/dgmann/document-manager/pdf-processor/pkg/pdf/gopdf"
	"github.com/dgmann/document-manager/pdf-processor/pkg/pdf/pdfcpu"
	"github.com/dgmann/document-manager/pdf-processor/pkg/pdf/poppler"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	formatter := &log.TextFormatter{}
	formatter.FullTimestamp = true
	formatter.TimestampFormat = time.RFC3339Nano
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
		log.WarnLevel,
		log.InfoLevel,
	)))
}

func main() {
	go func() {
		log.Println(http.ListenAndServe(":8080", nil))
	}()
	config := ConfigFromEnv()
	ctx := context.Background()
	otlProvider, err := opentelemetry.NewProvider(ctx, "pdf-processor", config.OtelCollectorUrl)
	if err != nil {
		log.WithError(err).Warnln("error creating OpenTelemetry exporter")
	}
	defer func(otlProvider *opentelemetry.Provider, ctx context.Context) {
		if otlProvider != nil {
			_ = otlProvider.Shutdown(ctx)
		}
	}(otlProvider, ctx)
	if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
		log.WithContext(ctx).WithError(err).Warnln("error initializing runtime metrics")
	}

	extractors, rasterizers := initProcessors()

	rotator := imaging.NewRotator()
	extractor, ok := extractors[config.Extractor]
	if !ok {
		log.Fatalf("%s is not a valid extractor. Valid values: poppler, pdfcpu", config.Extractor)
	}
	rasterizer, ok := rasterizers[config.Rasterizer]
	if !ok {
		log.Fatalf("%s is not a valid rasterizer. Valid values: poppler, mupdf", config.Extractor)
	}
	converter := pdf.NewConverter(extractor(), rasterizer(), pdfcpu.NewExtractor())

	creator := gopdf.NewPdfCreator()
	port := 9000
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to open socket: %v", err)
	}
	log.Info("Starting gRPC Server")
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(1024*1024*300),
		grpc.MaxSendMsgSize(1024*1024*300),
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	processor.RegisterPdfProcessorServer(grpcServer, NewGRPCServer(converter, rotator, creator))
	if err := grpcServer.Serve(lis); err != nil {
		log.Error(err)
	}
}

type initFunc func() pdf.ImageConverter

func initProcessors() (map[string]initFunc, map[string]initFunc) {
	extractors := make(map[string]initFunc)
	extractors["poppler"] = func() pdf.ImageConverter { return poppler.NewExtractor() }
	extractors["pdfcpu"] = func() pdf.ImageConverter { return pdfcpu.NewExtractor() }
	extractors["unipdf"] = func() pdf.ImageConverter { return unipdf.NewExtractor() }

	rasterizers := make(map[string]initFunc)
	rasterizers["poppler"] = func() pdf.ImageConverter { return poppler.NewRasterizer() }
	rasterizers["unipdf"] = func() pdf.ImageConverter { return unipdf.NewRasterizer() }
	return extractors, rasterizers
}
