package datastore

import (
	"fmt"
)

type NotFoundError struct {
	Collection string
	Id         string
	Err        error
}

func NewNotFoundError(Id string, Collection string, Error error) *NotFoundError {
	return &NotFoundError{Id: Id, Collection: Collection, Err: Error}
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("could not find ID %s in collection %s", e.Id, e.Collection)
}

func (e *NotFoundError) Unwrap() error {
	return e.Err
}
