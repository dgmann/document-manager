package http

import (
	"github.com/dgmann/document-manager-api/shared"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
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

	router.GET("tags", func(context *gin.Context) {
		tags, err := app.Tags.All()
		if err != nil {
			context.AbortWithError(500, err)
			return
		}

		context.JSON(200, tags)
	})

	router.Run()
}
