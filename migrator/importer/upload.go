package importer

import (
	"context"
	"fmt"
	"github.com/dgmann/document-manager/api/client"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

type Importer struct {
	Api        *client.HttpUploader
	retryCount int
}

func NewImporter(url string, retryCount int, timeout time.Duration) *Importer {
	return &Importer{Api: client.NewHttpUploader(url, timeout), retryCount: retryCount}
}

func (i *Importer) ImportRecords(ctx context.Context, records []ImportableRecord) (<-chan *ImportableRecord, <-chan ImportError) {
	return parallel(ctx, records, i.uploadFunc())
}

func (i *Importer) uploadFunc() func(r *ImportableRecord) error {
	return func(r *ImportableRecord) error {
		f, err := os.Open(r.Path)
		if err != nil {
			return err
		}
		defer f.Close()
		r.File = f

		logrus.WithField("record", r).Debug("upload file")
		for j := 1; j <= i.retryCount; j++ {

			err = i.Api.CreateRecord(&r.NewRecord)
			if err == nil {
				return nil
			}
			logrus.Warnf("error uploading file. Retry %d of %d, %s", j, i.retryCount, err.Error())
			if _, err := f.Seek(0, io.SeekStart); err != nil {
				return fmt.Errorf("error resetting PDF file stream: %w", err)
			}
		}
		return fmt.Errorf("error uploading file. %s: %w", r.Path, err)
	}
}
