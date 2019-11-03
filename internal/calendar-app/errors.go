package app

import "errors"

var (
	// ErrNotFound объект не найден
	ErrNotFound = errors.New("object not found")

	// ErrTimeBusy время уже занято
	ErrTimeBusy = errors.New("this time is busy")
)
