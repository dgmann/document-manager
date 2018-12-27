package app

import "io"

type ArchiveService interface {
	Get(id string) (io.Reader, error)
	ResourceWriter
}
