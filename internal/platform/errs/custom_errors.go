package errs

import "errors"

var ErrNotFound = errors.New("not found")

var ErrBadRequest = errors.New("bad request")

var ErrServerError = errors.New("")

var ErrAlreadyProcessed = errors.New("duplicate")

var ErrInProgress = errors.New("processing not finished")

var ErrInvalidEvent = errors.New("event's fields don't match required types or are missing")
