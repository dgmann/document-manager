package repositories

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"image"
	"io"
	"os"
	"path"
	"path/filepath"
)

const fileExtension = ".png"

type ImageRepository interface {
	Get(id string) map[string]image.Image
	Set(id string, images []*bytes.Buffer) ([]string, error)
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

func (f *FileSystemImageRepository) Copy(fromId string, toId string) error {
	sourceFolder := path.Join(f.directory, fromId)
	destinationFolder := path.Join(f.directory, toId)
	return copyFolder(sourceFolder, destinationFolder)
}

func (f *FileSystemImageRepository) Set(id string, images []*bytes.Buffer) (results []string, err error) {
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
		save(fp, img.Bytes())
		results = append(results, imgId)
	}
	return results, nil
}

func (f *FileSystemImageRepository) Delete(id string) error {
	p := path.Join(f.directory, id)
	return os.RemoveAll(p)
}

func save(filePath string, img []byte) {
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
	imageFile.Write(img)
}

func (f *FileSystemImageRepository) Serve(context *gin.Context, recordId string, imageId string) {
	p := path.Join(f.directory, recordId, imageId+fileExtension)
	context.File(p)
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
