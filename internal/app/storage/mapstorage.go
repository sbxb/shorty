package storage

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
)

// MapStorage defines a simple in-memory storage implemented as a wrapper
// aroung Go map
type MapStorage struct {
	data map[string]string
	file *os.File
	mu   sync.RWMutex
}

// MapStorage implements Storage interface
var _ Storage = (*MapStorage)(nil)

func NewMapStorage() *MapStorage {
	d := make(map[string]string)
	return &MapStorage{data: d}
}

// Open creates a file if missing, opens the file for reading and writing,
// and puts the file object into .file field
func (st *MapStorage) Open(filename string) error {
	if filename == "" {
		return nil
	}
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		return err
	}
	st.file = f
	st.tryLoadRecords()
	return nil
}

// tryLoadRecords tries to load the content of the opened file ignoring any errors
func (st *MapStorage) tryLoadRecords() {
	scanner := bufio.NewScanner(st.file)
	for scanner.Scan() {
		input := strings.Fields(scanner.Text())
		if len(input) != 2 {
			continue
		}
		st.AddURL(input[1], input[0])
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func (st *MapStorage) dumpData() {
	keys := make([]string, len(st.data))
	for k := range st.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Println("*****")
	for _, k := range keys {
		fmt.Println(k, st.data[k])
	}
	fmt.Println("*****")
}

// AddURL saves both url and its id
// MapStorage implementation never returns non-nil error
func (st *MapStorage) AddURL(url string, id string) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.data[id] = url

	if st.file != nil {
		if _, err := st.file.WriteString(fmt.Sprintf("%s\t%s\n", id, url)); err != nil {
			log.Println("MapStorage failed to write to a file", err)
		}
	}

	return nil
}

// GetURL searches for url by its id
// Returns url found or an empty string for a nonexistent id (valid url is
// never an empty string)
// MapStorage implementation never returns non-nil error
func (st *MapStorage) GetURL(id string) (string, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()
	url := st.data[id]

	return url, nil
}

func (st *MapStorage) Close() {
	log.Println("MapStorage closer activated")
	if st.file == nil {
		return
	}
	log.Println("MapStorage closing", st.file.Name())
	if err := st.file.Close(); err != nil {
		log.Println(err)
	}
}
