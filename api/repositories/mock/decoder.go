package mock

import (
	"github.com/dgmann/document-manager/api/repositories/database"
	"github.com/stretchr/testify/mock"
)

type Decoder struct {
	mock.Mock
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) Decode(decodable database.Decodable, data interface{}) error {
	args := d.Called(decodable, data)
	return args.Error(0)
}
