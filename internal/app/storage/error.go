package storage

import "fmt"

// IDConflictError represents "Record already exists" error
type IDConflictError struct {
	ID string
}

func (ice *IDConflictError) Error() string {
	return fmt.Sprintf("Storage already has a record with id %s", ice.ID)
}

func NewIDConflictError(id string) error {
	return &IDConflictError{id}
}

// URLDeletedError represents "Record marked as deleted" error
type URLDeletedError struct {
	ID string
}

func (ude *URLDeletedError) Error() string {
	return fmt.Sprintf("Storage marked as deleted the record with id %s", ude.ID)
}

func NewURLDeletedError(id string) error {
	return &URLDeletedError{id}
}
