package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgmann/document-manager/api/app"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
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
		return nil, app.NewNotFoundError(id, Records, err)
	}
	return r.findByObjectId(ctx, res)
}

func (r *RecordService) findByObjectId(ctx context.Context, id primitive.ObjectID) (*app.Record, error) {
	var record app.Record

	res := r.Records.FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		return nil, res.Err()
	}

	if err := res.Decode(&record); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, app.NewNotFoundError(id.Hex(), Records, err)
		}
		return nil, err
	}
	return &record, nil
}

func (r *RecordService) Query(ctx context.Context, query map[string]interface{}) ([]app.Record, error) {
	cursor, err := r.Records.Find(ctx, query)
	if err != nil {
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
		return app.NewNotFoundError(id, Records, err)
	}
	res, err := r.Records.DeleteOne(ctx, bson.M{"_id": key})
	if err != nil {
		return fmt.Errorf("while deleting record %s. %w", id, err)
	}
	if res.DeletedCount == 0 {
		return app.NewNotFoundError(id, Records, err)
	}

	r.Events.Send(app.EventDeleted, id)
	return err
}

func (r *RecordService) Update(ctx context.Context, id string, record app.Record) (*app.Record, error) {
	key, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, app.NewNotFoundError(id, Records, err)
	}
	// TODO: Remove deleted pages from the file system
	res := r.Records.FindOneAndUpdate(ctx, bson.M{"_id": key}, bson.M{"$set": record}, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if res.Err() != nil {
		return nil, err
	}
	var updated app.Record

	if err := res.Decode(&updated); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, app.NewNotFoundError(id, Records, err)
		}
		return nil, err
	}

	r.Events.Send(app.EventUpdated, updated)
	return &updated, nil
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
			return nil, fmt.Errorf("decoding records from database: %w", err)
		}
		records = append(records, r)
	}

	return records, nil
}
