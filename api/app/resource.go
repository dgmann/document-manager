package app

import "os"

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

type Statter interface {
	Stat(name string) (os.FileInfo, error)
}

type GenericResource struct {
	key    []string
	data   []byte
	format string
}

func NewGenericResource(data []byte, format string) *GenericResource {
	return &GenericResource{[]string{}, data, format}
}

func NewKeyedGenericResource(data []byte, format string, key ...string) *GenericResource {
	return &GenericResource{key, data, format}
}

func NewDirectoryResource(key ...string) *GenericResource {
	return &GenericResource{key: key}
}

func (g *GenericResource) Key() []string {
	return g.key
}

func (g *GenericResource) Data() []byte {
	return g.data
}

func (g *GenericResource) Format() string {
	return g.format
}
