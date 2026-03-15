package controllers

import "errors"

var (
	ErrNilDependency = errors.New("dependency is nil")
	ErrInvalidConfig = errors.New("invalid configuration")
)
