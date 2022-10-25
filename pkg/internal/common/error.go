package common

import (
	"errors"
	"fmt"
	"runtime"
)

func StringError(err error) error {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	trace := fmt.Sprintf("[%s:%d %s]", frame.File, frame.Line, frame.Function)
	msg := err.Error() + ": " + trace
	res := errors.New(msg)
	return res
}
