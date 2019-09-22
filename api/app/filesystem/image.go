package filesystem

import (
	"fmt"
	"github.com/dgmann/document-manager/api/app"
	"io"
	"os"
)

type ImageService struct {
	*DiskStorage
}

func NewImageService(directory string) (*ImageService, error) {
	repository, err := NewDiskStorage(directory)
	if err != nil {
		return nil, err
	}
	return &ImageService{DiskStorage: repository}, nil
}

func (f *ImageService) Get(id string) (map[string]*app.Image, error) {
	images := make(map[string]*app.Image, 0)
	p := app.NewKey(id)
	err := f.ForEach(p, func(resource app.KeyedResource, err error) error {
		fileName := resource.Key()[len(resource.Key())-1]
		images[fileName] = app.NewImage(resource.Data(), resource.Format())
		return err
	})
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (f *ImageService) Copy(fromId string, toId string) error {
	sourceFolder := f.Locate(app.NewKey(fromId))
	destinationFolder := f.Locate(app.NewKey(toId))
	return copyFolder(sourceFolder, destinationFolder)
}

func (f *ImageService) Path(locatable app.Locatable) string {
	return f.Locate(locatable)
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

	if _, err = io.Copy(destfile, sourcefile); err == nil {
		if sourceinfo, statErr := os.Stat(source); statErr != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}
	}
	return
}
