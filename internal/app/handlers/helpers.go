package handlers

import (
	"context"
	"errors"

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

// ConcurrentDeleteBatch takes a slice of ids to be deleted, process the
// slice chunk by chunk starting several (this number is limited by
// concurrentWorkers constant) concurrent workers that call Storage.DeleteBatch()
// Errors are logged and ignored
func ConcurrentDeleteBatch(store storage.Storage, ids []string, userID string) {
	const batchSize = 50        // 1-5 for debug, 50-100+ for testing/production
	const concurrentWorkers = 3 // 2-4 concurrent workers should be enough

	// worker puts {} to reserve a slot, gets {} back when done to free a slot
	pool := make(chan struct{}, concurrentWorkers)

	inputSize := len(ids)
	for beg := 0; beg < inputSize; beg += batchSize {
		end := beg + batchSize
		if end > inputSize {
			end = inputSize
		}

		pool <- struct{}{}
		logger.Debugf("Got permission to process: [%d - %d)", beg, end)

		go func(ids []string) {
			//time.Sleep(3 * time.Second) // makes debug easier
			err := store.DeleteBatch(context.Background(), ids, userID)
			if err != nil {
				logger.Warningf("Deletion worker: %v", err)
			}
			<-pool
		}(ids[beg:end]) // workers read (and only read) non-overlapping parts of the slice
	}
}
