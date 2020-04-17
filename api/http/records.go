package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgmann/document-manager/api/datastore"
	"github.com/dgmann/document-manager/api/pdf"
	"github.com/dgmann/document-manager/api/storage"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RecordController struct {
	records      datastore.RecordService
	images       storage.ImageService
	pdfs         storage.ArchiveService
	pdfProcessor pdf.Processor
}

func (controller *RecordController) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", controller.All)
	r.Get("/{recordId}", controller.One)
	r.Post("/", controller.Create)
	r.Patch("/{recordId}", controller.Update)
	r.Delete("/{recordId}", controller.Delete)

	r.Post("/{recordId}/duplicate", controller.Duplicate)
	r.Put("/{recordId}/reset", controller.Reset)
	r.Post("/{recordId}/append/{idtoappend}", controller.Append)

	r.Get("/{recordId}/pages/{imageId}", controller.Page)
	r.Post("/{recordId}/pages", controller.UpdatePages)
	r.Post("/{recordId}/pages/{imageId}/rotate/{degrees}", controller.RotatePage)
	return r
}

// swagger:route GET /records records
//
// Returns all records. Supports filtering by status and pagination
//
// This will show you all available records by default.
// You can filter the records by status and limit the result through pagination.
//
//	Consumes:
//		- application/json
//
//  Produces:
//		- application/json
//
//  Responses:
//		200: DataResponse
//		500: ErrorResponse
func (controller *RecordController) All(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	query := datastore.NewRecordQuery().SetStatus(datastore.Status(params.Get("status")))
	skip, err := strconv.ParseInt(params.Get("skip"), 10, 64)
	if err != nil {
		skip = 0
	}
	limit, err := strconv.ParseInt(params.Get("limit"), 10, 64)
	if err != nil {
		limit = 0
	}

	options := datastore.NewQueryOptions().SetSort(params.Get("sort")).SetSkip(skip).SetLimit(limit)
	records, err := controller.records.Query(req.Context(), query, options)
	if err != nil {
		NewErrorResponse(w, err, http.StatusInternalServerError).WriteJSON()
		return
	}
	withUrl := SetURLForRecordList(records, url.URL{Scheme: req.URL.Scheme, Host: req.Host}, controller.images)
	NewResponse(w, withUrl).WriteJSON()
}

func (controller *RecordController) One(w http.ResponseWriter, req *http.Request) {
	id := URLParamFromContext(req.Context(), "recordId")
	result, err := controller.records.Find(req.Context(), id)
	if err != nil {
		var notFoundErr *datastore.NotFoundError
		if ok := errors.As(err, &notFoundErr); ok {
			NewErrorResponse(w, notFoundErr, http.StatusNotFound).WriteJSON()
		} else {
			NewErrorResponse(w, err, http.StatusInternalServerError).WriteJSON()
		}

		return
	}

	withUrl := SetURLForRecord(result, url.URL{Scheme: req.URL.Scheme, Host: req.Host}, controller.images)
	NewResponse(w, withUrl).WriteJSON()
}

func (controller *RecordController) Create(w http.ResponseWriter, req *http.Request) {
	sender := req.FormValue("sender")
	patientId := req.FormValue("patientId")
	status := req.FormValue("status")
	category := req.FormValue("category")
	receivedAt := time.Now()
	var date time.Time
	if r := req.FormValue("receivedAt"); r != "" {
		if parsed, err := time.Parse(time.RFC3339, r); err != nil {
			receivedAt = parsed
		}
	}
	if r := req.FormValue("date"); r != "" {
		if parsed, err := time.Parse(time.RFC3339, r); err != nil {
			date = parsed
		}
	}
	newRecord := datastore.CreateRecord{Sender: sender, ReceivedAt: receivedAt, Date: date, PatientId: &patientId, Status: datastore.Status(status), Category: &category}

	file, _, err := req.FormFile("pdf")
	if err != nil {
		NewErrorResponse(w, fmt.Errorf("no file found. Please specify a pdf file in the field: pdf. %w", err), 400).WriteJSON()
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
		return
	}
	_ = file.Close()

	images, err := controller.pdfProcessor.Convert(req.Context(), bytes.NewBuffer(fileBytes))
	if err != nil {
		NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
		return
	}

	res, err := controller.records.Create(req.Context(), newRecord, images, bytes.NewBuffer(fileBytes))
	if err != nil {
		NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
		return
	}

	withUrl := SetURLForRecord(res, url.URL{Scheme: req.URL.Scheme, Host: req.Host}, controller.images)
	NewResponseWithStatus(w, withUrl, http.StatusCreated).WriteJSON()
}

func (controller *RecordController) Update(w http.ResponseWriter, req *http.Request) {
	var body datastore.Record

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
		return
	}

	if body.Status != nil && !body.Status.IsValid() {
		NewErrorResponse(w, fmt.Errorf("%s is not a valid status", *body.Status), http.StatusBadRequest).WriteJSON()
		return
	}

	id := URLParamFromContext(req.Context(), "recordId")
	updated, err := controller.records.Update(req.Context(), id, body)
	if err != nil {
		var e *datastore.NotFoundError
		statusCode := http.StatusBadRequest
		if errors.As(err, &e) {
			statusCode = http.StatusNotFound
		}
		NewErrorResponse(w, err, statusCode).WriteJSON()
		return
	}

	withUrl := SetURLForRecord(updated, url.URL{Scheme: req.URL.Scheme, Host: req.Host}, controller.images)
	NewResponse(w, withUrl).WriteJSON()
}

func (controller *RecordController) Delete(w http.ResponseWriter, req *http.Request) {
	id := URLParamFromContext(req.Context(), "recordId")
	err := controller.records.Delete(req.Context(), id)
	if err != nil {
		NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
		return
	}
	NewResponseWithStatus(w, nil, http.StatusNoContent).WriteJSON()
}

func (controller *RecordController) Duplicate(w http.ResponseWriter, req *http.Request) {
	id := URLParamFromContext(req.Context(), "recordId")
	recordToDuplicate, err := controller.records.Find(req.Context(), id)
	if err != nil {
		NewErrorResponse(w, err, http.StatusNotFound).WriteJSON()
		return
	}

	file, err := controller.pdfs.Get(id)
	if err != nil {
		NewErrorResponse(w, err, http.StatusInternalServerError).WriteJSON()
		return
	}

	newId := primitive.NewObjectID()

	err = controller.images.Copy(recordToDuplicate.Id.Hex(), newId.Hex())
	if err != nil {
		NewErrorResponse(w, err, http.StatusInternalServerError).WriteJSON()
		return
	}

	data := bytes.NewBuffer(file.Data())
	copiedRecord, err := controller.records.Create(req.Context(), datastore.CreateRecord{
		Id:         &newId,
		ReceivedAt: recordToDuplicate.ReceivedAt,
		Sender:     recordToDuplicate.Sender,
		Pages:      recordToDuplicate.Pages,
	}, nil, data)

	if err != nil {
		NewErrorResponse(w, err, http.StatusInternalServerError).WriteJSON()
		return
	}

	withUrl := SetURLForRecord(copiedRecord, url.URL{Scheme: req.URL.Scheme, Host: req.Host}, controller.images)
	NewResponseWithStatus(w, withUrl, http.StatusCreated).WriteJSON()
}

func (controller *RecordController) Reset(w http.ResponseWriter, req *http.Request) {
	id := URLParamFromContext(req.Context(), "recordId")
	f, err := controller.pdfs.Get(id)
	if err != nil {
		NewErrorResponse(w, err, http.StatusNotFound).WriteJSON()
		return
	}

	images, err := controller.pdfProcessor.Convert(req.Context(), bytes.NewBuffer(f.Data()))
	if err != nil {
		NewErrorResponse(w, err, http.StatusInternalServerError).WriteJSON()
		return
	}

	var pages []datastore.Page
	for _, img := range images {
		page := datastore.NewPage(img.Format)
		if err := controller.images.Write(storage.NewKeyedGenericResource(img.Image, img.Format, id, page.Id)); err != nil {
			NewErrorResponse(w, err, http.StatusInternalServerError).WriteJSON()
			return
		}
		pages = append(pages, *page)
	}

	updated, err := controller.records.Update(req.Context(), id, datastore.Record{Pages: pages})
	if err != nil {
		NewErrorResponse(w, err, http.StatusInternalServerError).WriteJSON()
		return
	}

	withUrl := SetURLForRecord(updated, url.URL{Scheme: req.URL.Scheme, Host: req.Host}, controller.images)
	NewResponse(w, withUrl).WriteJSON()
}

func (controller *RecordController) Append(w http.ResponseWriter, req *http.Request) {
	idToAppend := URLParamFromContext(req.Context(), "idtoappend")
	recordToAppend, err := controller.records.Find(req.Context(), idToAppend)
	if err != nil {
		NewErrorResponse(w, err, http.StatusNotFound).WriteJSON()
		return
	}

	targetId := URLParamFromContext(req.Context(), "recordId")
	targetRecord, err := controller.records.Find(req.Context(), targetId)
	if err != nil {
		NewErrorResponse(w, err, http.StatusNotFound).WriteJSON()
		return
	}

	pages := append(targetRecord.Pages, recordToAppend.Pages...)

	err = controller.images.Copy(idToAppend, targetId)
	if err != nil {
		NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
		return
	}

	updated, err := controller.records.Update(req.Context(), targetId, datastore.Record{Pages: pages})
	if err != nil {
		NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
		return
	}

	withUrl := SetURLForRecord(updated, url.URL{Scheme: req.URL.Scheme, Host: req.Host}, controller.images)
	NewResponse(w, withUrl).WriteJSON()
}

func (controller *RecordController) Page(w http.ResponseWriter, req *http.Request) {
	id := URLParamFromContext(req.Context(), "recordId")
	rec, err := controller.records.Find(req.Context(), id)
	if err != nil {
		NewErrorResponse(w, err, http.StatusNotFound).WriteJSON()
		return
	}

	imageId := URLParamFromContext(req.Context(), "imageId")
	for _, page := range rec.Pages {
		if page.Id == imageId {
			p := controller.images.Locate(storage.NewLocator(page.Format, id, imageId))
			http.ServeFile(w, req, p)
			return
		}
	}
	NewErrorResponse(w, errors.New("page not found"), 404).WriteJSON()
}

func (controller *RecordController) UpdatePages(w http.ResponseWriter, req *http.Request) {
	var updates []datastore.PageUpdate
	if err := json.NewDecoder(req.Body).Decode(&updates); err != nil {
		NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
		return
	}

	id := URLParamFromContext(req.Context(), "recordId")
	updated, err := controller.records.UpdatePages(req.Context(), id, updates)
	if err != nil {
		NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
		return
	}

	images, err := controller.images.Get(id)
	if err != nil {
		NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
		return
	}

	var errs []error
	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}

	for _, u := range updates {
		if u.Rotate == 0 {
			continue
		}

		wg.Add(1)
		go func(update datastore.PageUpdate) {
			defer wg.Done()

			if img, ok := images[update.Id]; ok {
				img, err := controller.pdfProcessor.Rotate(req.Context(), bytes.NewBuffer(img.Image), int(update.Rotate))
				if err != nil {
					mutex.Lock()
					errs = append(errs, fmt.Errorf("error rotating page %s: %w", update.Id, err))
					mutex.Unlock()
					return
				}
				err = controller.images.Write(storage.NewKeyedGenericResource(img.Image, img.Format, id, update.Id))
				if err != nil {
					mutex.Lock()
					errs = append(errs, fmt.Errorf("error saving page %s: %w", update.Id, err))
					mutex.Unlock()
					return
				}
			} else {
				mutex.Lock()
				errs = append(errs, fmt.Errorf("unknown error processing page %s: %w", update.Id, err))
				mutex.Unlock()
			}
		}(u)
	}

	wg.Wait()
	if len(errs) > 0 {
		errMessages := make([]string, len(errs))
		for i, e := range errs {
			errMessages[i] = e.Error()
		}
		NewErrorResponse(w, fmt.Errorf("error rotating pages: %s", strings.Join(errMessages, ", \n")), http.StatusInternalServerError).WriteJSON()
		return
	}

	withUrl := SetURLForRecord(updated, url.URL{Scheme: req.URL.Scheme, Host: req.Host}, controller.images)
	NewResponse(w, withUrl).WriteJSON()
}

func (controller *RecordController) RotatePage(w http.ResponseWriter, req *http.Request) {
	id := URLParamFromContext(req.Context(), "recordId")
	images, err := controller.images.Get(id)
	if err != nil {
		NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
		return
	}

	degrees, err := strconv.Atoi(URLParamFromContext(req.Context(), "degrees"))
	if err != nil {
		NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
		return
	}

	imageId := URLParamFromContext(req.Context(), "imageId")
	if img, ok := images[imageId]; ok {
		img, err := controller.pdfProcessor.Rotate(req.Context(), bytes.NewBuffer(img.Image), degrees)
		if err != nil {
			NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
			return
		}
		if err := controller.images.Write(storage.NewKeyedGenericResource(img.Image, img.Format, id, imageId)); err != nil {
			NewErrorResponse(w, err, http.StatusInternalServerError).WriteJSON()
			return
		}
		NewResponse(w, img)
	} else {
		NewErrorResponse(w, fmt.Errorf("image %s not found", imageId), http.StatusNotFound).WriteJSON()
		return
	}
}
