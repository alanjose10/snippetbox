package models

import "errors"

var (
	ErrSnippetNotFound = errors.New("no matching snippet found")

	ErrInvalidCredentials = errors.New("models: invalid credentials")

	ErrDuplicateEmail = errors.New("models: duplicate email")

	ErrUserDoesNotExist = errors.New("models: user does not exist")
)
