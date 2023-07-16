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
	Port              string
	HealthService     *status.HealthService
	StatisticsService *status.StatisticsService
	EventService      event.Subscriber
	RecordService     datastore.RecordService
	ImageService      storage.ImageService
	CategoryService   datastore.CategoryService
	ArchiveService    storage.ArchiveService
	TagService        datastore.TagService
	PdfProcessor      pdf.Processor
	ServiceName       string
	server            *http.Server
}

func (s *Server) Run() error {
	r := chi.NewRouter()

	r.Use(otelchi.Middleware(s.ServiceName, otelchi.WithChiRoutes(r)))
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

	recordController := &RecordController{
		records:      s.RecordService,
		images:       s.ImageService,
		pdfs:         s.ArchiveService,
		pdfProcessor: s.PdfProcessor,
	}
	patientController := &PatientController{
		records:    s.RecordService,
		categories: s.CategoryService,
		tags:       s.TagService,
		images:     s.ImageService,
	}
	categoryController := &CategoryController{
		categories: s.CategoryService,
	}

	health := HealthController{s.HealthService}
	statistics := StatisticsController{s.StatisticsService}
	tagController := NewTagController(s.TagService)
	archiveController := NewArchiveController(s.ArchiveService)
	exportController := NewExporterController(s.PdfProcessor, s.RecordService)

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, req.URL.String()+"api", http.StatusMovedPermanently)
	})

	r.Get(PathPrefix, func(w http.ResponseWriter, req *http.Request) {
		NewBinaryResponseWithStatus(w, []byte("Document Storage API"), 200).Write()
	})

	r.Mount(PathPrefix+"/notifications", getWebsocketHandler(s.EventService))
	r.Mount(PathPrefix+"/records", recordController.Router())
	r.Mount(PathPrefix+"/patients", patientController.Router())
	r.Mount(PathPrefix+"/categories", categoryController.Router())
	r.Get(PathPrefix+"/tags", tagController.All)
	r.Get(PathPrefix+"/archive/{recordId}", archiveController.One)
	r.Get(PathPrefix+"/export", exportController.Export)

	r.Mount(PathPrefix+"/debug", middleware.Profiler())
	r.Get(PathPrefix+"/status", health.Status)
	r.Get(PathPrefix+"/statistics", statistics.Statistics)

	server := &http.Server{Addr: ":" + s.Port, Handler: r}
	s.server = server
	return server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
