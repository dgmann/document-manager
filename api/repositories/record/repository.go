package record

import (
	"context"
	"github.com/dgmann/document-manager/api/models"
	"github.com/dgmann/document-manager/api/repositories"
	"github.com/dgmann/document-manager/api/services"
	"github.com/dgmann/document-manager/shared"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
)

type Repository interface {
	All(ctx context.Context) ([]*models.Record, error)
	Find(ctx context.Context, id string) (*models.Record, error)
	Query(ctx context.Context, query map[string]interface{}) ([]*models.Record, error)
	Create(ctx context.Context, data models.CreateRecord, images []*shared.Image, pdfData io.Reader) (*models.Record, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, record models.Record) (*models.Record, error)
	UpdatePages(ctx context.Context, id string, updates []*models.PageUpdate) (*models.Record, error)
}

type DatabaseRepository struct {
	records collection
	events  *services.EventService
	images  repositories.ResourceWriter
	pdfs    repositories.ResourceWriter
}

func NewDatabaseRepository(records collection, imageWriter repositories.ResourceWriter, pdfs repositories.ResourceWriter, eventService *services.EventService) *DatabaseRepository {
	return &DatabaseRepository{records: records, events: eventService, images: imageWriter, pdfs: pdfs}
}

func CreateIndexes(ctx context.Context, c indexer) error {
	_, err := c.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{"patientId", int32(1)}}})
	if err != nil {
		return err
	}
	return nil
}

func (r *DatabaseRepository) All(ctx context.Context) ([]*models.Record, error) {
	return r.Query(ctx, bson.M{})
}

func (r *DatabaseRepository) Find(ctx context.Context, id string) (*models.Record, error) {
	res, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return r.findByObjectId(ctx, res)
}

func (r *DatabaseRepository) findByObjectId(ctx context.Context, id primitive.ObjectID) (*models.Record, error) {
	var record models.Record

	res := r.records.FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		log.WithField("error", res.Err()).Error("Cannot find record")
		return nil, res.Err()
	}

	if err := res.Decode(record); err != nil {
		log.WithField("error", res.Err()).Error("error decoding record")
		return nil, res.Err()
	}
	return &record, nil
}

func (r *DatabaseRepository) Query(ctx context.Context, query map[string]interface{}) ([]*models.Record, error) {
	cursor, err := r.records.Find(ctx, query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return castToSlice(ctx, cursor)
}

func (r *DatabaseRepository) Create(ctx context.Context, data models.CreateRecord, images []*shared.Image, pdfData io.Reader) (*models.Record, error) {
	record := models.NewRecord(data)

	if len(record.Pages) == 0 {
		var pages []*models.Page
		for _, img := range images {
			page := models.NewPage(img.Format)

			if err := r.images.Write(repositories.NewKeyedGenericResource(img.Image, img.Format, record.Id.Hex(), page.Id)); err != nil {
				return nil, err
			}
			pages = append(pages, page)
		}
		record.Pages = pages
	}

	pdfBytes, err := ioutil.ReadAll(pdfData)
	if err != nil {
		return nil, err
	}

	if err := r.pdfs.Write(repositories.NewKeyedGenericResource(pdfBytes, "pdf", record.Id.Hex())); err != nil {
		e := r.images.Delete(repositories.NewDirectoryResource(record.Id.Hex()))
		log.Error(e)
		return nil, err
	}

	res, err := r.records.InsertOne(ctx, record)
	if err != nil {
		e := r.images.Delete(repositories.NewDirectoryResource(record.Id.Hex()))
		log.Error(e)
		e = r.pdfs.Delete(repositories.NewDirectoryResource(record.Id.Hex()))
		log.Error(e)
		return nil, err
	}
	created, err := r.findByObjectId(ctx, res.InsertedID.(primitive.ObjectID))
	if err != nil {
		return nil, err
	}

	r.events.Send(services.EventCreated, created)
	return created, nil
}

func (r *DatabaseRepository) Delete(ctx context.Context, id string) error {
	err := r.images.Delete(repositories.NewDirectoryResource(id))
	if err != nil {
		return err
	}

	err = r.pdfs.Delete(repositories.NewDirectoryResource(id))
	if err != nil {
		return err
	}

	key, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.records.DeleteOne(ctx, bson.M{"_id": key})

	r.events.Send(services.EventDeleted, id)
	return err
}

func (r *DatabaseRepository) Update(ctx context.Context, id string, record models.Record) (*models.Record, error) {
	key, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	// TODO: Remove deleted pages from the file system
	res := r.records.FindOneAndUpdate(ctx, bson.M{"_id": key}, bson.M{"$set": record})
	if res.Err() != nil {
		return nil, err
	}
	var updated *models.Record

	if err := res.Decode(updated); err != nil {
		return nil, err
	}

	r.events.Send(services.EventUpdated, updated)
	return updated, nil
}

// UpdatePages updates the pages specified while keeping the rest of the original pages
func (r *DatabaseRepository) UpdatePages(ctx context.Context, id string, updates []*models.PageUpdate) (*models.Record, error) {
	record, err := r.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	pages := make(map[string]*models.Page)
	for _, page := range record.Pages {
		pages[page.Id] = page
	}
	var updated []*models.Page
	for _, update := range updates {
		updated = append(updated, pages[update.Id])
	}
	return r.Update(ctx, id, models.Record{Pages: updated})
}
