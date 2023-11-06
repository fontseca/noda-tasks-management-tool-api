package types

import (
	"github.com/google/uuid"
)

type ContextKey struct{}

type JWTPayload struct {
	UserID   uuid.UUID
	UserRole Role
}
