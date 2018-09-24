package main

import (
	"fmt"
	"github.com/dgmann/document-manager/pdf-processor/converter"
	imagickConverter "github.com/dgmann/document-manager/pdf-processor/converter/imagick"
	"github.com/dgmann/document-manager/pdf-processor/image"
	"github.com/dgmann/document-manager/pdf-processor/processor"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/gographics/imagick.v3/imagick"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

func init() {
	imagick.Initialize()
	defer imagick.Terminate()
}

func main() {
	pdfToImage := imagickConverter.Converter{}
	go startGRPC(pdfToImage)

	router := gin.Default()
	router.Use(cors.Default())
	router.POST("images/convert", func(c *gin.Context) {
		images, err := pdfToImage.ToImages(c.Request.Body)
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
		rotated, err := image.Rotate(img, degree)
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

func startGRPC(converter converter.PdfToImageConverter) {
	port := 9000
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	processor.RegisterPdfProcessorServer(grpcServer, processor.NewGRPCServer(converter))
	grpcServer.Serve(lis)
}
