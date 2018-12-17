package main

import (
	"github.com/Shopify/logrus-bugsnag"
	"github.com/bugsnag/bugsnag-go"
	"github.com/dgmann/document-manager/api/http"
	"github.com/dgmann/document-manager/api/services"
	"github.com/globalsign/mgo"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

var bugsnagConfig bugsnag.Configuration

func init() {
	apiKey := envOrDefault("BUGSNAG_API_KEY", "")
	stage := envOrDefault("RELEASE_STAGE", "production")
	bugsnagConfig = bugsnag.Configuration{
		// Your Bugsnag project API key
		APIKey:       apiKey,
		ReleaseStage: stage,
		// The import paths for the Go packages
		// containing your source files
		ProjectPackages: []string{"main", "github.com/dgmann/document-manager/api"},
		PanicHandler:    func() {},
	}
	bugsnag.Configure(bugsnagConfig)

	hook, err := logrus_bugsnag.NewBugsnagHook()
	if err != nil {
		log.Error("Error registering bugsnag hook")
	}
	log.StandardLogger().Hooks.Add(hook)
}

func main() {
	recordDir := envOrDefault("RECORD_DIR", "/records")
	archiveDir := envOrDefault("ARCHIVE_DIR", "/archive")
	dbHost := envOrDefault("DB_HOST", "localhost")
	dbname := envOrDefault("DB_NAME", "manager")
	pdfprocessorUrl := envOrDefault("PDFPROCESSOR_URL", "127.0.0.1:9000")

	session, err := mgo.Dial(dbHost)
	if err != nil {
		log.Errorf("Error connecting to database: %s", err)
		os.Exit(1)
		return
	}
	defer session.Close()

	config := &Config{
		Db:              session.DB(dbname),
		RecordDir:       recordDir,
		PDFDir:          archiveDir,
		PdfProcessorUrl: pdfprocessorUrl,
		Bugsnag:         bugsnagConfig,
	}

	services.InitHealthService(dbHost, pdfprocessorUrl)
	factory := NewFactory(config)
	http.Run(factory, bugsnagConfig)
}

func envOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}
