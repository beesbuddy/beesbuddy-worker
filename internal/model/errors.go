package model

import "errors"

var (
	ErrDuplicateEmail     = errors.New("model: duplicate email")
	ErrDuplicateName      = errors.New("model: duplicate name")
	ErrDuplicateUsername  = errors.New("model: duplicate username")
	ErrInvalidCredentials = errors.New("model: invalid credentials")
	ErrInvalidValue       = errors.New("model: invalid value")
	ErrNoRecord           = errors.New("model: no matching record found")
	ErrUnableUpdateRecord = errors.New("model: unable to update record")
	ErrTimestamp          = errors.New("model: invalid timestamp type")
)
