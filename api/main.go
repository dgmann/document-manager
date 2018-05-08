package main

import (
	"github.com/dgmann/document-manager/api/http"
	"github.com/dgmann/document-manager/api/shared"
	"github.com/globalsign/mgo"
	log "github.com/sirupsen/logrus"
	"os"
	"github.com/dgmann/document-manager/api/services"
	"github.com/Shopify/logrus-bugsnag"
	"github.com/bugsnag/bugsnag-go"
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
	dbHost := envOrDefault("DB_HOST", "localhost")
	dbname := envOrDefault("DB_NAME", "manager")
	pdfprocessorUrl := envOrDefault("PDFPROCESSOR_URL", "http://localhost:8181")
	baseUrl := envOrDefault("BASE_URL", "http://localhost:8080")

	session, err := mgo.Dial(dbHost)
	if err != nil {
		log.Errorf("Error connecting to database: %s", err)
	}
	defer session.Close()

	config := &shared.Config{
		Db:              session.DB(dbname),
		RecordDir:       recordDir,
		PdfProcessorUrl: pdfprocessorUrl,
		BaseUrl:         baseUrl,
		Bugsnag:         bugsnagConfig,
	}

	services.InitHealthService(dbHost, pdfprocessorUrl)
	factory := http.NewFactory(config)
	http.Run(factory, config)
}

func envOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}
