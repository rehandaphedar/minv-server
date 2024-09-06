package token

import (
	"os"
	"testing"
	"time"

	"git.sr.ht/~rehandaphedar/minv-server/config"
	"git.sr.ht/~rehandaphedar/minv-server/util"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	config.InitialiseConfig("..")

	returnCode := m.Run()
	os.Exit(returnCode)
}

func TestNormalPASETO(t *testing.T) {
	require := require.New(t)

	username := util.RandomUsername()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := CreateToken(username, duration)
	require.NoError(err)
	require.NotEmpty(token)

	payload, err := VerifyToken(token)
	require.NoError(err)
	require.NotEmpty(payload)

	require.NotZero(payload.ID)
	require.Equal(username, payload.Username)
	require.WithinDuration(issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPASETO(t *testing.T) {
	require := require.New(t)

	token, err := CreateToken(util.RandomUsername(), -time.Minute)
	require.NoError(err)
	require.NotEmpty(token)

	payload, err := VerifyToken(token)
	require.Error(err)
	require.EqualError(err, ErrExpiredToken.Error())
	require.Nil(payload)
}

func TestInvalidPASETO(t *testing.T) {
	require := require.New(t)

	token := util.RandomString(20)

	payload, err := VerifyToken(token)
	require.Error(err)
	require.EqualError(err, ErrInvalidToken.Error())
	require.Nil(payload)
}
