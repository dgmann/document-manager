package repositories

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"github.com/dgmann/document-manager/shared"
)

type ImageRepository interface {
	Get(id string) map[string]io.Reader
	Set(id string, images []*shared.Image) ([]string, error)
	SetImage(id string, fileName string, image *shared.Image) error
	Delete(id string) error
	Serve(context *gin.Context, recordId string, imageId string)
}

type FileSystemImageRepository struct {
	directory string
}

func NewFileSystemImageRepository(directory string) *FileSystemImageRepository {
	return &FileSystemImageRepository{directory: directory}
}

func (f *FileSystemImageRepository) Get(id string) map[string]io.Reader {
	images := make(map[string]io.Reader, 0)
	p := path.Join(f.directory, id)
	filepath.Walk(p, func(d string, info os.FileInfo, err error) error {
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
			}
			fileName := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
			images[fileName] = f
		}
		return nil
	})
	return images
}

func (f *FileSystemImageRepository) Copy(fromId string, toId string) error {
	sourceFolder := path.Join(f.directory, fromId)
	destinationFolder := path.Join(f.directory, toId)
	return copyFolder(sourceFolder, destinationFolder)
}

func (f *FileSystemImageRepository) Set(id string, images []*shared.Image) (results []string, err error) {
	results = make([]string, 0)
	p := path.Join(f.directory, id)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		os.MkdirAll(p, os.ModePerm)
	}
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovering: %v", r)
			os.RemoveAll(p)
			err = errors.New("failed to save images")
			results = make([]string, 0)
		}
	}()
	for _, img := range images {
		imgId := uuid.New().String()
		fp := path.Join(p, imgId)
		save(fp, img)
		results = append(results, imgId)
	}
	return results, nil
}

func (f *FileSystemImageRepository) SetImage(recordId string, pageId string, image *shared.Image) error {
	p := f.getPath(recordId, pageId)
	err := save(p, image)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileSystemImageRepository) Delete(id string) error {
	p := path.Join(f.directory, id)
	return os.RemoveAll(p)
}

func save(filePath string, img *shared.Image) error {
	filePath = filePath + "." + img.Format
	imageFile, err := os.Create(filePath)
	defer imageFile.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"name":      imageFile.Name(),
			"directory": filePath,
			"error":     err,
		}).Error("Error creating image file")
		return err
	}
	_, err = imageFile.Write(img.Image)
	return err
}

func (f *FileSystemImageRepository) Serve(context *gin.Context, recordId string, imageId string) {
	p := f.getPath(recordId, imageId)
	context.File(p)
}

func (f *FileSystemImageRepository) getPath(recordId string, imageId string) string {
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
