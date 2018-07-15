package splitter

import (
	"io/ioutil"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
	"github.com/dgmann/document-manager/migrator/records/models"
)

func Split(path string) ([]*models.SubRecord, string, error) {
	tmpDir, err := ioutil.TempDir("", "migration_")
	if err != nil {
		return nil, tmpDir, errors.Wrap(err, "error creating tmp dir")
	}
	subrecords, err := splitByBookmarks(path, tmpDir)
	return subrecords, tmpDir, err
}

func splitByBookmarks(inputFile, outDir string) ([]*models.SubRecord, error) {
	cmd := exec.Command("java", "-jar", "./SplitPDF.jar", "-iFile", inputFile, " -CleanOutputFolder", "-oFolder", outDir)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return getSubfiles(outDir)
}

func getSubfiles(directory string) ([]*models.SubRecord, error) {
	var pdfList = make([]*models.SubRecord, 0)
	err := filepath.Walk(directory, func(path string, fi os.FileInfo, err error) error {
		if filepath.Ext(fi.Name()) == ".pdf" {
			pdfFile := &models.SubRecord{Path: path}
			pdfList = append(pdfList, pdfFile)
		}
		return nil
	})
	return pdfList, err
}
