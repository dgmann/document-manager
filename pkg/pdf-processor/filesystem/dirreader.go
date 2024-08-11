package filesystem

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/dgmann/document-manager/pkg/pdf-processor/pdf"
	"github.com/dgmann/document-manager/pkg/pdf-processor/pool"
	"github.com/dgmann/document-manager/pkg/pdf-processor/processor"
)

func ReadImagesFromDirectory(dirname string, writer pdf.ImageSender) (int, error) {
	files, err := os.ReadDir(dirname)
	if err != nil {
		return 0, err
	}
	imagesSent := 0
	for _, f := range files {
		p := path.Join(dirname, f.Name())
		if err := sendFile(f, p, writer); err != nil {
			return imagesSent, err
		}
		imagesSent++
	}
	return imagesSent, nil
}

func sendFile(f os.DirEntry, filePath string, writer pdf.ImageSender) error {
	extension := path.Ext(f.Name())
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	_, err = buf.ReadFrom(file)
	if err != nil {
		return err
	}
	return writer.Send(&processor.Image{Content: buf.Bytes(), Format: strings.Trim(extension, ".")})
}

func ReadFiles(dirname string) ([]string, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	paths := make([]string, len(files))
	for i, f := range files {
		paths[i] = path.Join(dirname, f.Name())
	}
	return paths, nil
}

func ToImage(file string) (*processor.Image, error) {
	extension := path.Ext(file)
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return &processor.Image{Content: content, Format: strings.Trim(extension, ".")}, nil
}
