package service

type OperationResult int

const (
	Ok OperationResult = iota
	InternalError
	BadRequest
)
