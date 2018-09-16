package main

import (
	"github.com/dgmann/document-manager/m1-adapter/m1"
	"github.com/gin-gonic/gin"
	"os"
)

var dsn string

func init() {
	dsn = getEnv("DSN", "")
	if len(dsn) == 0 {
		panic("invalid connection string: " + dsn)
	}
}

func main() {
	adapter := m1.NewDatabaseAdapter(dsn)
	err := adapter.Connect()
	if err != nil {
		println(err)
	}
	defer adapter.Close()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "M1-Adapter",
		})
	})

	r.GET("/patients", func(c *gin.Context) {
		pats, err := adapter.GetAllPatients()
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, pats)
	})

	r.GET("/patients/:id", func(c *gin.Context) {
		patId := c.Param("id")
		pat, err := adapter.GetPatient(patId)
		if err != nil {
			c.AbortWithStatusJSON(404, gin.H{
				"error":   err.Error(),
				"message": "Patient not found",
			})
			return
		}
		c.JSON(200, pat)
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
