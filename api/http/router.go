package http

import (
	"github.com/dgmann/document-manager/api/shared"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/dgmann/document-manager/api/services"
)

var app *shared.App

func Run(a *shared.App) {
	app = a
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowMethods("PATCH", "DELETE")
	router.Use(cors.New(config))
	router.Use(location.Default())

	registerWebsocket(router)
	registerRecords(router.Group("/records"))
	registerPatients(router.Group("/patients"))

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
		tags, err := app.Tags.All()
		if err != nil {
			context.AbortWithError(500, err)
			return
		}

		context.JSON(200, tags)
	})

	router.GET("categories", func(context *gin.Context) {
		categories, err := app.Categories.All()
		if err != nil {
			context.AbortWithError(500, err)
			return
		}

		context.JSON(200, categories)
	})

	router.Run()
}
