package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"net/http"
)

func Run(pagesPath string) {
	router := gin.Default()
	router.Use(cors.Default())
	registerRecords(router.Group("/records"), pagesPath)
	router.Run()
}

func registerStatic(g *gin.RouterGroup, staticDir string) {
	g.StaticFS("/static/", http.Dir(staticDir))
}
