package main

import (
	"fmt"
	"github.com/dgmann/document-manager/pdf-processor/api"
	"github.com/dgmann/document-manager/pdf-processor/imagick"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	processor := imagick.NewProcessor()
	go startGRPC(processor, processor)

	router := gin.Default()
	router.Use(cors.Default())
	router.POST("images/convert", func(c *gin.Context) {
		images, err := processor.ToImages(c.Request.Body)
		defer c.Request.Body.Close()
		if err != nil {
			c.Status(400)
			c.Error(err)
			c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
			log.Error(err)
			return
		}
		c.JSON(200, images)
	})
	router.POST("images/rotate/:degree", func(c *gin.Context) {
		degree, err := strconv.ParseFloat(c.Param("degree"), 64)
		if err != nil {
			c.AbortWithError(400, err)
		}

		img, err := ioutil.ReadAll(c.Request.Body)
		defer c.Request.Body.Close()
		if err != nil || len(img) == 0 {
			c.AbortWithError(400, err)
		}
		rotated, err := processor.Rotate(img, degree)
		if err != nil {
			c.AbortWithError(400, err)
		}
		c.JSON(200, rotated)
	})
	router.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "PDFProcessor")
	})
	router.Run(":8181")
}

func startGRPC(converter api.PdfToImageConverter, rotator api.Rotator) {
	port := 9000
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	api.RegisterPdfProcessorServer(grpcServer, api.NewGRPCServer(converter, rotator))
	grpcServer.Serve(lis)
}
