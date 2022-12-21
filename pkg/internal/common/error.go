package common

import (
	"github.com/pkg/errors"
)

func StringError(err error, optionalMsg ...string) error {
	if err == nil {
		return nil
	}

	concat := ""

	for _, msgs := range optionalMsg {
		concat += msgs + " "
	}

	if errors.Cause(err) == nil || errors.Cause(err) == err {
		// fmt.Printf("\nWARNING: Error does not implement StackTracer\n")
		return errors.Wrap(errors.New(err.Error()), concat)
	}

	return errors.Wrap(err, concat)
}
