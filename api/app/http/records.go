package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgmann/document-manager/api/app"
	"github.com/go-chi/chi"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type RecordController struct {
	records      app.RecordService
	images       app.ImageService
	pdfs         app.ArchiveService
	pdfProcessor app.PdfProcessor
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

func (controller *RecordController) All(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	query := make(map[string]interface{})
	for k, v := range params {
		query[k] = v[0]
	}
	records, err := controller.records.Query(req.Context(), query)
	if err != nil {
		NewErrorResponse(w, err, 400).WriteJSON()
		return
	}
	withUrl := SetURLForRecordList(records, url.URL{Scheme: req.URL.Scheme, Host: req.Host}, controller.images)
	NewResponse(w, withUrl).WriteJSON()
}

func (controller *RecordController) One(w http.ResponseWriter, req *http.Request) {
	id := URLParamFromContext(req.Context(), "recordId")
	result, err := controller.records.Find(req.Context(), id)
	if err != nil {
		NewErrorResponse(w, err, 404)
		return
	}

	withUrl := SetURLForRecord(result, url.URL{Scheme: req.URL.Scheme, Host: req.Host}, controller.images)
	NewResponse(w, withUrl).WriteJSON()
}

func (controller *RecordController) Create(w http.ResponseWriter, req *http.Request) {
	sender := req.FormValue("sender")
	newRecord := app.CreateRecord{Sender: sender}

	_, fh, err := req.FormFile("pdf")
	if err != nil {
		NewErrorResponse(w, fmt.Errorf("no file found. Please specify a pdf file in the field: pdf. %w", err), 400).WriteJSON()
		return
	}

	f, err := fh.Open()
	if err != nil {
		NewErrorResponse(w, fmt.Errorf("coud not open file: %w", err), 400).WriteJSON()
		return
	}
	defer f.Close()

	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
		return
	}

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
	var body app.Record

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
		return
	}

	id := URLParamFromContext(req.Context(), "recordId")
	updated, err := controller.records.Update(req.Context(), id, body)
	if err != nil {
		var e *app.NotFoundError
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

	copiedRecord, err := controller.records.Create(req.Context(), app.CreateRecord{
		Id:         &newId,
		ReceivedAt: recordToDuplicate.ReceivedAt,
		Sender:     recordToDuplicate.Sender,
		Pages:      recordToDuplicate.Pages,
	}, nil, file)

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

	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		NewErrorResponse(w, err, http.StatusInternalServerError).WriteJSON()
		return
	}

	images, err := controller.pdfProcessor.Convert(req.Context(), bytes.NewBuffer(fileBytes))
	if err != nil {
		NewErrorResponse(w, err, http.StatusInternalServerError).WriteJSON()
		return
	}

	var pages []app.Page
	for _, img := range images {
		page := app.NewPage(img.Format)
		if err := controller.images.Write(app.NewKeyedGenericResource(img.Image, img.Format, id, page.Id)); err != nil {
			NewErrorResponse(w, err, http.StatusInternalServerError).WriteJSON()
			return
		}
		pages = append(pages, *page)
	}

	updated, err := controller.records.Update(req.Context(), id, app.Record{Pages: pages})
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

	updated, err := controller.records.Update(req.Context(), targetId, app.Record{Pages: pages})
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
			p := controller.images.Path(id, imageId, page.Format)
			http.ServeFile(w, req, p)
			return
		}
	}
	NewErrorResponse(w, errors.New("page not found"), 404).WriteJSON()
}

func (controller *RecordController) UpdatePages(w http.ResponseWriter, req *http.Request) {
	var updates []app.PageUpdate
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

	var errorIds []string
	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}

	for _, u := range updates {
		if u.Rotate == 0 {
			continue
		}

		wg.Add(1)
		go func(update app.PageUpdate) {
			defer wg.Done()

			if img, ok := images[update.Id]; ok {
				img, err := controller.pdfProcessor.Rotate(req.Context(), bytes.NewBuffer(img.Image), int(update.Rotate))
				if err != nil {
					mutex.Lock()
					errorIds = append(errorIds, update.Id)
					mutex.Unlock()
					return
				}
				err = controller.images.Write(app.NewKeyedGenericResource(img.Image, img.Format, id, update.Id))
				if err != nil {
					mutex.Lock()
					errorIds = append(errorIds, update.Id)
					mutex.Unlock()
					return
				}
			} else {
				mutex.Lock()
				errorIds = append(errorIds, update.Id)
				mutex.Unlock()
			}
		}(u)
	}

	wg.Wait()
	if len(errorIds) > 0 {
		NewErrorResponse(w, fmt.Errorf("error rotating pages %s", strings.Join(errorIds, ",")), http.StatusInternalServerError).WriteJSON()
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
		if err := controller.images.Write(app.NewKeyedGenericResource(img.Image, img.Format, id, imageId)); err != nil {
			NewErrorResponse(w, err, http.StatusInternalServerError).WriteJSON()
			return
		}
		NewResponse(w, img)
	} else {
		NewErrorResponse(w, fmt.Errorf("image %s not found", imageId), http.StatusNotFound).WriteJSON()
		return
	}
}
