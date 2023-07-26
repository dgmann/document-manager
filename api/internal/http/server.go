// Package http DocumentManager API.
//
// # The purpose of this application is to provide an api to store uploaded PDFs as images and serve them to various frontends
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//	Schemes: http
//	Host: localhost
//	BasePath: /api
//	Version: 0.0.1
//
//	Consumes:
//	- application/json
//	- multipart/form-data
//
//	Produces:
//	- application/json
//
// swagger:meta
package http

import (
	"context"
	"github.com/dgmann/document-manager/api/internal/datastore"
	"github.com/dgmann/document-manager/api/internal/event"
	"github.com/dgmann/document-manager/api/internal/pdf"
	"github.com/dgmann/document-manager/api/internal/status"
	"github.com/dgmann/document-manager/api/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/handlers"
	"github.com/riandyrn/otelchi"
	"net/http"
)

var (
	PathPrefix = "/api"
)

type Server struct {
	mux    *chi.Mux
	server *http.Server
}

func NewServer(serviceName string, options ...ControllerOption) *Server {
	r := chi.NewRouter()

	r.Use(otelchi.Middleware(serviceName, otelchi.WithChiRoutes(r)))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(handlers.ProxyHeaders)

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, req.URL.String()+"api", http.StatusMovedPermanently)
	})

	r.Get(PathPrefix, func(w http.ResponseWriter, req *http.Request) {
		NewBinaryResponseWithStatus(w, []byte("Document Storage API"), 200).Write()
	})

	r.Mount(PathPrefix+"/debug", middleware.Profiler())

	// Register Controllers
	for _, opt := range options {
		opt(r)
	}

	return &Server{
		mux: r,
	}
}

func (s *Server) Run(port string) error {
	server := &http.Server{Addr: ":" + port, Handler: s.mux}
	s.server = server
	return server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

type ControllerOption = func(rootMux *chi.Mux)

func WithRecordController(recordService datastore.RecordService, imageService storage.ImageService, archiveService storage.ArchiveService, pdfProcessor pdf.Processor) ControllerOption {
	return func(rootMux *chi.Mux) {
		recordController := &RecordController{
			records:      recordService,
			images:       imageService,
			pdfs:         archiveService,
			pdfProcessor: pdfProcessor,
		}
		rootMux.Mount(PathPrefix+"/records", recordController.Router())
	}
}

func WithPatientController(recordService datastore.RecordService, imageService storage.ImageService, categoryService datastore.CategoryService, tagService datastore.TagService) ControllerOption {
	return func(rootMux *chi.Mux) {
		patientController := &PatientController{
			records:    recordService,
			categories: categoryService,
			tags:       tagService,
			images:     imageService,
		}
		rootMux.Mount(PathPrefix+"/patients", patientController.Router())
	}
}

func WithCategoryController(service datastore.CategoryService) ControllerOption {
	return func(rootMux *chi.Mux) {
		categoryController := &CategoryController{
			categories: service,
		}
		rootMux.Mount(PathPrefix+"/categories", categoryController.Router())
	}
}

func WithTagController(tagService datastore.TagService) ControllerOption {
	return func(rootMux *chi.Mux) {
		tagController := &TagController{tags: tagService}
		rootMux.Get(PathPrefix+"/tags", tagController.All)
	}
}

func WithArchiveController(archiveService storage.ResourceLocator) ControllerOption {
	return func(rootMux *chi.Mux) {
		archiveController := &ArchiveController{pdfs: archiveService}
		rootMux.Get(PathPrefix+"/archive/{recordId}", archiveController.One)
	}
}

func WithHealthController(healthService *status.HealthService) ControllerOption {
	return func(rootMux *chi.Mux) {
		healthController := HealthController{healthService: healthService}
		rootMux.Get(PathPrefix+"/status", healthController.Status)
	}
}

func WithStatisticController(statisticsService *status.StatisticsService) ControllerOption {
	return func(rootMux *chi.Mux) {
		statisticsController := StatisticsController{statisticService: statisticsService}
		rootMux.Get(PathPrefix+"/statistics", statisticsController.Statistics)
	}
}

func WithExportController(recordService datastore.RecordService, pdfCreator pdf.Creator) ControllerOption {
	return func(rootMux *chi.Mux) {
		exportController := &ExporterController{creator: pdfCreator, records: recordService}
		rootMux.Get(PathPrefix+"/export", exportController.Export)
	}
}

func WithNotificationController(eventService event.Subscriber) ControllerOption {
	return func(rootMux *chi.Mux) {
		webSocketController := &WebsocketController{Subscriber: eventService}
		rootMux.Mount(PathPrefix+"/notifications", webSocketController.getWebsocketHandler())
	}
}

func WithCounterController(counter datastore.Counter) ControllerOption {
	return func(rootMux *chi.Mux) {
		rootMux.Get(PathPrefix+"/counts/{resource}", func(writer http.ResponseWriter, req *http.Request) {
			resource := URLParamFromContext(req.Context(), "resource")
			count, err := counter.Count(req.Context(), resource)
			if err != nil {
				NewErrorResponse(writer, err, http.StatusInternalServerError).WriteJSON()
				return
			}
			NewResponse(writer, count).WriteJSON()
		})
	}
}
