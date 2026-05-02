package db

import "time"

// Key is the database representation of an API key.
// SecretHash and Salt are never exposed outside this package.
type Key struct {
	ID         int64
	ClientID   string
	SecretHash string
	Salt       string
	CreatedAt  time.Time
	RevokedAt  *time.Time
}

// Issued is returned once at creation time. The plaintext Secret
// is never stored and cannot be recovered after this point.
type Issued struct {
	ClientID string `json:"client_id"`
	Secret   string `json:"secret"`
}
