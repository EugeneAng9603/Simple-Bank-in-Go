package token

import (
	"time"
)

// Maker interface manages tokens
type Maker interface {
	// Creates a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	// Checks if the token is valid
	VerifyToken(token string) (*Payload, error)
}
