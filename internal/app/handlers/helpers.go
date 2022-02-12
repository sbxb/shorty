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
