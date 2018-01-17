package http

import (
	"github.com/dgmann/document-manager-api/repositories"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var records *repositories.RecordRepository

func Run(repository *repositories.RecordRepository, pagesPath string) {
	records = repository
	router := gin.Default()
	router.Use(cors.Default())
	registerRecords(router.Group("/records"), pagesPath)
	router.Run()
}

func registerStatic(g *gin.RouterGroup, staticDir string) {
	g.StaticFS("/static/", http.Dir(staticDir))
}
