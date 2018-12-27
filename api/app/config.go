package app

import (
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Config struct {
	Db               *mongo.Database
	RecordDirectory  string
	ArchiveDirectory string
	PdfProcessorUrl  string
}
