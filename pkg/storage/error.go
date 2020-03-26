package storage

import "github.com/pkg/errors"

// type DuplicateError struct {
//   error
// }

type NotFoundError struct {
	error
}

// func Duplicate(message string) DuplicateError {
//   return DuplicateError{errors.New(message)}
// }

func NotFound(message string) NotFoundError {
	return NotFoundError{errors.New(message)}
}

// var (
//   DuplicateError = errors.New("duplicate")
//   NotFoundError  = errors.New("not found")
// )
