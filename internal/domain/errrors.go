package domain

import "errors"

var (
	ErrEmptyAnswer     = errors.New("empty answer")
	ErrIncorrectAnswer = errors.New("incorrect answer")
	ErrSyntaxError     = errors.New("syntax error")
	ErrInternalError   = errors.New("internal error")
)
