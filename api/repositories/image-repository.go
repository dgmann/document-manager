package repositories

import (
	"errors"
	"fmt"
	"github.com/dgmann/document-manager/api/models"
	"github.com/dgmann/document-manager/api/services"
	"github.com/dgmann/document-manager/shared"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ImageRepository interface {
	services.FileInfoService
	Get(id string) (map[string]*shared.Image, error)
	Set(id string, images []*shared.Image) ([]*models.Page, error)
	SetImage(id string, fileName string, image *shared.Image) error
	Delete(id string) error
	Serve(context *gin.Context, recordId string, imageId string, format string)
	Copy(fromId string, toId string) error
}

type FileSystemImageRepository struct {
	directory string
}

func newFileSystemImageRepository(directory string) *FileSystemImageRepository {
	return &FileSystemImageRepository{directory: directory}
}

func (f *FileSystemImageRepository) Get(id string) (map[string]*shared.Image, error) {
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
			images[fileName] = shared.NewImage(data, filepath.Ext(info.Name()))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (f *FileSystemImageRepository) GetFileInfo(recordId, pageId string, format string) (os.FileInfo, error) {
	p := f.getPath(recordId, pageId+"."+format)
	return os.Stat(p)
}

func (f *FileSystemImageRepository) Copy(fromId string, toId string) error {
	sourceFolder := path.Join(f.directory, fromId)
	destinationFolder := path.Join(f.directory, toId)
	return copyFolder(sourceFolder, destinationFolder)
}

func (f *FileSystemImageRepository) Set(id string, images []*shared.Image) (results []*models.Page, err error) {
	results = make([]*models.Page, 0)
	p := path.Join(f.directory, id)
	if err := os.RemoveAll(p); err != nil {
		log.WithField("recordId", id).Debug("image folder did not exist yet")
	}
	if _, err := os.Stat(p); os.IsNotExist(err) {
		os.MkdirAll(p, os.ModePerm)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovering: %v", r)
			os.RemoveAll(p)
			err = errors.New("failed to save images")
			results = make([]*models.Page, 0)
		}
	}()
	for _, img := range images {
		imgId := bson.NewObjectId().Hex()
		fp := path.Join(p, imgId)
		save(fp, img)
		results = append(results, models.NewPage(imgId, img.Format))
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
	err := os.RemoveAll(p)
	if !os.IsNotExist(err) {
		return err
	}
	log.Infof("%s cannot be deleted as it does not exist", p)
	return nil
}

func save(filePath string, img *shared.Image) error {
	filePath = filePath + "." + normalizeExtension(img.Format)
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

func (f *FileSystemImageRepository) Serve(context *gin.Context, recordId string, imageId string, format string) {
	p := f.getPath(recordId, imageId+"."+format)
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

func normalizeExtension(extension string) string {
	if extension == "jpeg" {
		return "jpg"
	}
	return extension
}