package datastore

import "context"

type Decodable interface {
	Decode(interface{}) error
	Err() error
}

type Cursor interface {
	Next(context.Context) bool
	Decodable

	// Close the cursor.
	Close(context.Context) error
}
