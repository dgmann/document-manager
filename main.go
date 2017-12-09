package main

import (
	"github.com/dgmann/document-manager-api/http"
	"os"
)

func main() {
	recordDir := envOrDefault("RECORD_DIR", "./records")
	http.Run(recordDir)
}

func envOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}
