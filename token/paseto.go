package token

import (
	"time"

	"github.com/o1egl/paseto"
	"github.com/spf13/viper"
)

var pasetoInstance *paseto.V2 = paseto.NewV2()

func CreateToken(username string, duration time.Duration) (string, error) {
	payload := NewPayload(username, duration)
	return pasetoInstance.Encrypt([]byte(viper.GetString("paseto_key")), payload, nil)
}

func VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := pasetoInstance.Decrypt(token, []byte(viper.GetString("paseto_key")), payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.isValid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
