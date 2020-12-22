package filesystem

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/dgmann/document-manager/pdf-processor/pkg/pdf"
	"github.com/dgmann/document-manager/pdf-processor/pkg/pool"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
)

func ReadImagesFromDirectory(dirname string, writer pdf.ImageSender) (int, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return 0, err
	}
	imagesSent := 0
	for _, f := range files {
		p := path.Join(dirname, f.Name())
		if err := sendFile(f, p, writer); err != nil {
			return imagesSent, err
		}
	}
	return imagesSent, nil
}

func sendFile(f os.FileInfo, p string, writer pdf.ImageSender) error {
	extension := path.Ext(f.Name())
	file, err := os.Open("file.go")
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
