package repositories

import (
	"github.com/dgmann/document-manager/api/models"
	"github.com/dgmann/document-manager/api/services"
	"github.com/dgmann/document-manager/shared"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
)

type RecordRepository interface {
	All() ([]*models.Record, error)
	Find(id string) (*models.Record, error)
	Query(query map[string]interface{}) ([]*models.Record, error)
	Create(data models.CreateRecord, images []*shared.Image, pdfData io.Reader) (*models.Record, error)
	Delete(id string) error
	Update(id string, record models.Record) (*models.Record, error)
	UpdatePages(id string, updates []*models.PageUpdate) (*models.Record, error)
}

type DBRecordRepository struct {
	records *mgo.Collection
	events  *services.EventService
	images  ResourceWriter
	pdfs    ResourceWriter
}

func NewDBRecordRepository(records *mgo.Collection, imageWriter ResourceWriter, pdfs ResourceWriter, eventService *services.EventService) *DBRecordRepository {
	processedIndex := mgo.Index{
		Key:        []string{"patientId", "-date", "tags", "status"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}

	err := records.EnsureIndex(processedIndex)
	if err != nil {
		log.Panicf("Error setting indices %s", err)
	}

	return &DBRecordRepository{records: records, events: eventService, images: imageWriter, pdfs: pdfs}
}

func (r *DBRecordRepository) All() ([]*models.Record, error) {
	return r.Query(bson.M{})
}

func (r *DBRecordRepository) Find(id string) (*models.Record, error) {
	return r.findByObjectId(bson.ObjectIdHex(id))
}

func (r *DBRecordRepository) findByObjectId(id bson.ObjectId) (*models.Record, error) {
	var record models.Record

	if err := r.records.FindId(id).One(&record); err != nil {
		log.WithField("error", err).Debug("Cannot find record")
		return nil, err
	}
	return &record, nil
}

func (r *DBRecordRepository) Query(query map[string]interface{}) ([]*models.Record, error) {
	records := make([]*models.Record, 0)

	if err := r.records.Find(query).All(&records); err != nil {
		log.Error(err)
		return nil, err
	}

	return records, nil
}

func (r *DBRecordRepository) Create(data models.CreateRecord, images []*shared.Image, pdfData io.Reader) (*models.Record, error) {
	record := models.NewRecord(data)

	if len(record.Pages) == 0 {
		var pages []*models.Page
		for _, img := range images {
			page := models.NewPage(img.Format)

			if err := r.images.Write(NewKeyedGenericResource(img.Image, img.Format, record.Id.Hex(), page.Id)); err != nil {
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

	if err := r.pdfs.Write(NewKeyedGenericResource(pdfBytes, "pdf", record.Id.Hex())); err != nil {
		e := r.images.Delete(NewDirectoryResource(record.Id.Hex()))
		log.Error(e)
		return nil, err
	}

	if err := r.records.Insert(&record); err != nil {
		e := r.images.Delete(NewDirectoryResource(record.Id.Hex()))
		log.Error(e)
		e = r.pdfs.Delete(NewDirectoryResource(record.Id.Hex()))
		log.Error(e)
		return nil, err
	}
	created, err := r.findByObjectId(record.Id)
	if err != nil {
		return nil, err
	}

	r.events.Send(services.EventCreated, created)
	return created, nil
}

func (r *DBRecordRepository) Delete(id string) error {
	err := r.images.Delete(NewDirectoryResource(id))
	if err != nil {
		return err
	}

	err = r.pdfs.Delete(NewDirectoryResource(id))
	if err != nil {
		return err
	}

	key := bson.ObjectIdHex(id)
	err = r.records.RemoveId(key)

	r.events.Send(services.EventDeleted, id)
	return err
}

func (r *DBRecordRepository) Update(id string, record models.Record) (*models.Record, error) {
	key := bson.ObjectIdHex(id)
	// TODO: Remove delted pages from the file system
	if err := r.records.UpdateId(key, bson.M{"$set": record}); err != nil {
		return nil, err
	}
	updated, err := r.findByObjectId(key)
	if err != nil {
		return nil, err
	}
	r.events.Send(services.EventUpdated, updated)
	return updated, nil
}

// UpdatePages updates the pages specified while keeping the rest of the original pages
func (r *DBRecordRepository) UpdatePages(id string, updates []*models.PageUpdate) (*models.Record, error) {
	record, err := r.Find(id)
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
	return r.Update(id, models.Record{Pages: updated})
}
