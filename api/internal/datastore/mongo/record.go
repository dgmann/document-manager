package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgmann/document-manager/api/internal/datastore"
	"github.com/dgmann/document-manager/api/pkg/api"
	log "github.com/sirupsen/logrus"
	"io"
	"time"

	"github.com/dgmann/document-manager/api/internal/event"
	"github.com/dgmann/document-manager/api/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RecordService struct {
	RecordServiceConfig
}

type recordCollection interface {
	finder
	oneInserter
	oneDeleter
	oneFinderUpdater
}

type RecordServiceConfig struct {
	Records recordCollection
	Events  event.Sender
	Images  storage.ResourceWriter
	Pdfs    storage.ResourceWriter
}

func NewRecordService(config RecordServiceConfig) *RecordService {
	return &RecordService{config}
}

func (r *RecordService) All(ctx context.Context) ([]api.Record, error) {
	return r.Query(ctx, datastore.NewRecordQuery())
}

func (r *RecordService) Find(ctx context.Context, id string) (*api.Record, error) {
	res, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, datastore.NewNotFoundError(id, Records, err)
	}
	return r.findByObjectId(ctx, res)
}

func (r *RecordService) findByObjectId(ctx context.Context, id primitive.ObjectID) (*api.Record, error) {
	var record api.Record

	res := r.Records.FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		return nil, res.Err()
	}

	if err := res.Decode(&record); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, datastore.NewNotFoundError(id.Hex(), Records, err)
		}
		return nil, err
	}
	return &record, nil
}

func (r *RecordService) Query(ctx context.Context, recordQuery *datastore.RecordQuery, queryOptions ...*datastore.QueryOptions) ([]api.Record, error) {
	op := make([]*options.FindOptions, len(queryOptions))
	for index, option := range queryOptions {
		op[index] = options.Find().SetSkip(option.Skip).SetLimit(option.Limit).SetSort(option.Sort)
	}

	query, err := recordQuery.ToMap()
	if err != nil {
		return nil, err
	}

	cursor, err := r.Records.Find(ctx, query, op...)
	if err != nil {
		return nil, err
	}

	return castToRecordSlice(ctx, cursor)
}

func (r *RecordService) Create(ctx context.Context, data api.CreateRecord, images []storage.Image, pdfData io.Reader) (*api.Record, error) {
	record := api.NewRecord(data)

	if len(record.Pages) == 0 {
		var pages []api.Page
		for _, img := range images {
			page := datastore.NewPage(img.Format)

			if err := r.Images.Write(storage.NewKeyedGenericResource(img.Image, img.Format, record.Id, page.Id)); err != nil {
				return nil, err
			}
			pages = append(pages, *page)
		}
		record.Pages = pages
	}

	pdfBytes, err := io.ReadAll(pdfData)
	if err != nil {
		return nil, err
	}

	if err := r.Pdfs.Write(storage.NewKeyedGenericResource(pdfBytes, "pdf", record.Id)); err != nil {
		if e := r.Images.Delete(storage.NewKey(record.Id)); e != nil {
			err = fmt.Errorf("%s. %w", e, err)
		}
		return nil, fmt.Errorf("error storing pdf in filesystem. %w", err)
	}

	res, err := r.Records.InsertOne(ctx, record)
	if err != nil {
		if e := r.Images.Delete(storage.NewKey(record.Id)); e != nil {
			err = fmt.Errorf("%s. %w", e, err)
		}
		if e := r.Pdfs.Delete(storage.NewKey(record.Id)); e != nil {
			err = fmt.Errorf("%s. %w", e, err)
		}
		return nil, fmt.Errorf("error storing pdf in database. %w", err)
	}
	created, err := r.findByObjectId(ctx, res.InsertedID.(primitive.ObjectID))
	if err != nil {
		return nil, fmt.Errorf("error finding newly created document in database. %w", err)
	}

	if err := r.Events.Send(ctx, api.New(api.RecordTopic, api.TypeCreated, created.Id, created)); err != nil {
		log.WithError(err).Info("error sending event")
	}
	return created, nil
}

func (r *RecordService) Delete(ctx context.Context, id string) error {
	err := r.Images.Delete(storage.NewKey(id))
	if err != nil {
		return err
	}

	err = r.Pdfs.Delete(storage.NewKey(id))
	if err != nil {
		return err
	}

	key, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return datastore.NewNotFoundError(id, Records, err)
	}
	res, err := r.Records.DeleteOne(ctx, bson.M{"_id": key})
	if err != nil {
		return fmt.Errorf("while deleting record %s. %w", id, err)
	}
	if res.DeletedCount == 0 {
		return datastore.NewNotFoundError(id, Records, err)
	}

	if err := r.Events.Send(ctx, api.New(api.RecordTopic, api.TypeDeleted, id, nil)); err != nil {
		log.WithError(err).Info("error sending event")
	}
	return err
}

func (r *RecordService) Update(ctx context.Context, id string, record api.Record) (*api.Record, error) {
	key, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, datastore.NewNotFoundError(id, Records, err)
	}
	record.UpdatedAt = time.Now()
	// TODO: Remove deleted pages from the file system
	res := r.Records.FindOneAndUpdate(ctx, bson.M{"_id": key}, bson.M{"$set": record}, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if res.Err() != nil {
		return nil, err
	}
	var updated api.Record

	if err := res.Decode(&updated); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, datastore.NewNotFoundError(id, Records, err)
		}
		return nil, err
	}

	if err := r.Events.Send(ctx, api.New(api.RecordTopic, api.TypeUpdated, updated.Id, updated)); err != nil {
		log.WithError(err).Info("error sending event")
	}
	return &updated, nil
}

// UpdatePages updates the pages of a document without modifying the pages themselves.
// Useful for deleting multiple pages
func (r *RecordService) UpdatePages(ctx context.Context, id string, updates []api.PageUpdate) (*api.Record, error) {
	record, err := r.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	pages := make(map[string]api.Page)
	for _, page := range record.Pages {
		pages[page.Id] = page
	}
	var updated []api.Page
	for _, update := range updates {
		page := pages[update.Id]
		// Check if the page was really modified or just included without modification
		if update.Rotate != 0 {
			page.UpdatedAt = time.Now()
		}
		if len(update.Content) > 0 {
			page.Content = update.Content
		}
		updated = append(updated, page)
	}
	return r.Update(ctx, id, api.Record{Pages: updated})
}

func castToRecordSlice(ctx context.Context, cursor datastore.Cursor) ([]api.Record, error) {
	defer cursor.Close(ctx)
	records := make([]api.Record, 0)

	for cursor.Next(ctx) {
		r := api.Record{}
		if err := cursor.Decode(&r); err != nil {
			return nil, fmt.Errorf("decoding records from database: %w", err)
		}
		records = append(records, r)
	}

	if cursor.Err() != nil {
		return nil, cursor.Err()
	}

	return records, nil
}