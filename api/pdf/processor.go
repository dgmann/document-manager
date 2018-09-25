package pdf

import (
	"context"
	"github.com/dgmann/document-manager/pdf-processor/api"
	"github.com/dgmann/document-manager/shared"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
)

type Processor struct {
	baseUrl string
	conn    *grpc.ClientConn
}

func NewPDFProcessor(baseUrl string) (*Processor, error) {
	insecure := grpc.WithInsecure()
	conn, err := grpc.Dial(baseUrl, insecure)
	if err != nil {
		return nil, err
	}

	return &Processor{conn: conn, baseUrl: baseUrl}, nil
}

func (p *Processor) Close() {
	p.conn.Close()
}

func (p *Processor) Convert(f io.Reader) ([]*shared.Image, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	client := api.NewPdfProcessorClient(p.conn)
	stream, err := client.ConvertPdfToImage(context.Background(), &api.Pdf{Content: b})

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

func (p *Processor) Rotate(image io.Reader, degrees int) (*shared.Image, error) {
	b, err := ioutil.ReadAll(image)
	if err != nil {
		return nil, err
	}

	client := api.NewPdfProcessorClient(p.conn)
	result, err := client.RotateImage(context.Background(), &api.Rotate{Content: b, Degree: float64(degrees)})
	return shared.NewImage(result.Content, result.Format), err
}
