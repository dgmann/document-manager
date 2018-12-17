package image

import (
	"fmt"
	"github.com/dgmann/document-manager/api/repositories"
	"github.com/dgmann/document-manager/api/repositories/filesystem"
	"github.com/dgmann/document-manager/api/services"
	"github.com/dgmann/document-manager/shared"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Repository interface {
	services.FileInfoService
	repositories.ResourceWriter
	Get(id string) (map[string]*shared.Image, error)
	Serve(context *gin.Context, recordId string, imageId string, format string)
	Copy(fromId string, toId string) error
}

type FileSystemRepository struct {
	*filesystem.Repository
	directory string
}

func NewFileSystemRepository(directory string) *FileSystemRepository {
	return &FileSystemRepository{directory: directory, Repository: filesystem.NewRepository(directory)}
}

func (f *FileSystemRepository) Get(id string) (map[string]*shared.Image, error) {
	images := make(map[string]*shared.Image, 0)
	p := path.Join(f.directory, id)
	err := filepath.Walk(p, func(d string, info os.FileInfo, err error) error {
		if err != nil {
			log.Error(err)
			return err
		}
		if !info.IsDir() {
			f, err := os.Open(d)
			if err != nil {
				log.WithFields(log.Fields{
					"name":      f.Name(),
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
			images[fileName] = shared.NewImage(data, strings.Trim(ext, "."))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (f *FileSystemRepository) GetFileInfo(recordId, pageId string, format string) (os.FileInfo, error) {
	p := f.getPath(recordId, pageId+"."+format)
	return os.Stat(p)
}

func (f *FileSystemRepository) Copy(fromId string, toId string) error {
	sourceFolder := path.Join(f.directory, fromId)
	destinationFolder := path.Join(f.directory, toId)
	return copyFolder(sourceFolder, destinationFolder)
}

func (f *FileSystemRepository) Serve(context *gin.Context, recordId string, imageId string, format string) {
	p := f.getPath(recordId, imageId+"."+format)
	context.File(p)
}

func (f *FileSystemRepository) getPath(recordId string, imageId string) string {
	return path.Join(f.directory, recordId, imageId)
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
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}

func normalizeExtension(extension string) string {
	if extension == "jpeg" {
		return "jpg"
	}
	return extension
}
