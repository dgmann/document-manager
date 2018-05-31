package shared

import "github.com/pkg/errors"

func WrapError(err error, msg string) error {
	if err == nil {
		return errors.New(msg)
	} else {
		return errors.Wrap(err, msg)
	}
}
