package service

import (
	"github.com/Pallinder/go-randomdata"
	pblib "github.com/hwsc-org/hwsc-api-blocks/protobuf/lib"
	"github.com/hwsc-org/hwsc-lib/auth"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClientUserSvcDial(t *testing.T) {
	err := userSvc.userSvcConn.Close()
	assert.Nil(t, err, "test closing connection")
	err = userSvc.dial()
	assert.Nil(t, err, "test dialing with after closing connection")
}

func TestClientUserGetConnection(t *testing.T) {
	assert.NotNil(t, userSvc.getConnection())
}

func TestClientUserGetStatus(t *testing.T) {
	resp, err := userSvc.getStatus()
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestClientUserCreateUser(t *testing.T) {
	validEmail := randomdata.Email()
	cases := []struct {
		desc     string
		user     *pblib.User
		isExpErr bool
		errStr   string
	}{
		{
			"test for valid user registration",
			&pblib.User{
				FirstName:    randomdata.FirstName(randomdata.Male),
				LastName:     randomdata.LastName(),
				Email:        validEmail,
				Password:     "Abcd!123@",
				Organization: testOrg,
			},
			false,
			"",
		},
		{
			"test for duplicate user email",
			&pblib.User{
				FirstName:    randomdata.FirstName(randomdata.Male),
				LastName:     randomdata.LastName(),
				Email:        validEmail,
				Password:     "Abcd!123@",
				Organization: testOrg,
			},
			true,
			`rpc error: code = Internal desc = pq: duplicate key value violates unique constraint "accounts_email_key"`,
		},
		{
			"test for input empty string",
			&pblib.User{
				FirstName: "",
			},
			true,
			"rpc error: code = Internal desc = invalid User first name",
		},
		{
			"test for missing password",
			&pblib.User{
				FirstName:    randomdata.FirstName(randomdata.Male),
				LastName:     randomdata.LastName(),
				Email:        validEmail,
				Password:     "",
				Organization: testOrg,
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

func TestClientMakeNewAuthSecret(t *testing.T) {
	oldAuthSecret := currAuthSecret
	assert.Nil(t, userSvc.makeNewAuthSecret())
	newAuthSecret, err := userSvc.getAuthSecret()
	assert.Nil(t, err)
	assert.NotEqual(t, oldAuthSecret, newAuthSecret)
}

func TestClientGetAuthSecret(t *testing.T) {
	assert.Nil(t, userSvc.makeNewAuthSecret())
	newSecret, err := userSvc.getAuthSecret()
	assert.Nil(t, err)
	assert.Nil(t, auth.ValidateSecret(newSecret))
	assert.Equal(t, currAuthSecret, newSecret, "test auth secret is set")
}

func TestClientAuthenticateUser(t *testing.T) {
	validEmail := randomdata.Email()
	validPassword := "Abcd!123@"
	resp, err := userSvc.createUser(
		&pblib.User{
			FirstName:    randomdata.FirstName(randomdata.Male),
			LastName:     randomdata.LastName(),
			Email:        validEmail,
			Password:     validPassword,
			Organization: testOrg,
		})
	assert.Nil(t, err, "verify no err for create user")
	assert.NotNil(t, resp, "verify resp is not nil for create use")
	emailToken := resp.GetIdentification().GetToken()
	err = userSvc.verifyEmailToken(emailToken)
	assert.Nil(t, err, "verify the email")
	cases := []struct {
		desc     string
		email    string
		password string
		isExpErr bool
		errStr   string
	}{
		{
			"test for non-existing user",
			randomdata.Email(),
			"ASD1231!",
			true,
			"rpc error: code = Unauthenticated desc = email does not exist in db",
		},
		{
			"test for missing email",
			"",
			"ASD1231!",
			true,
			"rpc error: code = InvalidArgument desc = invalid User email",
		},
		{
			"test for missing password",
			validEmail,
			"",
			true,
			"rpc error: code = InvalidArgument desc = invalid User password",
		},
		{
			"test for valid email password",
			validEmail,
			validPassword,
			false,
			"",
		},
	}
	for _, c := range cases {
		resp, err := userSvc.authenticateUser(c.email, c.password)
		if c.isExpErr {
			assert.EqualError(t, err, c.errStr, c.desc)
		} else {
			assert.Nil(t, err, c.desc)
			assert.NotNil(t, resp, c.desc)
			assert.Equal(t, c.email, resp.GetUser().GetEmail(), c.desc)
			assert.Empty(t, resp.GetUser().GetPassword(), c.desc)
			assert.Equal(t, currAuthSecret, resp.GetIdentification().GetSecret(), c.desc)
		}
	}
}

func TestClientVerifyAuthToken(t *testing.T) {
	validEmail := randomdata.Email()
	validPassword := "Abcd!123@"
	resp, err := userSvc.createUser(
		&pblib.User{
			FirstName:    randomdata.FirstName(randomdata.Male),
			LastName:     randomdata.LastName(),
			Email:        validEmail,
			Password:     validPassword,
			Organization: testOrg,
		})
	assert.Nil(t, err, "Test_verifyAuthToken")
	assert.NotNil(t, resp, "Test_verifyAuthToken")
	emailToken := resp.GetIdentification().GetToken()
	err = userSvc.verifyEmailToken(emailToken)
	assert.Nil(t, err, "verify the email")

	resp, err = userSvc.authenticateUser(validEmail, validPassword)
	assert.Nil(t, err, "Test_verifyAuthToken")
	assert.NotNil(t, resp, "Test_verifyAuthToken")
	authToken := resp.GetIdentification().GetToken()
	authSecret := resp.GetIdentification().GetSecret()

	cases := []struct {
		desc     string
		token    string
		isExpErr bool
		errStr   string
	}{
		{
			"test for valid auth token",
			authToken,
			false,
			"",
		},
		{
			"test for invalid token type",
			emailToken,
			true,
			"rpc error: code = Unauthenticated desc = no matching auth token were found with given token",
		},
		{
			"test for empty string",
			"",
			true,
			"rpc error: code = Unauthenticated desc = empty token string",
		},
		{
			"test for fake token",
			fakeAuthToken,
			true,
			"rpc error: code = Unauthenticated desc = no matching auth token were found with given token",
		},
		{
			// cannot really test unless we can somehow insert an auth token in the db
			// may require user-svc changes
			"test for expired token",
			expiredUserToken,
			true,
			"rpc error: code = Unauthenticated desc = no matching auth token were found with given token",
		},
	}
	for _, c := range cases {
		resp, err := userSvc.verifyAuthToken(c.token)
		if c.isExpErr {
			assert.EqualError(t, err, c.errStr)
			assert.Nil(t, resp)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, c.token, resp.GetIdentification().GetToken())
			assert.Equal(t, authSecret, resp.GetIdentification().GetSecret())
		}
	}
}

func TestClientVerifyEmailToken(t *testing.T) {
	validCase := "test valid email token"
	validEmail := randomdata.Email()
	validPassword := "Abcd!123@"
	resp, err := userSvc.createUser(
		&pblib.User{
			FirstName:    randomdata.FirstName(randomdata.Male),
			LastName:     randomdata.LastName(),
			Email:        validEmail,
			Password:     validPassword,
			Organization: testOrg,
		})
	assert.Nil(t, err, validCase)
	assert.NotNil(t, resp, validCase)

	emailToken := resp.GetIdentification().GetToken()
	expNilErr := userSvc.verifyEmailToken(emailToken)
	assert.Nil(t, expNilErr, validCase)

	expErr := userSvc.verifyEmailToken(emailToken)
	assert.EqualError(t, expErr,
		"rpc error: code = Internal desc = no matching email token were found with given token",
		"test verify email token that was already verified")

	err = userSvc.verifyEmailToken("")
	assert.EqualError(t, err,
		"rpc error: code = InvalidArgument desc = empty token string",
		"test for empty string")

}

func TestClientRefreshCurrAuthSecret(t *testing.T) {
	cases := []struct {
		desc     string
		input    *pblib.Secret
		isExpErr bool
		err      error
	}{
		{
			"test setting nil current auth secret",
			nil,
			false,
			nil,
		},
		{
			"test setting expired current auth secret",
			expiredAuthSecret,
			false,
			nil,
		},
		{
			"test not replacing a valid current auth secret",
			validAuthSecret,
			false,
			nil,
		},
	}
	for _, c := range cases {
		currAuthSecret = c.input
		err := userSvc.refreshCurrAuthSecret()
		if c.isExpErr {
			assert.EqualError(t, err, c.err.Error(), c.desc)
		} else {
			assert.Nil(t, err, c.desc)
		}
	}
}

func TestClientReplaceCurrAuthSecret(t *testing.T) {
	oldAuthSecret := currAuthSecret
	assert.Nil(t, userSvc.makeNewAuthSecret())
	err := userSvc.replaceCurrAuthSecret()
	assert.Nil(t, err)
	assert.NotEqual(t, oldAuthSecret, currAuthSecret)
}
