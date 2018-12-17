package repositories

type ResourceWriter interface {
	Deleter
	Writer
}

type KeyedResource interface {
	Keyed
	Resource
}

type Resource interface {
	Data() []byte
	Format() string
}

type Keyed interface {
	Key() []string
}

type Writer interface {
	Write(resource KeyedResource) error
}

type Deleter interface {
	Delete(resource KeyedResource) error
}
