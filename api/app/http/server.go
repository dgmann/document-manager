package http

import (
	"github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/gin"
	"github.com/dgmann/document-manager/api/app"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Healthchecker   map[string]app.HealthChecker
	EventService    app.Subscriber
	RecordService   app.RecordService
	ImageService    app.ImageService
	CategoryService app.CategoryService
	ArchiveService  app.ArchiveService
	TagService      app.TagService
	PdfProcessor    app.PdfProcessor
	Bug             bugsnag.Configuration
}

func (s *Server) Run() error {
	router := gin.Default()
	pprof.Register(router)
	router.Use(bugsnaggin.AutoNotify(s.Bug))
	router.Use(gin.ErrorLogger())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowMethods("PATCH", "DELETE")
	router.Use(cors.New(corsConfig))
	router.Use(location.New(location.Config{
		Host:             "localhost:8080",
		Scheme:           "http",
		ForwardingHeader: "X-Forwarded-Host",
	}))

	responder := NewResponseFactory(s.ImageService)

	recordController := &RecordController{
		records:         s.RecordService,
		images:          s.ImageService,
		pdfs:            s.ArchiveService,
		pdfProcessor:    s.PdfProcessor,
		responseService: responder,
	}
	patientController := &PatientController{
		records:         s.RecordService,
		categories:      s.CategoryService,
		responseService: responder,
		tags:            s.TagService,
	}
	categoryController := &CategoryController{
		responseService: responder,
		categories:      s.CategoryService,
	}

	registerWebsocket(router, s.EventService)
	registerRecords(router.Group("/records"), recordController)
	registerPatients(router.Group("/patients"), patientController)
	registerCategories(router.Group("/categories"), categoryController)

	health := HealthController{s.Healthchecker}
	tagController := NewTagController(s.TagService)
	archiveController := NewArchiveController(s.ArchiveService)

	router.GET("", func(c *gin.Context) {
		c.String(200, "Document Storage API")
	})
	router.GET("status", health.Status)

	router.GET("tags", tagController.All)

	router.GET("archive/:recordId", archiveController.One)

	return router.Run()
}
