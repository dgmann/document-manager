package directorywatcher

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type Config struct {
	Sources        []Source
	RetryCount     int           `yaml:"retryCount"`
	DestinationUrl string        `yaml:"destinationUrl"`
	ScanInterval   time.Duration `yaml:"scanInterval"`
	Timeout        time.Duration `yaml:"timeout"`
}

type Source struct {
	Parser    string
	Directory string
	Sender    string
}

func LoadConfig(data []byte) (Config, error) {
	config := Config{}

	err := yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		return config, fmt.Errorf("error parsing config file: %w", err)
	}
	return config, err
}
