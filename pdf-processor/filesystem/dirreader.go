package filesystem

import (
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	"io/ioutil"
	"path"
	"strings"
)

func ReadImagesFromDirectory(dirname string) ([]*processor.Image, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	var images []*processor.Image
	for _, f := range files {
		p := path.Join(dirname, f.Name())
		extension := path.Ext(f.Name())
		content, err := ioutil.ReadFile(p)
		if err != nil {
			return nil, err
		}
		images = append(images, &processor.Image{Content: content, Format: strings.Trim(extension, ".")})
	}
	return images, nil
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
