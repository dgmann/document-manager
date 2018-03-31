package pdfprocessor

import (
	"io/ioutil"
	"os"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	_ "image/jpeg"
	_ "image/png"
	"bytes"
)

type FileWrapper struct {
	Dir string
}

type Image struct {
	*bytes.Buffer
	Type string
}

func NewFileWrapper() (*FileWrapper, error) {
	dir, err := ioutil.TempDir("", "images")
	if err != nil {
		return nil, err
	}
	return &FileWrapper{Dir: dir}, nil
}

func (f *FileWrapper) Dispose() {
	err := os.RemoveAll(f.Dir)
	if err != nil {
		log.Error(err)
	}
}

func (f *FileWrapper) GetFilePath(fp string) string {
	return filepath.Join(f.Dir, fp)
}

func (f *FileWrapper) GetImagesAsBuffer(imageType ...string) ([]*Image, error) {
	imgList := []*Image{}
	err := filepath.Walk(f.Dir, func(path string, fi os.FileInfo, err error) error {
		if contains(imageType, filepath.Ext(fi.Name())) {
			imgFile, _ := os.Open(path)
			buf := new(bytes.Buffer)
			buf.ReadFrom(imgFile)
			imgBuff := Image{Buffer: buf}
			imgBuff.Type = filepath.Ext(fi.Name())[:1]
			imgFile.Close()
			imgList = append(imgList, &imgBuff)
		}
		return nil
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return imgList, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
