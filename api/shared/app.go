package shared

import (
	"github.com/globalsign/mgo"
	"github.com/bugsnag/bugsnag-go"
)

type Config struct {
	Db              *mgo.Database
	RecordDir       string
	PDFDir          string
	PdfProcessorUrl string
	BaseUrl         string
	Bugsnag         bugsnag.Configuration
}

type DatabaseConfig interface {
	GetDatabase() *mgo.Database
}

type RecordDirectoryConfig interface {
	GetRecordDirectory() string
}

type PdfProcessorConfig interface {
	GetPdfProcessorUrl() string
}

type HttpConfig interface {
	GetBaseUrl() string
}

type BugsnagConfig interface {
	GetBugsnagConfig() bugsnag.Configuration
}

type PDFDirectoryConfig interface {
	GetPDFDirectory() string
}

func (c *Config) GetDatabase() *mgo.Database {
	return c.Db
}

func (c *Config) GetRecordDirectory() string {
	return c.RecordDir
}

func (c *Config) GetPdfProcessorUrl() string {
	return c.PdfProcessorUrl
}

func (c *Config) GetBaseUrl() string {
	return c.BaseUrl
}

func (c *Config) GetBugsnagConfig() bugsnag.Configuration {
	return c.Bugsnag
}

func (c *Config) GetPDFDirectory() string {
	return c.PDFDir
}
