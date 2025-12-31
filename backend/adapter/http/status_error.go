// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package httpadapter

import (
	"errors"
	"net/http"

	"github.com/FlyingOnion/consee/backend/buffer"
	"github.com/FlyingOnion/consee/backend/service"
)

var (
	errUnknown           = errors.New("unknown error")
	errTokenEmpty        = errors.New("token is empty")
	errInvalidToken      = errors.New("invalid token")
	errIdEmpty           = errors.New("id is empty")
	errInvalidFile       = errors.New("invalid file")
	errInvalidFileFormat = errors.New("invalid file format")
)

func unknownError() *StatusError {
	return &StatusError{
		Status: http.StatusInternalServerError,
		Err:    errUnknown,
	}
}

type StatusError struct {
	Process string
	Status  int
	Err     error
}

func (e StatusError) Error() string {
	process := e.Process
	if process == "" {
		process = "processing request"
	}
	var b buffer.Buffer
	return b.WriteString("an error occurred while ").WriteString(process).
		WriteString(": ").
		WriteString(e.Err.Error()).
		WriteString(" (status: ").WriteInt(e.Status).WriteString(")").
		String()
}

func NewStatusError(err error) *StatusError {
	if statusErr, ok := err.(*StatusError); ok {
		return statusErr
	}
	if statusError, ok := err.(StatusError); ok {
		return &statusError
	}
	if derr, ok := err.(*service.DomainError); ok {
		switch derr.Code {
		case service.DomainErrorCodeNotImplemented:
			return &StatusError{Err: err, Status: http.StatusNotFound}
		case service.DomainErrorCodeAlreadyExists:
			return &StatusError{Err: err, Status: http.StatusConflict}
		case service.DomainErrorCodeNotFound:
			return &StatusError{Err: err, Status: http.StatusNotFound}
		case service.DomainErrorCodeInvalidInput:
			return &StatusError{Err: err, Status: http.StatusBadRequest}
		case service.DomainErrorCodePermissionDenied:
			return &StatusError{Err: err, Status: http.StatusForbidden}
		case service.DomainErrorCodeInternalError:
			return &StatusError{Err: err, Status: http.StatusInternalServerError}
		}
	}
	return unknownError()
}
