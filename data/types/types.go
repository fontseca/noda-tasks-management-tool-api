package types

import (
	"github.com/google/uuid"
	"time"
)

// Role represents the role of a user.
type Role uint8

const (
	// RoleAdmin represents the administrator role.
	RoleAdmin Role = 1 + iota
	// RoleUser represents the regular user role.
	RoleUser
)

// TaskPriority represents the priority level of a task.
type TaskPriority string

const (
	// TaskPriorityHigh represents a high-priority task.
	TaskPriorityHigh TaskPriority = "high"
	// TaskPriorityMedium represents a medium-priority task.
	TaskPriorityMedium TaskPriority = "medium"
	// TaskPriorityLow represents a low-priority task.
	TaskPriorityLow TaskPriority = "low"
)

// TaskStatus represents the status of a task.
type TaskStatus string

const (
	// TaskStatusIncomplete represents an incomplete task.
	TaskStatusIncomplete TaskStatus = "incomplete"
	// TaskStatusComplete represents a complete task.
	TaskStatusComplete TaskStatus = "complete"
	// TaskStatusDeferred represents a deferred task.
	TaskStatusDeferred TaskStatus = "deferred"
)

// Position represents a position in a sequence.
type Position uint32

// ContextKey is used as a key for context values.
type ContextKey struct{}

type JWTPayload struct {
	UserID   uuid.UUID // UserID is the unique identifier for a user.
	UserRole Role      // UserRole represents the role of the user.
}

// TokenExpires represents the expiration details of a token.
type TokenExpires struct {
	At     time.Time `json:"at"`     // At is the absolute expiration time.
	Within float64   `json:"within"` // Within is the relative expiration time.
	Unit   string    `json:"unit"`   // Unit is the unit of time for the relative expiration.
}

// TokenPayload represents the payload of a token.
type TokenPayload struct {
	ID       string       `json:"token_id"`  // ID is the unique identifier for the token.
	Token    string       `json:"token"`     // Token is the actual token value.
	Subject  string       `json:"subject"`   // Subject is the subject of the token.
	Issuer   string       `json:"issuer"`    // Issuer is the entity that issued the token.
	IssuedAt time.Time    `json:"issued_at"` // IssuedAt is the time when the token was issued.
	Expires  TokenExpires `json:"expires"`   // Expires represents the expiration details of the token.
}

// Pagination represents the pagination parameters for querying a collection.
type Pagination struct {
	Page int64 // Page is the page number.
	RPP  int64 // RPP (Records Per Page) is the number of results to be returned per page.
}

// Result represents the result of a paginated query.
type Result[T any] struct {
	Page      int64 `json:"page"`      // Page is the current page number.
	RPP       int64 `json:"rpp"`       // RPP is the number of Records Per Page for the query.
	Retrieved int64 `json:"retrieved"` // Retrieved is the total number of items retrieved.
	Payload   []*T  `json:"payload"`   // Payload is the actual data payload.
}
