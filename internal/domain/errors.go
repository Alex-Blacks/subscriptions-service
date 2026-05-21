package domain

import "errors"

var (
	ErrAlreadyExists    = errors.New("already exists")
	ErrNotFound         = errors.New("not found")
	ErrConflict         = errors.New("conflict")
	ErrNoFieldsToUpdate = errors.New("no fields to update")
)
