// Copyright 2020 Ethersphere.
// Copyright 2022 Beezim Authors.
// All rights reserved.
// Use of this source code is originally governed by
// BSD 3-Clause and our modifications by GPLv3
// license that can be found in the LICENSE file.
//
// This code is based on the beekeeper beeclient api.
// The http client was split to its own package.
// The bee api and debug api were modified and
// simplified to fit the purposes of Beezim.

package httpclient

import (
	"errors"
	"strings"
)

// BadRequestError holds list of errors from http response that represent
// invalid data submitted by the user.
type BadRequestError struct {
	errors []string
}

// NewBadRequestError constructs a new BadRequestError with provided errors.
func NewBadRequestError(errors ...string) (err *BadRequestError) {
	return &BadRequestError{
		errors: errors,
	}
}

func (e *BadRequestError) Error() (s string) {
	return strings.Join(e.errors, " ")
}

// Errors returns a list of error messages.
func (e *BadRequestError) Errors() (errs []string) {
	return e.errors
}

// Errors that are returned by the API.
var (
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrNotFound            = errors.New("not found")
	ErrMethodNotAllowed    = errors.New("method not allowed")
	ErrTooManyRequests     = errors.New("too many requests")
	ErrInternalServerError = errors.New("internal server error")
	ErrServiceUnavailable  = errors.New("service unavailable")
	ErrRecoveryInitiated   = errors.New("try again later")
)
