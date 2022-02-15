package inmemory

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sbxb/shorty/internal/app/logger"
	"github.com/sbxb/shorty/internal/app/storage"
	"github.com/sbxb/shorty/internal/app/url"
)

// FileMapStorage defines a persistent in-memory storage that loads / saves data
// from / to a file during storage construction / close
type FileMapStorage struct {
	*MapStorage

	file *os.File
}

// FileMapStorage implements Storage interface
var _ storage.Storage = (*FileMapStorage)(nil)

func NewFileMapStorage(filename string) (*FileMapStorage, error) {
	ms, _ := NewMapStorage()
	if filename == "" {
		return &FileMapStorage{MapStorage: ms}, nil
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		return nil, err
	}
	logger.Info("FileMapStorage opened", f.Name())
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
		if len(input) != 2 {
			continue
		}
		parts := strings.SplitN(input[1], "|", 3)
		ue := url.URLEntry{
			ShortURL:    input[0],
			OriginalURL: parts[2],
		}
		userID := parts[0]

		st.data[ue.ShortURL] = userID + "|" + parts[1] + "|" + ue.OriginalURL
		logger.Debugf("Loaded from file ==> [%s] :: [%s]", ue.ShortURL, st.data[ue.ShortURL])
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

	logger.Info("FileMapStorage closing", st.file.Name())

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
		_, wErr = w.WriteString(fmt.Sprintf("%s\t%s\n", id, uu))
		if wErr != nil {
			break
		}
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return wErr
}
