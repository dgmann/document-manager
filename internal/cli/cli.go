package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/dgmann/document-manager/pkg/client"
)

var ErrUnknownCommand = errors.New("unkown command")

var dmClient *client.HTTPClient

func Init(args []string) error {
	url := os.Getenv("DM_URL")
	if url == "" {
		return fmt.Errorf("DM_URL unset, example: https://dm.example.com/api")
	}
	client, err := client.NewHTTPClient(url, 0)
	if err != nil {
		return fmt.Errorf("error initializing client: %w", err)
	}
	dmClient = client
	return nil
}
