package common

import (
	"github.com/pkg/errors"
)

// func StringError(err error) error {
// 	pc := make([]uintptr, 15)
// 	n := runtime.Callers(2, pc)
// 	frames := runtime.CallersFrames(pc[:n])
// 	frame, _ := frames.Next()
// 	trace := fmt.Sprintf("[%s:%d %s]", frame.File, frame.Line, frame.Function)
// 	msg := err.Error() + ": " + trace
// 	res := errors.New(msg)
// 	return res
// }

// func PrintError(err pkgerrors.Error) {
// 	type stackTracer interface {
// 		StackTrace() pkgerrors.StackTrace
// 	}

// 	err, ok := pkgerrors.Cause(err).(stackTracer)
// 	if !ok {
// 		panic("error does not implement stackTracer!")
// 	}

// 	st := err.StackTrace()
// 	fmt.Printf("%+v", st[0:2]) // top two frames
// }

// func StringError(args ...any) error {
// 	if len(args) == 0 {
// 		return errors.Wrap(errors.New("No Args Provided"), "No Args Provided")
// 	}
// 	err := args[0].(error)
// 	msg := err.Error() + " "
// 	for v := 1; v < len(args); v++ {
// 		msg += args[v].(string) + " "
// 	}
// 	wrapped := errors.Wrap(err, msg)
// 	PrintError(wrapped)
// 	return wrapped
// }

// func PrintError(err error) {
// 	type stackTracer interface {
// 		StackTrace() errors.StackTrace
// 	}

// 	cause, ok := errors.Cause(err).(stackTracer)
// 	if !ok {
// 		panic("ERR DOES NOT IMPLEMENT STACKTRACER")
// 	}
// 	st := cause.StackTrace()
// 	fmt.Printf("%+v", st[0:2]) // top two frames
// }

func StringError(err error, optionalMsg ...string) error {
	concat := ""
	for _, msgs := range optionalMsg {
		concat += msgs + " "
	}
	return errors.Wrap(err, concat)
}
