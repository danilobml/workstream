package errs

import "errors"

var ErrNotFound = errors.New("not found")

var ErrBadRequest = errors.New("bad request")

var ErrServerError = errors.New("")

var ErrAlreadyProcessed = errors.New("duplicate")

var ErrInProgress = errors.New("processing not finished")

var ErrInvalidEvent = errors.New("event's fields don't match required types or are missing")

var ErrMailServiceDisabled = errors.New("one or more email config variables are missing")

var ErrAlreadyExists = errors.New("user with this email already exists")

var ErrParsingRoles = errors.New("user with this email already exists")

var ErrInvalidToken = errors.New("invalid user token")

var ErrParsingToken = errors.New("could not parse user token")

var ErrInvalidCredentials = errors.New("invalid credentials")

var ErrUnauthorized = errors.New("unauthorized")

