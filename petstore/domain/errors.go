package domain

import "errors"

var (
	// Err400BadRequest variable
	Err400BadRequest = errors.New("Requested Parameter or Body Not Valid")
	// Err404NotFound variable
	Err404NotFound = errors.New("Requested Resource Not Found")
	// Err500InternalServerError variable
	Err500InternalServerError = errors.New("Internal Server Error")
)
