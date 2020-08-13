package model

import (
	"errors"
	"reflect"
)

func checkNonzero(errs []error, value interface{}, message string) []error {
	if reflect.ValueOf(value).IsZero() {
		errs = append(errs, errors.New(message))
	}

	return errs
}
