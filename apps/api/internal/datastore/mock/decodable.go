package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type Decodable struct {
	mock.Mock
}

func NewDecodable() *Decodable {
	return &Decodable{}
}

func (d *Decodable) Decode(data interface{}) error {
	args := d.Called(data)
	return args.Error(0)
}

func (d *Decodable) Err() error {
	args := d.Called()
	return args.Error(0)
}

func (d *Decodable) Close(ctx context.Context) error {
	args := d.Called(ctx)
	return args.Error(0)
}
