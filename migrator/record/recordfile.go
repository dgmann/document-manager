package record

import (
	"os/exec"
	"io/ioutil"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

type RecordFile struct {
	filePath       string
	splittedPdfDir string
	SubFiles       []*PdfFile
}

func NewRecordFile(path string) (*RecordFile, error) {
	tmpDir, err := ioutil.TempDir("", "migration_")
	if err != nil {
		return nil, errors.Wrap(err, "error creating tmp dir")
	}
	r := &RecordFile{filePath: path, splittedPdfDir: tmpDir}
	r.SubFiles, err = splitByBookmarks(path, tmpDir)
	if err != nil {
		return nil, errors.Wrap(err, "error splitting record in subfiles")
	}
	return r, nil
}

func (r *RecordFile) Destroy() error {
	return os.RemoveAll(r.splittedPdfDir)
}

func splitByBookmarks(inputFile, outDir string) ([]*PdfFile, error) {
	cmd := exec.Command("java", "-jar", "./SplitPDF.jar", "-iFile", inputFile, " -CleanOutputFolder", "-oFolder", outDir)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return getSubfiles(outDir)
}

func getSubfiles(directory string) ([]*PdfFile, error) {
	var pdfList = make([]*PdfFile, 0)
	err := filepath.Walk(directory, func(path string, fi os.FileInfo, err error) error {
		if filepath.Ext(fi.Name()) == ".pdf" {
			pdfFile := NewPdfFile(path)
			pdfList = append(pdfList, pdfFile)
		}
		return nil
	})
	return pdfList, err
}
