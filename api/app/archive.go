package app

type ArchiveService interface {
	Get(id string) (KeyedResource, error)
	ResourceWriter
}
