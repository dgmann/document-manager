package datastore

import (
	"errors"
	"fmt"
)

type NotFoundError struct {
	Collection string
	Id         string
	Err        error
}

var ErrNoDocuments = errors.New("no results found")

func NewNotFoundError(Id string, Collection string, Error error) *NotFoundError {
	return &NotFoundError{Id: Id, Collection: Collection, Err: Error}
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("could not find ID %s in collection %s", e.Id, e.Collection)
}

func (e *NotFoundError) Is(tgt error) bool {
	target, ok := tgt.(*NotFoundError)
	if !ok {
		return false
	}
	return e.Id == target.Id && e.Collection == target.Collection
}

func (e *NotFoundError) Unwrap() error {
	return e.Err
}
