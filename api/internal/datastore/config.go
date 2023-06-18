package datastore

import "fmt"

type DatabaseConfig struct {
	Host string
	Port string
	Name string
}

func (config DatabaseConfig) String() string {
	return fmt.Sprintf("%s:%s/%s", config.Host, config.Port, config.Name)
}
