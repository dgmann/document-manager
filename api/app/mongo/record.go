package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgmann/document-manager/api/app"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
)

type RecordService struct {
	RecordServiceConfig
}

type recordCollection interface {
	finder
	oneFinder
	oneInserter
	oneDeleter
	oneFinderUpdater
}

type RecordServiceConfig struct {
	Records recordCollection
	Events  app.Sender
	Images  app.ResourceWriter
	Pdfs    app.ResourceWriter
}

func NewRecordService(config RecordServiceConfig) *RecordService {
	return &RecordService{config}
}

func (r *RecordService) All(ctx context.Context) ([]app.Record, error) {
	return r.Query(ctx, bson.M{})
}

func (r *RecordService) Find(ctx context.Context, id string) (*app.Record, error) {
	res, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return r.findByObjectId(ctx, res)
}

func (r *RecordService) findByObjectId(ctx context.Context, id primitive.ObjectID) (*app.Record, error) {
	var record app.Record

	res := r.Records.FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		logrus.WithField("error", res.Err()).Error("Cannot find record")
		return nil, res.Err()
	}

	if err := res.Decode(&record); err != nil {
		logrus.WithField("error", err).Error("error decoding record")
		return nil, errors.New(fmt.Sprintf("record with ID %s not found", id.Hex()))
	}
	return &record, nil
}

func (r *RecordService) Query(ctx context.Context, query map[string]interface{}) ([]app.Record, error) {
	cursor, err := r.Records.Find(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return castToRecordSlice(ctx, cursor)
}

func (r *RecordService) Create(ctx context.Context, data app.CreateRecord, images []app.Image, pdfData io.Reader) (*app.Record, error) {
	record := app.NewRecord(data)

	if len(record.Pages) == 0 {
		var pages []app.Page
		for _, img := range images {
			page := app.NewPage(img.Format)

			if err := r.Images.Write(app.NewKeyedGenericResource(img.Image, img.Format, record.Id.Hex(), page.Id)); err != nil {
				return nil, err
			}
			pages = append(pages, *page)
		}
		record.Pages = pages
	}

	pdfBytes, err := ioutil.ReadAll(pdfData)
	if err != nil {
		return nil, err
	}

	if err := r.Pdfs.Write(app.NewKeyedGenericResource(pdfBytes, "pdf", record.Id.Hex())); err != nil {
		e := r.Images.Delete(app.NewDirectoryResource(record.Id.Hex()))
		logrus.Error(e)
		return nil, err
	}

	res, err := r.Records.InsertOne(ctx, record)
	if err != nil {
		e := r.Images.Delete(app.NewDirectoryResource(record.Id.Hex()))
		logrus.Error(e)
		e = r.Pdfs.Delete(app.NewDirectoryResource(record.Id.Hex()))
		logrus.Error(e)
		return nil, err
	}
	created, err := r.findByObjectId(ctx, res.InsertedID.(primitive.ObjectID))
	if err != nil {
		return nil, err
	}

	r.Events.Send(app.EventCreated, created)
	return created, nil
}

func (r *RecordService) Delete(ctx context.Context, id string) error {
	err := r.Images.Delete(app.NewDirectoryResource(id))
	if err != nil {
		return err
	}

	err = r.Pdfs.Delete(app.NewDirectoryResource(id))
	if err != nil {
		return err
	}

	key, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.Records.DeleteOne(ctx, bson.M{"_id": key})

	r.Events.Send(app.EventDeleted, id)
	return err
}

func (r *RecordService) Update(ctx context.Context, id string, record app.Record) (*app.Record, error) {
	key, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	// TODO: Remove deleted pages from the file system
	res := r.Records.FindOneAndUpdate(ctx, bson.M{"_id": key}, bson.M{"$set": record})
	if res.Err() != nil {
		return nil, err
	}
	var updated *app.Record

	if err := res.Decode(updated); err != nil {
		return nil, err
	}

	r.Events.Send(app.EventUpdated, updated)
	return updated, nil
}

// UpdatePages updates the pages specified while keeping the rest of the original pages
func (r *RecordService) UpdatePages(ctx context.Context, id string, updates []app.PageUpdate) (*app.Record, error) {
	record, err := r.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	pages := make(map[string]app.Page)
	for _, page := range record.Pages {
		pages[page.Id] = page
	}
	var updated []app.Page
	for _, update := range updates {
		updated = append(updated, pages[update.Id])
	}
	return r.Update(ctx, id, app.Record{Pages: updated})
}

func castToRecordSlice(ctx context.Context, cursor mongo.Cursor) ([]app.Record, error) {
	records := make([]app.Record, 0)

	for cursor.Next(ctx) {
		r := app.Record{}
		if err := cursor.Decode(&r); err != nil {
			logrus.WithError(err).Error("error decoding category from database")
			return nil, err
		}
		records = append(records, r)
	}

	return records, nil
}
