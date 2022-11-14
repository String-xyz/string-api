package common

import (
	"github.com/pkg/errors"
)

func StringError(err error, optionalMsg ...string) error {
	concat := ""
	for _, msgs := range optionalMsg {
		concat += msgs + " "
	}
	return errors.Wrap(err, concat)
}
