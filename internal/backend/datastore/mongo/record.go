package mongo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/dgmann/document-manager/internal/backend/datastore"
	"github.com/dgmann/document-manager/pkg/api"
	log "github.com/sirupsen/logrus"

	"github.com/dgmann/document-manager/internal/backend/event"
	"github.com/dgmann/document-manager/internal/backend/storage"
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
	Events  event.Sender[*api.Record]
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
		return nil, datastore.NewNotFoundError(id, CollectionRecords, err)
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
			return nil, datastore.NewNotFoundError(id.Hex(), CollectionRecords, datastore.ErrNoDocuments)
		}
		return nil, err
	}
	return &record, nil
}

func (r *RecordService) Query(ctx context.Context, recordQuery datastore.RecordQuery, queryOptions ...*datastore.QueryOptions) ([]api.Record, error) {
	op := make([]*options.FindOptions, len(queryOptions))
	for index, option := range queryOptions {
		op[index] = options.Find().SetSkip(option.Skip).SetLimit(option.Limit).SetSort(option.Sort)
	}

	cursor, err := r.Records.Find(ctx, recordQuery, op...)
	if err != nil {
		return nil, err
	}

	return castToRecordSlice(ctx, cursor)
}

func (r *RecordService) Create(ctx context.Context, data api.CreateRecord, images []storage.Image, pdfData io.Reader) (created *api.Record, resErr error) {
	record := newRecord(data)
	recordId := record.Id.Hex()

	// ID IS not set here!
	if len(record.Pages) == 0 {
		var pages []api.Page
		for _, img := range images {
			page := datastore.NewPage(img.Format)

			if err := r.Images.Write(storage.NewKeyedGenericResource(img.Image, img.Format, recordId, page.Id)); err != nil {
				return nil, err
			}
			pages = append(pages, *page)
		}
		defer func() {
			if resErr != nil {
				if e := r.Images.Delete(storage.NewKey(recordId)); e != nil {
					log.Debugf("error deleting images during cleanup: %s", e)
				}
			}
		}()
		record.Pages = pages
	}

	pdfBytes, err := io.ReadAll(pdfData)
	if err != nil {
		return nil, err
	}

	if err := r.Pdfs.Write(storage.NewKeyedGenericResource(pdfBytes, "pdf", recordId)); err != nil {
		return nil, fmt.Errorf("error storing pdf in filesystem. %w", err)
	}
	defer func() {
		if resErr != nil {
			if e := r.Pdfs.Delete(storage.NewKey(recordId)); e != nil {
				log.Debugf("error deleting pdf during cleanup: %s", e)
			}
		}
	}()

	if _, err := r.Records.InsertOne(ctx, record); err != nil {
		return nil, fmt.Errorf("error storing pdf in database. %w", err)
	}
	created, err = r.findByObjectId(ctx, record.Id)
	if err != nil {
		return nil, fmt.Errorf("error finding newly created document in database. %w", err)
	}

	if err := r.Events.Send(ctx, api.NewEvent(api.RecordTopic, api.EventTypeCreated, created.Id, created)); err != nil {
		log.WithError(err).Info("error sending event")
	}
	return created, nil
}

func newRecord(data api.CreateRecord) *datastore.Record {
	record := &datastore.Record{
		Id: primitive.NewObjectID(),
		Record: &api.Record{
			Date:       nil,
			ReceivedAt: time.Now(),
			Comment:    data.Comment,
			PatientId:  data.PatientId,
			Sender:     data.Sender,
			Tags:       &data.Tags,
			Pages:      data.Pages,
			Status:     &data.Status,
			Category:   data.Category,
			UpdatedAt:  time.Now(),
		},
	}
	if !data.Date.IsZero() {
		record.Date = &data.Date
	}
	if !data.ReceivedAt.IsZero() {
		record.ReceivedAt = data.ReceivedAt
	}
	if len(*record.Tags) == 0 {
		record.Tags = &[]string{}
	}
	if len(record.Pages) == 0 {
		record.Pages = []api.Page{}
	}
	if record.Status.IsNone() {
		status := api.StatusInbox
		record.Status = &status
	}
	return record
}

func (r *RecordService) Delete(ctx context.Context, id string) error {
	key, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return datastore.NewNotFoundError(id, CollectionRecords, err)
	}
	res, err := r.Records.DeleteOne(ctx, bson.M{"_id": key})
	if err != nil {
		return fmt.Errorf("while deleting record %s. %w", id, err)
	}
	if res.DeletedCount == 0 {
		return datastore.NewNotFoundError(id, CollectionRecords, err)
	}

	if err := r.Images.Delete(storage.NewKey(id)); err != nil {
		log.WithError(err).WithField("recordId", id).Warn("error deleting images")
	}
	if err := r.Pdfs.Delete(storage.NewKey(id)); err != nil {
		log.WithError(err).WithField("recordId", id).Warn("error deleting PDF")
	}

	if err := r.Events.Send(ctx, api.NewEvent[*api.Record](api.RecordTopic, api.EventTypeDeleted, id, nil)); err != nil {
		log.WithError(err).Info("error sending event")
	}
	return nil
}

func (r *RecordService) Update(ctx context.Context, id string, record api.Record, updateOptions ...datastore.UpdateOption) (*api.Record, error) {
	key, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, datastore.NewNotFoundError(id, CollectionRecords, err)
	}
	// Unset record.Id as it is an immutable field
	record.Id = ""
	record.UpdatedAt = time.Now()
	// TODO: Remove deleted pages from the file system
	query := bson.M{"_id": key}
	for _, opts := range updateOptions {
		opts(query)
	}
	res := r.Records.FindOneAndUpdate(ctx, query, bson.M{"$set": record}, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if err := res.Err(); res.Err() != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, datastore.ErrNoDocuments
		}
		return nil, err
	}

	var updated api.Record
	if err := res.Decode(&updated); err != nil {
		return nil, err
	}

	if err := r.Events.Send(ctx, api.NewEvent(api.RecordTopic, api.EventTypeUpdated, updated.Id, &updated)); err != nil {
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
		if update.Content != nil && len(*update.Content) > 0 {
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
