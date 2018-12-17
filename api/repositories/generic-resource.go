package repositories

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
