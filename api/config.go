package main

import (
	"github.com/bugsnag/bugsnag-go"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Config struct {
	Db              *mongo.Database
	RecordDir       string
	PDFDir          string
	PdfProcessorUrl string
	Bugsnag         bugsnag.Configuration
}

func (c *Config) GetDatabase() *mongo.Database {
	return c.Db
}

func (c *Config) GetRecordDirectory() string {
	return c.RecordDir
}

func (c *Config) GetPdfProcessorUrl() string {
	return c.PdfProcessorUrl
}

func (c *Config) GetBugsnagConfig() bugsnag.Configuration {
	return c.Bugsnag
}

func (c *Config) GetPDFDirectory() string {
	return c.PDFDir
}
