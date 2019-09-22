package app

import "context"

type Checkable interface {
	Check(context.Context) (string, error)
}
