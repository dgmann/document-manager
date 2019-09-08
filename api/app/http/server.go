package http

import (
	"github.com/dgmann/document-manager/api/app"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/handlers"
	"net/http"
)

type Server struct {
	Healthchecker      map[string]app.HealthChecker
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

	r.Mount("/debug", middleware.Profiler())

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

	r.Mount("/notifications", getWebsocketHandler(s.EventService))
	r.Mount("/records", recordController.Router())
	r.Mount("/patients", patientController.Router())
	r.Mount("/categories", categoryController.Router())

	health := HealthController{s.Healthchecker}
	statistics := StatisticsController{s.StatisticProviders}
	tagController := NewTagController(s.TagService)
	archiveController := NewArchiveController(s.ArchiveService)

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		NewResponseWithStatus(w, []byte("Document Storage API"), 200).Write()
	})
	r.Get("/status", health.Status)
	r.Get("/statistics", statistics.Statistics)

	r.Get("/tags", tagController.All)

	r.Get("/archive/{recordId}", archiveController.One)

	return http.ListenAndServe(":8080", r)
}
