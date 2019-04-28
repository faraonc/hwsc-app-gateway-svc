package service

import (
	"github.com/Pallinder/go-randomdata"
	pblib "github.com/hwsc-org/hwsc-api-blocks/protobuf/lib"
	"github.com/hwsc-org/hwsc-lib/auth"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_userDial(t *testing.T) {
	assert.Nil(t, userSvc.userSvcConn.Close(), "test closing connection")
	assert.Nil(t, userSvc.dial(), "test dialing with error")
}

func Test_userGetConnection(t *testing.T) {
	assert.NotNil(t, userSvc.getConnection())
}

func Test_userGetStatus(t *testing.T) {
	resp, err := userSvc.getStatus()
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func Test_userCreateUser(t *testing.T) {
	validEmail := randomdata.Email()

	cases := []struct {
		desc     string
		user     *pblib.User
		isExpErr bool
		errStr   string
	}{
		{
			"Test for valid user registration",
			&pblib.User{
				FirstName:    randomdata.FirstName(randomdata.Male),
				LastName:     randomdata.LastName(),
				Email:        validEmail,
				Password:     "Abcd!123@",
				Organization: "TestOrg",
			},
			false,
			"",
		},
		{
			"Test for duplicate user email",
			&pblib.User{
				FirstName:    randomdata.FirstName(randomdata.Male),
				LastName:     randomdata.LastName(),
				Email:        validEmail,
				Password:     "Abcd!123@",
				Organization: "TestOrg",
			},
			true,
			`rpc error: code = Internal desc = pq: duplicate key value violates unique constraint "accounts_email_key"`,
		},
		{
			"Test for empty string",
			&pblib.User{
				FirstName: "",
			},
			true,
			"rpc error: code = Internal desc = invalid User first name",
		},
		{
			"Test for missing password",
			&pblib.User{
				FirstName:    randomdata.FirstName(randomdata.Male),
				LastName:     randomdata.LastName(),
				Email:        validEmail,
				Password:     "",
				Organization: "TestOrg",
			},
			true,
			"rpc error: code = Internal desc = invalid User password",
		},
	}

	for _, c := range cases {
		resp, err := userSvc.createUser(c.user)
		if c.isExpErr {
			assert.Nil(t, resp, c.desc)
			assert.EqualError(t, err, c.errStr, c.desc)
		} else {
			assert.Nil(t, err, c.desc)
			assert.NotNil(t, resp.GetUser().GetUuid(), c.desc)
			assert.Equal(t, auth.PermissionStringMap[auth.NoPermission], resp.GetUser().GetPermissionLevel(), c.desc)
		}
	}
}

func Test_makeNewAuthSecret(t *testing.T) {
	oldAuthSecret := currAuthSecret
	assert.Nil(t, userSvc.makeNewAuthSecret())
	newAuthSecret, err := userSvc.getAuthSecret()
	assert.Nil(t, err)
	assert.NotEqual(t, oldAuthSecret, newAuthSecret)
}

func Test_getAuthSecret(t *testing.T) {
	assert.Nil(t, userSvc.makeNewAuthSecret())
	authSecret, err := userSvc.getAuthSecret()
	assert.Nil(t, err)
	assert.NotNil(t, auth.ValidateSecret(authSecret))
}

func Test_refreshCurrAuthSecret(t *testing.T) {
	cases := []struct {
		input    *pblib.Secret
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

func Test_replaceCurrAuthSecret(t *testing.T) {
	oldAuthSecret := currAuthSecret
	assert.Nil(t, userSvc.makeNewAuthSecret())
	err := userSvc.replaceCurrAuthSecret()
	assert.Nil(t, err)
	assert.NotEqual(t, oldAuthSecret, currAuthSecret)
}
