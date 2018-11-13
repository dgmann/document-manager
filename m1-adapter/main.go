package main

import (
	"fmt"
	"github.com/dgmann/document-manager/m1-adapter/m1"
	"github.com/gin-gonic/gin"
	"os"
)

var username string
var password string
var host string
var port string
var dbname string

func init() {
	username = getEnv("DB_USERNAME", "")
	if len(username) == 0 {
		panic("invalid username: " + username)
	}

	password = getEnv("DB_PASSWORD", "")
	if len(password) == 0 {
		panic("invalid password: " + password)
	}

	host = getEnv("DB_HOST", "")
	if len(host) == 0 {
		panic("invalid host: " + host)
	}

	port = getEnv("DB_PORT", "1521")
	if len(port) == 0 {
		panic("invalid port: " + port)
	}

	dbname = getEnv("DB_NAME", "M1DB")
	if len(dbname) == 0 {
		panic("invalid host: " + dbname)
	}
}

func main() {
	dsn := fmt.Sprintf("%s/%s@%s:%s/%s", username, password, host, port, dbname)
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
		firstname := c.Query("firstname")
		lastname := c.Query("lastname")
		var pats []*m1.Patient
		var err error

		if firstname != "" || lastname != "" {
			pats, err = adapter.FindPatientsByName(firstname, lastname)
		} else {
			pats, err = adapter.GetAllPatients()
		}
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
