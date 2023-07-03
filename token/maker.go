package token

import (
	"time"
)

// Maker is interface for managing tokens
type Maker interface {
	// 'CreateToken' will create and sign a new token for spesific username and duration
	CreateToken(username string, duration time.Duration) (string, error)
	// 'VerifyToken' is to checks if the input token is valid or not
	VerifyToken(token string) (*Payload, error)
}
