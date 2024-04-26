package storage

import "context"

type ArchiveService interface {
	Get(ctx context.Context, id string) (KeyedResource, error)
	ResourceWriter
	ResourceLocator
}
