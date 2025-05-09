package domain

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrEmptyAnswer     Error = Error("empty answer")
	ErrIncorrectAnswer Error = Error("incorrect answer")
	ErrSyntaxError     Error = Error("syntax error")
	ErrInternalError   Error = Error("internal error")
)
