package storage

import "fmt"

// Record already exists in storage
type IDConflictError struct {
	ID string
}

func (ice *IDConflictError) Error() string {
	return fmt.Sprintf("Storage already has a record with id %s", ice.ID)
}

func NewIDConflictError(id string) error {
	return &IDConflictError{id}
}

// Record marked as deleted in storage
type URLDeletedError struct {
	ID string
}

func (ude *URLDeletedError) Error() string {
	return fmt.Sprintf("Storage marked as deleted the record with id %s", ude.ID)
}

func NewURLDeletedError(id string) error {
	return &IDConflictError{id}
}
