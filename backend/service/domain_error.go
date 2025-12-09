// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package service

type DomainErrorCode string

const (
	DomainErrorCodeNotImplemented   DomainErrorCode = "NOT_IMPLEMENTED"
	DomainErrorCodeAlreadyExists    DomainErrorCode = "ALREADY_EXISTS"
	DomainErrorCodeNotFound         DomainErrorCode = "NOT_FOUND"
	DomainErrorCodeInvalidInput     DomainErrorCode = "INVALID_INPUT"
	DomainErrorCodeInternalError    DomainErrorCode = "INTERNAL_ERROR"
	DomainErrorCodePermissionDenied DomainErrorCode = "PERMISSION_DENIED"
	DomainErrorCodeMultiple         DomainErrorCode = "MULTIPLE_ERRORS_OCCURED"
	DomainErrorCodeUnknown          DomainErrorCode = "UNKNOWN"
)

type DomainError struct {
	Code    DomainErrorCode
	Message string
}

func (e DomainError) Error() string {
	return e.Message
}

var (
	errNotImplemented        = &DomainError{Code: DomainErrorCodeNotImplemented, Message: "service function not implemented"}
	errFailedToConnectConsul = &DomainError{Code: DomainErrorCodeInternalError, Message: "failed to connect to consul"}
	errPermissionDenied      = &DomainError{Code: DomainErrorCodePermissionDenied, Message: "permission denied"}
	errAdminPermissionDenied = &DomainError{Code: DomainErrorCodeInternalError, Message: "permission denied"}
	errFailedToParse         = &DomainError{Code: DomainErrorCodeInternalError, Message: "failed to parse value"}
	errNotAdmin              = &DomainError{Code: DomainErrorCodePermissionDenied, Message: "token should have admin permission"}
	errUnknown               = &DomainError{Code: DomainErrorCodeUnknown, Message: "unknown error"}
)
