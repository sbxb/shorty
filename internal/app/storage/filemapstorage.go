package storage

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sbxb/shorty/internal/app/url"
)

// FileMapStorage defines a persistent in-memory storage that loads / saves data
// from / to a file during storage construction / close
type FileMapStorage struct {
	*MapStorage

	file *os.File
}

// TODO ??? use json serialized strings instead of tab separated ones: the current
// format relies on the order of fields and the parsing rule is not flexible at all

// FileMapStorage implements Storage interface
var _ Storage = (*FileMapStorage)(nil)

func NewFileMapStorage(filename string) (*FileMapStorage, error) {
	ms, _ := NewMapStorage()
	if filename == "" {
		return &FileMapStorage{MapStorage: ms}, nil
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		return nil, err
	}

	storage := &FileMapStorage{MapStorage: ms, file: f}
	if err := storage.LoadRecordsFromFile(); err != nil {
		return nil, err
	}

	return storage, nil
}

// tryLoadRecords tries to load the content of the opened file ignoring any errors
func (st *FileMapStorage) LoadRecordsFromFile() error {
	scanner := bufio.NewScanner(st.file)
	for scanner.Scan() {
		input := strings.Fields(scanner.Text())
		if len(input) != 3 {
			continue
		}
		ue := url.URLEntry{
			ShortURL:    input[0],
			OriginalURL: input[2],
		}
		userID := input[1]
		st.AddURL(context.Background(), ue, userID)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (st *FileMapStorage) Close() error {
	if st.file == nil {
		return nil
	}

	if err := st.SaveRecordsToFile(); err != nil {
		return fmt.Errorf("FileMapStorage failed to save data to %s: %v",
			st.file.Name(), err)
	}

	log.Println("MapStorage closing", st.file.Name())

	if err := st.file.Close(); err != nil {
		return err
	}

	st.file = nil

	return nil
}

func (st *FileMapStorage) SaveRecordsToFile() error {
	// Will catch every possible error to make sure the data properly written
	// and the file is in a consistent state
	if err := st.file.Truncate(0); err != nil {
		return err
	}

	if _, err := st.file.Seek(0, 0); err != nil {
		return err
	}

	var wErr error = nil
	w := bufio.NewWriter(st.file)
	for id, uu := range st.data {
		_, wErr = w.WriteString(fmt.Sprintf("%s\t%s\n", id, strings.Replace(uu, ";", "\t", 1)))
		if wErr != nil {
			break
		}
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return wErr
}
