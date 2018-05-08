package shared

import (
	"github.com/globalsign/mgo"
)

type Config struct {
	Db              *mgo.Database
	RecordDir       string
	PdfProcessorUrl string
	BaseUrl         string
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
