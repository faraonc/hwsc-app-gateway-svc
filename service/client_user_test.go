package service

import (
	pbauth "github.com/hwsc-org/hwsc-api-blocks/protobuf/lib"
	"github.com/hwsc-org/hwsc-lib/auth"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeNewAuthSecret(t *testing.T) {
	oldAuthSecret := currAuthSecret
	assert.Nil(t, userSvc.makeNewAuthSecret())
	newAuthSecret, err := userSvc.getAuthSecret()
	assert.Nil(t, err)
	assert.NotEqual(t, oldAuthSecret, newAuthSecret)
}

func TestGetAuthSecret(t *testing.T) {
	assert.Nil(t, userSvc.makeNewAuthSecret())
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
		// this case does not replace the currAuthSecret
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

func TestReplaceCurrAuthSecret(t *testing.T) {
	oldAuthSecret := currAuthSecret
	assert.Nil(t, userSvc.makeNewAuthSecret())
	err := userSvc.replaceCurrAuthSecret()
	assert.Nil(t, err)
	assert.NotEqual(t, oldAuthSecret, currAuthSecret)
}
