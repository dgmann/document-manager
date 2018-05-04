package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/dgmann/document-manager/api/services"
	"github.com/dgmann/document-manager/api/repositories"
	"github.com/dgmann/document-manager/api/pdf"
	"github.com/dgmann/document-manager/api/shared"
)

type Factory struct {
	repositories.Factory
	pdfProcessorUrl string
}

func (f *Factory) GetPdfProcessor() *pdf.PDFProcessor {
	return pdf.NewPDFProcessor(f.pdfProcessorUrl)
}

func NewFactory(config *shared.Config) *Factory {
	return &Factory{
		repositories.NewFactory(config),
		config.GetPdfProcessorUrl(),
	}
}

func Run(factory *Factory) {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowMethods("PATCH", "DELETE")
	router.Use(cors.New(config))
	router.Use(location.Default())

	registerWebsocket(router)
	registerRecords(router.Group("/records"), factory)
	registerPatients(router.Group("/patients"), factory)
	registerCategories(router.Group("/categories"), factory)

	router.GET("", func(context *gin.Context) {
		context.String(200, "Document Manager API")
	})

	router.GET("status", func(context *gin.Context) {
		hs := services.GetHealthService()
		if err := hs.Check(); err != nil {
			context.String(500, "Status: Error, Message: %s", err)
		} else {
			context.String(200, "Status: Ok")
		}
	})

	router.GET("tags", func(context *gin.Context) {
		tags, err := factory.GetTagRepository().All()
		if err != nil {
			context.AbortWithError(500, err)
			return
		}

		context.JSON(200, tags)
	})

	router.Run()
}
