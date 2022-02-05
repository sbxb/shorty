package storage

import "fmt"

type IDConflictError struct {
	ID string
}

func (ice *IDConflictError) Error() string {
	return fmt.Sprintf("Storage already has a record with id %s", ice.ID)
}

func NewIDConflictError(id string) error {
	return &IDConflictError{id}
}
