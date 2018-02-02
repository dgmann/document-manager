package http

import (
	"github.com/dgmann/document-manager-api/shared"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var app *shared.App

func Run(a *shared.App) {
	app = a
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowMethods("PATCH", "DELETE")
	router.Use(cors.New(config))
	registerRecords(router.Group("/records"))
	router.Run()
}

func registerStatic(g *gin.RouterGroup, staticDir string) {
	g.StaticFS("/static/", http.Dir(staticDir))
}
