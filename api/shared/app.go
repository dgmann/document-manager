package shared

import (
	"github.com/globalsign/mgo"
)

type Config struct {
	db              *mgo.Database
	recordDir       string
	pdfProcessorUrl string
}

func NewConfig(db *mgo.Database, recordDir string, pdfProcessorUrl string) *Config {
	return &Config{
		db:              db,
		recordDir:       recordDir,
		pdfProcessorUrl: pdfProcessorUrl,
	}
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

func (c *Config) GetDatabase() *mgo.Database {
	return c.db
}

func (c *Config) GetRecordDirectory() string {
	return c.recordDir
}

func (c *Config) GetPdfProcessorUrl() string {
	return c.pdfProcessorUrl
}
