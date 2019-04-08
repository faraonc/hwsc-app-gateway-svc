package service

import (
	pbauth "github.com/hwsc-org/hwsc-api-blocks/protobuf/lib"
	"github.com/hwsc-org/hwsc-lib/auth"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeNewAuthSecret(t *testing.T) {
	assert.Nil(t, userSvc.makeNewAuthSecret())
}

func TestGetAuthAuthSecret(t *testing.T) {
	authSecret, err := userSvc.getAuthSecret()
	assert.Nil(t, err)
	assert.NotNil(t, auth.ValidateSecret(authSecret))
}

func TestRefreshCurrAuthSecret(t *testing.T) {
	cases := []struct {
		input    *pbauth.Secret
		isExpErr bool
		err      error
	}{
		{nil, false, nil},
		{expiredAuthSecret, false, nil},
		{validAuthSecret, false, nil},
	}
	for _, c := range cases {
		currAuthSecret = c.input
		err := userSvc.refreshCurrAuthSecret()
		if c.isExpErr {
			assert.EqualError(t, err, c.err.Error())
		} else {
			assert.Nil(t, err)
		}
	}
}
