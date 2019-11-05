package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type Cursor struct {
	mock.Mock
}

func NewCursor() *Cursor {
	return &Cursor{}
}

func (c *Cursor) Next(ctx context.Context) bool {
	args := c.Called(ctx)
	return args.Bool(0)
}

func (c *Cursor) Decode(data interface{}) error {
	args := c.Called(data)
	return args.Error(0)
}

func (c *Cursor) Err() error {
	args := c.Called()
	return args.Error(0)
}

func (c *Cursor) Close(ctx context.Context) error {
	args := c.Called(ctx)
	return args.Error(0)
}
