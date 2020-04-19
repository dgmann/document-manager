package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgmann/document-manager/api/datastore"
	"github.com/dgmann/document-manager/api/storage"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"io"
	"io/ioutil"
)

type PdfProcessor struct {
	baseUrl    string
	conn       *grpc.ClientConn
	images     storage.ImageService
	categories datastore.CategoryService
}

func NewPDFProcessor(baseUrl string, images storage.ImageService, cateogories datastore.CategoryService) (*PdfProcessor, error) {
	conn, err := grpc.Dial(
		baseUrl,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*300), grpc.MaxCallSendMsgSize(1024*1024*300)),
	)
	if err != nil {
		return nil, err
	}

	return &PdfProcessor{conn: conn, baseUrl: baseUrl, images: images, categories: cateogories}, nil
}

func (p *PdfProcessor) Close() error {
	return p.conn.Close()
}

func (p *PdfProcessor) Convert(ctx context.Context, f io.Reader) ([]storage.Image, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	client := processor.NewPdfProcessorClient(p.conn)
	stream, err := client.ConvertPdfToImage(ctx, &processor.Pdf{Content: b})
	if err != nil {
		return nil, err
	}

	var images []storage.Image
	for {
		image, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		images = append(images, *storage.NewImage(image.Content, image.Format))
	}
	return images, nil
}

func (p *PdfProcessor) Rotate(ctx context.Context, image io.Reader, degrees int) (*storage.Image, error) {
	b, err := ioutil.ReadAll(image)
	if err != nil {
		return nil, err
	}

	client := processor.NewPdfProcessorClient(p.conn)
	result, err := client.RotateImage(ctx, &processor.Rotate{Content: b, Degree: float64(degrees)})
	if err != nil {
		return nil, err
	}
	return storage.NewImage(result.Content, result.Format), err
}

func (p *PdfProcessor) Create(ctx context.Context, title string, records []datastore.Record) ([]byte, error) {
	client := processor.NewPdfProcessorClient(p.conn)

	categories, err := p.categories.All(ctx)
	doc, err := NewDocument(title, records, p.images, categories)
	if err != nil {
		return nil, err
	}
	pdf, err := client.CreatePdf(ctx, doc)
	if err != nil {
		return nil, err
	}
	return pdf.Content, err
}

func (p *PdfProcessor) Check(ctx context.Context) (string, error) {
	state := p.conn.GetState()
	if state == connectivity.TransientFailure || state == connectivity.Shutdown {
		return state.String(), errors.New(fmt.Sprintf("grpc error. Connection state: %v", state))
	}
	return state.String(), nil
}