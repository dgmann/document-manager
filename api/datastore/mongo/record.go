package mongo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/dgmann/document-manager/api/datastore"
	"github.com/dgmann/document-manager/api/event"
	"github.com/dgmann/document-manager/api/storage"
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

func (r *RecordService) All(ctx context.Context) ([]datastore.Record, error) {
	return r.Query(ctx, datastore.NewRecordQuery())
}

func (r *RecordService) Find(ctx context.Context, id string) (*datastore.Record, error) {
	res, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, datastore.NewNotFoundError(id, Records, err)
	}
	return r.findByObjectId(ctx, res)
}

func (r *RecordService) findByObjectId(ctx context.Context, id primitive.ObjectID) (*datastore.Record, error) {
	var record datastore.Record

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

func (r *RecordService) Query(ctx context.Context, recordQuery *datastore.RecordQuery, queryOptions ...*datastore.QueryOptions) ([]datastore.Record, error) {
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

func (r *RecordService) Create(ctx context.Context, data datastore.CreateRecord, images []storage.Image, pdfData io.Reader) (*datastore.Record, error) {
	record := datastore.NewRecord(data)

	if len(record.Pages) == 0 {
		var pages []datastore.Page
		for _, img := range images {
			page := datastore.NewPage(img.Format)

			if err := r.Images.Write(storage.NewKeyedGenericResource(img.Image, img.Format, record.Id.Hex(), page.Id)); err != nil {
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

	if err := r.Pdfs.Write(storage.NewKeyedGenericResource(pdfBytes, "pdf", record.Id.Hex())); err != nil {
		if e := r.Images.Delete(storage.NewKey(record.Id.Hex())); e != nil {
			err = fmt.Errorf("%s. %w", e, err)
		}
		return nil, fmt.Errorf("error storing pdf in filesystem. %w", err)
	}

	res, err := r.Records.InsertOne(ctx, record)
	if err != nil {
		if e := r.Images.Delete(storage.NewKey(record.Id.Hex())); e != nil {
			err = fmt.Errorf("%s. %w", e, err)
		}
		if e := r.Pdfs.Delete(storage.NewKey(record.Id.Hex())); e != nil {
			err = fmt.Errorf("%s. %w", e, err)
		}
		return nil, fmt.Errorf("error storing pdf in database. %w", err)
	}
	created, err := r.findByObjectId(ctx, res.InsertedID.(primitive.ObjectID))
	if err != nil {
		return nil, fmt.Errorf("error finding newly created document in database. %w", err)
	}

	r.Events.Send(event.New(event.RecordTopic, event.Created, created.Id.Hex()))
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

	r.Events.Send(event.New(event.RecordTopic, event.Deleted, id))
	return err
}

func (r *RecordService) Update(ctx context.Context, id string, record datastore.Record) (*datastore.Record, error) {
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
	var updated datastore.Record

	if err := res.Decode(&updated); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, datastore.NewNotFoundError(id, Records, err)
		}
		return nil, err
	}

	r.Events.Send(event.New(event.RecordTopic, event.Updated, updated.Id.Hex()))
	return &updated, nil
}

// UpdatePages updates the pages of a document without modifying the pages themselves.
// Useful for deleting multiple pages
func (r *RecordService) UpdatePages(ctx context.Context, id string, updates []datastore.PageUpdate) (*datastore.Record, error) {
	record, err := r.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	pages := make(map[string]datastore.Page)
	for _, page := range record.Pages {
		pages[page.Id] = page
	}
	var updated []datastore.Page
	for _, update := range updates {
		page := pages[update.Id]
		// Check if the page was really modified or just included without modification
		if update.Rotate != 0 {
			page.UpdatedAt = time.Now()
		}
		updated = append(updated, page)
	}
	return r.Update(ctx, id, datastore.Record{Pages: updated})
}

func castToRecordSlice(ctx context.Context, cursor datastore.Cursor) ([]datastore.Record, error) {
	defer cursor.Close(ctx)
	records := make([]datastore.Record, 0)

	for cursor.Next(ctx) {
		r := datastore.Record{}
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
