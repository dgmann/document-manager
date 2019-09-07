package filesystem

import (
	"fmt"
	"github.com/dgmann/document-manager/api/app"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ImageService struct {
	*Storage
}

func NewImageService(directory string) (*ImageService, error) {
	repository, err := New(directory)
	if err != nil {
		return nil, err
	}
	return &ImageService{Storage: repository}, nil
}

func (f *ImageService) Get(id string) (map[string]*app.Image, error) {
	images := make(map[string]*app.Image, 0)
	p := path.Join(f.baseDirectory, id)
	err := filepath.Walk(p, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.WithError(err).Error("error getting files")
			return err
		}
		if !info.IsDir() {
			f, err := os.Open(currentPath)

			if err != nil {
				d, filename := filepath.Split(currentPath)
				logrus.WithFields(logrus.Fields{
					"name":      filename,
					"directory": d,
					"error":     err,
				}).Error("Error reading image")
				return err
			}
			fileName := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
			data, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}
			ext := filepath.Ext(info.Name())
			images[fileName] = app.NewImage(data, strings.Trim(ext, "."))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (f *ImageService) Copy(fromId string, toId string) error {
	sourceFolder := path.Join(f.baseDirectory, fromId)
	destinationFolder := path.Join(f.baseDirectory, toId)
	return copyFolder(sourceFolder, destinationFolder)
}

func (f *ImageService) Serve(context *gin.Context, recordId string, imageId string, format string) {
	p := f.getPath(recordId, imageId+"."+format)
	context.File(p)
}

func (f *ImageService) getPath(recordId string, imageId string) string {
	return path.Join(f.baseDirectory, recordId, imageId)
}

func copyFolder(source string, dest string) (err error) {

	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			err = copyFolder(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err = copyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

func copyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return
	}

	defer func() {
		cerr := destfile.Close()
		if err == nil {
			err = cerr
		}
	}()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}
	return
}
