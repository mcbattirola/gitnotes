package errflags

import "errors"

// ErrorFlag defines a list of flags you can set on errors.
type ErrorFlag int

const (
	NotFound = iota
	NotAuthorized
	BadParameter
)

// New creates a new error with a flag
func New(text string, flag ErrorFlag) error {
	return Flag(errors.New(text), flag)
}

// Flag wraps err with an error that will return true from HasFlag(err, flag).
func Flag(err error, flag ErrorFlag) error {
	if err == nil {
		return nil
	}
	return flagged{error: err, flag: flag}
}

// HasFlag reports if err has been flagged with the given flag.
func HasFlag(err error, flag ErrorFlag) bool {
	for {
		if f, ok := err.(flagged); ok && f.flag == flag {
			return true
		}
		if err = errors.Unwrap(err); err == nil {
			return false
		}
	}
}

type flagged struct {
	error
	flag ErrorFlag
}

func (f flagged) Unwrap() error {
	return f.error
}
