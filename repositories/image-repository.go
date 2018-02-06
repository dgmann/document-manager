package repositories

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"image"
	"image/png"
	"os"
	"path"
	"path/filepath"
)

const fileExtension = ".png"

type ImageRepository interface {
	Get(id string) map[string]image.Image
	Set(id string, images []image.Image) ([]string, error)
	Delete(id string) error
	Serve(context *gin.Context, recordId string, imageId string)
}

type FileSystemImageRepository struct {
	directory string
}

func NewFileSystemImageRepository(directory string) *FileSystemImageRepository {
	return &FileSystemImageRepository{directory: directory}
}

func (f *FileSystemImageRepository) Get(id string) map[string]image.Image {
	images := make(map[string]image.Image, 0)
	p := path.Join(f.directory, id)
	filepath.Walk(p, func(d string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			dir := path.Join(d, info.Name())
			f, err := os.Open(dir)
			if err != nil {
				log.WithFields(log.Fields{
					"name":      f.Name(),
					"directory": dir,
					"error":     err,
				}).Error("Error reading image")
			}
			img, _, err := image.Decode(f)
			images[info.Name()] = img
		}
		return nil
	})
	return images
}

func (f *FileSystemImageRepository) Set(id string, images []image.Image) (results []string, err error) {
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

func (f *FileSystemImageRepository) Delete(id string) error {
	p := path.Join(f.directory, id)
	return os.RemoveAll(p)
}

func save(filePath string, img image.Image) {
	imageFile, err := os.Create(filePath + fileExtension)
	defer imageFile.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"name":      imageFile.Name(),
			"directory": filePath,
			"error":     err,
		}).Error("Error creating image file")
		panic(err)
	}
	err = png.Encode(imageFile, img)
	if err != nil {
		log.WithFields(log.Fields{
			"name":      imageFile.Name(),
			"directory": filePath,
			"error":     err,
		}).Error("Error saving image")
		panic(err)
	}
}

func (f *FileSystemImageRepository) Serve(context *gin.Context, recordId string, imageId string) {
	p := path.Join(f.directory, recordId, imageId+fileExtension)
	context.File(p)
}
