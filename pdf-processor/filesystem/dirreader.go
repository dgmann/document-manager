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
