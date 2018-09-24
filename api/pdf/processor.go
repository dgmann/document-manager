package pdf

import (
	"context"
	"encoding/json"
	"github.com/dgmann/document-manager/pdf-processor/processor"
	"github.com/dgmann/document-manager/shared"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"strconv"
)

type Processor struct {
	baseUrl string
	conn    *grpc.ClientConn
}

func NewPDFProcessor(baseUrl string) (*Processor, error) {
	conn, err := grpc.Dial(baseUrl)
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

	client := processor.NewPdfProcessorClient(p.conn)
	stream, err := client.ConvertPdfToImage(context.Background(), &processor.Pdf{Content: b})

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

func (p *Processor) httpRequest(f io.Reader) ([]*shared.Image, error) {
	requester := NewHttpRequester(p.baseUrl + "/images/convert")
	result, err := p.Upload(requester, f)
	if err != nil {
		log.Errorf("Error fetching images: %s", err)
		return nil, err
	}
	return result, nil
}

func (p *Processor) Upload(requester Requester, file io.Reader) ([]*shared.Image, error) {
	result, err := requester.Do(file)
	if err != nil {
		log.WithField("error", err).Error("Error transforming pdf to images")
		return nil, err
	}
	defer result.Close()
	var images []*shared.Image
	if json.NewDecoder(result).Decode(&images) != nil {
		log.WithField("error", err).Error("Error decoding response")
		return nil, err
	}
	return images, nil
}

func (p *Processor) Rotate(image io.Reader, degrees int) (*shared.Image, error) {
	requester := NewHttpRequester(p.baseUrl + "/images/rotate/" + strconv.Itoa(degrees))
	result, err := requester.Do(image)
	if err != nil {
		log.WithField("error", err).Error("Error rotating image")
		return nil, err
	}
	defer result.Close()
	var img shared.Image
	if json.NewDecoder(result).Decode(&img) != nil {
		log.WithField("error", err).Error("Error decoding response")
		return nil, err
	}
	return &img, nil
}
