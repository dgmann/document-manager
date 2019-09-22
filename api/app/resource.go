package app

import (
	"os"
	"strings"
)

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
	Formatted
}

type Keyed interface {
	Key() []string
}

type Formatted interface {
	Format() string
}

type Locatable interface {
	Keyed
	Formatted
}

type Writer interface {
	Write(resource KeyedResource) error
}

type Deleter interface {
	Delete(resource Locatable) error
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

//NewKeyedGenericResource creates a new GenericResource with the provided data, format and keys
func NewKeyedGenericResource(data []byte, format string, key ...string) *GenericResource {
	format = strings.TrimLeft(format, ".")
	return &GenericResource{key, data, format}
}

//NewKey creates a new Locatable only with the provided keys set.
//Can be used as an identifier
func NewKey(key ...string) Locatable {
	return &GenericResource{key: key}
}

//NewLocator creates a new Locatable only the provided keys and format.
//Can be used as an identifier
func NewLocator(format string, key ...string) Locatable {
	return &GenericResource{key: key, format: format}
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
