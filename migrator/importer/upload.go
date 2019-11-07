package importer

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgmann/document-manager/api/client"
	"github.com/sirupsen/logrus"
	"os"
)

type Importer struct {
	Api        *client.HttpUploader
	retryCount int
}

func NewImporter(url string, retryCount int) *Importer {
	return &Importer{Api: client.NewHttpUploader(url), retryCount: retryCount}
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
		for j := 0; j < i.retryCount; j++ {
			err = i.Api.CreateRecord(&r.NewRecord)
			if err == nil {
				return nil
			}
			logrus.WithError(err).Errorf("error uploading file. Retry %d of %d", j, i.retryCount)
		}
		return errors.New(fmt.Sprintf("%s: %s", r.Path, err.Error()))
	}
}
