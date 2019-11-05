package mongo

type Decoder interface {
	Decode(decodable Decodable, data interface{}) error
}

type Decodable interface {
	Decode(interface{}) error
}

type DefaultDecoder struct {
}

func NewDefaultDecoder() *DefaultDecoder {
	return &DefaultDecoder{}
}

func (d *DefaultDecoder) Decode(decodable Decodable, data interface{}) error {
	return decodable.Decode(data)
}
