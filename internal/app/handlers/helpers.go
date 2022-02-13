package handlers

import (
	"context"
	"errors"
	"sync"

	"github.com/sbxb/shorty/internal/app/auth"
	"github.com/sbxb/shorty/internal/app/logger"
	"github.com/sbxb/shorty/internal/app/storage"
)

func GetUserID(ctx context.Context) string {
	UserID, _ := ctx.Value(auth.ContextUserIDKey).(string)
	if UserID == "" {
		logger.Warning("User ID not found, check if authCookie middleware was enabled")
	}

	return UserID
}

func IsConflictError(err error) bool {
	var conflictError *storage.IDConflictError

	return errors.As(err, &conflictError)
}

func IsDeletedError(err error) bool {
	var deletedError *storage.URLDeletedError

	return errors.As(err, &deletedError)
}

func ConcurrentDeleteBatch(store storage.Storage, ids []string, userID string) {
	logger.Info("ConcurrentDeleteBatch started")

	inputSize := len(ids)

	// too few ids, delete everything straightforward
	if inputSize < 4 {
		err := store.DeleteBatch(context.Background(), ids, userID)
		if err != nil {
			logger.Warningf("ConcurrentDeleteBatch simple version failed: %v", err)
		} else {
			logger.Info("ConcurrentDeleteBatch simple version successfully completed")
		}
		return
	}

	// Split ids into 4 concurrently processed parts (of almost equal size)
	const numWorkers = 4

	splitChannels := make([]chan string, numWorkers)

	beg := 0
	for i := 0; i < numWorkers; i++ {
		splitChannels[i] = make(chan string)

		chunkSize := inputSize / numWorkers
		if inputSize%numWorkers > i {
			chunkSize++
		}

		go func(out chan<- string, ids []string) {
			for _, id := range ids {
				out <- id
			}
			close(out)
		}(splitChannels[i], ids[beg:beg+chunkSize])

		logger.Debugf("%d:%d", beg, beg+chunkSize)
		beg += chunkSize
	}

	d, _ := NewDeleter(store, userID)
	for id := range fanIn(splitChannels...) {
		logger.Debug(id)

		d.buffer = append(d.buffer, id)
		if cap(d.buffer) == len(d.buffer) {
			if err := d.Flush(); err != nil {
				logger.Warningf("Flush failed: %v", err)
			}
		}
	}
	d.Flush()
}

// fanIn combines multiple input channels into the single output channel
func fanIn(inputChs ...chan string) chan string {
	outCh := make(chan string)

	go func() {
		wg := &sync.WaitGroup{}

		for _, inputCh := range inputChs {
			wg.Add(1)

			go func(inputCh chan string) {
				defer wg.Done()
				for item := range inputCh {
					outCh <- item
				}
			}(inputCh)
		}

		wg.Wait()
		close(outCh)
	}()

	return outCh
}

type Deleter struct {
	buffer []string
	store  storage.Storage
	userID string
}

func NewDeleter(store storage.Storage, userID string) (*Deleter, error) {
	return &Deleter{
		buffer: make([]string, 0, 10), // FIXME 10 for debug; 100 or 1000 for prod
		store:  store,
		userID: userID,
	}, nil
}

func (d *Deleter) Flush() error {
	if len(d.buffer) == 0 {
		return nil
	}
	logger.Debug("Time to flush")
	err := d.store.DeleteBatch(context.Background(), d.buffer, d.userID)
	if err != nil {
		logger.Warningf("Flush failed: %v", err)
		return err
	}
	d.buffer = d.buffer[:0]
	return nil
}
