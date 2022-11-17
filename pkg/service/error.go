package service

import (
	"fmt"

	"github.com/pkg/errors"
)

func PrintError(err error) {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	cause, ok := errors.Cause(err).(stackTracer)
	if !ok {
		panic("error does not implement stackTracer!")
	}

	st := cause.StackTrace()
	fmt.Printf("%+v", st[0:2]) // top two frames
}
