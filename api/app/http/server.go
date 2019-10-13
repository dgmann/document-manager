package http

import (
	"github.com/dgmann/document-manager/api/app"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/handlers"
	"net/http"
)

var (
    PathPrefix = "/api"
)

type Server struct {
	Healthchecker      map[string]app.Checkable
	StatisticProviders map[string]app.StatisticProvider
	EventService       app.Subscriber
	RecordService      app.RecordService
	ImageService       app.ImageService
	CategoryService    app.CategoryService
	ArchiveService     app.ArchiveService
	TagService         app.TagService
	PdfProcessor       app.PdfProcessor
}

func (s *Server) Run() error {
	r := chi.NewRouter()

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

	health := HealthController{s.Healthchecker}
	statistics := StatisticsController{s.StatisticProviders}
	tagController := NewTagController(s.TagService)
	archiveController := NewArchiveController(s.ArchiveService)
	exportController := NewExporterController(s.PdfProcessor, s.RecordService)

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, req.URL.String() + "api", http.StatusMovedPermanently)
	})

	r.Get(PathPrefix, func(w http.ResponseWriter, req *http.Request) {
		NewResponseWithStatus(w, []byte("Document Storage API"), 200).Write()
	})

	r.Mount(PathPrefix + "/notifications", getWebsocketHandler(s.EventService))
	r.Mount(PathPrefix + "/records", recordController.Router())
	r.Mount(PathPrefix + "/patients", patientController.Router())
	r.Mount(PathPrefix + "/categories", categoryController.Router())
	r.Get(PathPrefix + "/tags", tagController.All)
	r.Get(PathPrefix + "/archive/{recordId}", archiveController.One)
	r.Get(PathPrefix + "/export", exportController.Export)

	r.Mount(PathPrefix + "/debug", middleware.Profiler())
	r.Get(PathPrefix + "/status", health.Status)
	r.Get(PathPrefix + "/statistics", statistics.Statistics)

	return http.ListenAndServe(":80", r)
}
