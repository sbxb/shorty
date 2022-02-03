package handlers

import (
	"context"
	"errors"
	"log"

	"github.com/sbxb/shorty/internal/app/auth"
	"github.com/sbxb/shorty/internal/app/storage"
)

func GetUserID(ctx context.Context) string {
	UserID, _ := ctx.Value(auth.ContextUserIDKey).(string)
	if UserID == "" {
		log.Println("User ID not found, check if authCookie middleware was enabled")
	}

	return UserID
}

func IsConflictError(err error) bool {
	var conflictError *storage.IDConflictError
	if errors.As(err, &conflictError) {
		return true
	}

	return false
}
