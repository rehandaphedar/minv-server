package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Different types of errors returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("Invalid token")
	ErrExpiredToken = errors.New("Expired token")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) *Payload {
	tokenID, _ := uuid.NewRandom()

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload
}

func (payload *Payload) isValid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
