package log

import (
	"log/slog"
	"os"
)

func init() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)
}

var Logger = slog.Default()

func ErrAttr(err error) slog.Attr {
	return slog.Any("error", err)
}
