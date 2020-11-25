package main

import (
	"os"
)

type Config struct {
	Extractor  string
	Rasterizer string
}

func ConfigFromEnv() Config {
	extractor := envOrDefault("EXTRACTOR", "pdfcpu")
	rasterizer := envOrDefault("RASTERIZER", "poppler")
	return Config{
		Extractor:  extractor,
		Rasterizer: rasterizer,
	}
}

func envOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}
