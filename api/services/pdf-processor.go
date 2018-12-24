package services

import (
	"context"
	"github.com/dgmann/document-manager/pdf-processor/api"
	"github.com/dgmann/document-manager/shared"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
)

type PdfProcessor struct {
	baseUrl string
	conn    *grpc.ClientConn
}

func NewPDFProcessor(baseUrl string) (*PdfProcessor, error) {
	conn, err := grpc.Dial(
		baseUrl,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*100), grpc.MaxCallSendMsgSize(1024*1024*100)),
	)
	if err != nil {
		return nil, err
	}

	return &PdfProcessor{conn: conn, baseUrl: baseUrl}, nil
}

func (p *PdfProcessor) Close() error {
	return p.conn.Close()
}

func (p *PdfProcessor) Convert(f io.Reader) ([]*shared.Image, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	client := api.NewPdfProcessorClient(p.conn)
	stream, err := client.ConvertPdfToImage(context.Background(), &api.Pdf{Content: b})
	if err != nil {
		return nil, err
	}

	var images []*shared.Image
	for {
		image, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		images = append(images, shared.NewImage(image.Content, image.Format))
	}
	return images, nil
}

func (p *PdfProcessor) Rotate(image io.Reader, degrees int) (*shared.Image, error) {
	b, err := ioutil.ReadAll(image)
	if err != nil {
		return nil, err
	}

	client := api.NewPdfProcessorClient(p.conn)
	result, err := client.RotateImage(context.Background(), &api.Rotate{Content: b, Degree: float64(degrees)})
	return shared.NewImage(result.Content, result.Format), err
}
