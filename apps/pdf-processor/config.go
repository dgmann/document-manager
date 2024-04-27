package main

import (
	"os"
)

type Config struct {
	Extractor        string
	Rasterizer       string
	OtelCollectorUrl string
}

func ConfigFromEnv() Config {
	extractor := envOrDefault("EXTRACTOR", "pdfcpu")
	rasterizer := envOrDefault("RASTERIZER", "poppler")
	otelCollectorUrl := envOrDefault("OTEL_COLLECTOR_URL", "")
	return Config{
		Extractor:        extractor,
		Rasterizer:       rasterizer,
		OtelCollectorUrl: otelCollectorUrl,
	}
}

func envOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}
