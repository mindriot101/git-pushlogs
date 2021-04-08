package errors

import (
	"bytes"
)

const sep = ":\n  "

func Errorf(msg string, err error) *Error {
	return &Error{
		Msg: msg,
		Err: err,
	}
}

// Error wraps all errors that occur
type Error struct {
	// The top level error message
	Msg string
	// The underlying error if any
	Err error
}

func (e *Error) Error() string {
	b := new(bytes.Buffer)
	e.nestedErrorMessages(b)
	return b.String()
}

func (e *Error) nestedErrorMessages(b *bytes.Buffer) {
	pad(b, sep)
	b.WriteString(e.Msg)
	cur := e
	for {
		if cur.Err != nil {
			if prevErr, ok := cur.Err.(*Error); ok {
				cur = prevErr
				cur.nestedErrorMessages(b)
			} else {
				break
			}
		} else {
			break
		}
	}
}

func (e *Error) isZero() bool {
	return e.Err == nil
}

func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}
