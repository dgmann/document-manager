package http

import (
	"github.com/bugsnag/bugsnag-go/gin"
	"github.com/dgmann/document-manager/api/pdf"
	"github.com/dgmann/document-manager/api/repositories"
	"github.com/dgmann/document-manager/api/services"
	"github.com/dgmann/document-manager/api/shared"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type Factory struct {
	repositories.Factory
	pdfProcessorUrl string
	baseUrl         string
	eventService    *services.EventService
}

func (f *Factory) GetPdfProcessor() (*pdf.Processor, error) {
	return pdf.NewPDFProcessor(f.pdfProcessorUrl)
}

func (f *Factory) GetResponseService() *services.ResponseService {
	return services.NewResponseService(f.baseUrl, f.GetImageRepository())
}

func (f *Factory) GetEventService() *services.EventService {
	return f.eventService
}

func NewFactory(config *shared.Config) *Factory {
	fileInfoService := repositories.NewImageRepository(config)
	responseService := services.NewResponseService(config.GetBaseUrl(), fileInfoService)
	eventService := services.NewEventService(responseService)
	eventService.Log()
	f := &Factory{
		Factory:         repositories.NewFactory(config, eventService),
		pdfProcessorUrl: config.GetPdfProcessorUrl(),
		baseUrl:         config.GetBaseUrl(),
		eventService:    eventService,
	}
	return f
}

func Run(factory *Factory, c *shared.Config) {
	router := gin.Default()
	pprof.Register(router)
	router.Use(bugsnaggin.AutoNotify(c.GetBugsnagConfig()))
	router.Use(gin.ErrorLogger())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowMethods("PATCH", "DELETE")
	router.Use(cors.New(config))
	router.Use(location.Default())

	registerWebsocket(router, factory.GetEventService())
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

	router.GET("archive/:recordId", func(context *gin.Context) {
		pdfs := factory.GetPDFRepository()
		file, err := pdfs.Get(context.Param("recordId"))
		if err != nil {
			context.AbortWithError(404, err)
			return
		}
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		if err != nil {
			context.AbortWithError(500, err)
			return
		}
		context.Data(200, "application/pdf", data)
	})

	router.Run()
}
