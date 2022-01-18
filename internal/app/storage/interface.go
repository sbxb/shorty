package storage

// Storage provides API for writing/reading URLs to/from a data store
type Storage interface {
	Open(credentials string) error
	AddURL(url string, id string) error
	GetURL(id string) (string, error)
	Close()
}
