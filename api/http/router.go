package http

import (
	"github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Run(factory Factory, bug bugsnag.Configuration) {
	router := gin.Default()
	pprof.Register(router)
	router.Use(bugsnaggin.AutoNotify(bug))
	router.Use(gin.ErrorLogger())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowMethods("PATCH", "DELETE")
	router.Use(cors.New(config))
	router.Use(location.New(location.Config{
		Host:             "localhost:8080",
		Scheme:           "http",
		ForwardingHeader: "X-Forwarded-Host",
	}))

	registerWebsocket(router, factory.GetEventService())
	registerRecords(router.Group("/records"), factory)
	registerPatients(router.Group("/patients"), factory)
	registerCategories(router.Group("/categories"), factory)

	generalController := NewGeneralController()
	tagController := NewTagController(factory)
	archiveController := NewArchiveController(factory)

	router.GET("", generalController.Home)
	router.GET("status", generalController.Status)

	router.GET("tags", tagController.All)

	router.GET("archive/:recordId", archiveController.One)

	router.Run()
}
