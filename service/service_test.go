package service

import (
	"github.com/Pallinder/go-randomdata"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	pbsvc "github.com/hwsc-org/hwsc-api-blocks/protobuf/hwsc-app-gateway-svc/app"
	pbuser "github.com/hwsc-org/hwsc-api-blocks/protobuf/hwsc-user-svc/user"
	pblib "github.com/hwsc-org/hwsc-api-blocks/protobuf/lib"
	"github.com/hwsc-org/hwsc-lib/auth"
	"github.com/hwsc-org/hwsc-lib/hosts"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
	"time"
)

var (
	// UserDB contains user database configs grabbed from env vars
	UserDB hosts.UserDBHost
)

func TestGetStatus(t *testing.T) {
	cases := []struct {
		req         *pbsvc.AppGatewayServiceRequest
		serverState state
		expMsg      string
	}{
		{&pbsvc.AppGatewayServiceRequest{}, available, "OK"},
		{&pbsvc.AppGatewayServiceRequest{}, unavailable, "Unavailable"},
	}

	for _, c := range cases {
		serviceStateLocker.currentServiceState = c.serverState
		s := Service{}
		res, _ := s.GetStatus(context.TODO(), c.req)
		assert.Equal(t, c.expMsg, res.GetMessage())
	}
}

func TestCreateUser(t *testing.T) {
	validEmail := randomdata.Email()

	cases := []struct {
		desc        string
		user        *pblib.User
		isExpErr    bool
		errStr      string
		serverState state
	}{
		{
			"Test for unavailable service",
			&pblib.User{
				FirstName:    randomdata.FirstName(randomdata.Male),
				LastName:     randomdata.LastName(),
				Email:        validEmail,
				Password:     "Abcd!123@",
				Organization: testOrg,
			},
			true,
			"rpc error: code = Unavailable desc = service is unavailable",
			unavailable,
		},
		{
			"Test for nil user",
			nil,
			true,
			"rpc error: code = InvalidArgument desc = nil request",
			available,
		},
		{
			"Test for valid user registration",
			&pblib.User{
				FirstName:    randomdata.FirstName(randomdata.Male),
				LastName:     randomdata.LastName(),
				Email:        validEmail,
				Password:     "Abcd!123@",
				Organization: testOrg,
			},
			false,
			"",
			available,
		},
		{
			"Test for duplicate user email",
			&pblib.User{
				FirstName:    randomdata.FirstName(randomdata.Male),
				LastName:     randomdata.LastName(),
				Email:        validEmail,
				Password:     "Abcd!123@",
				Organization: testOrg,
			},
			true,
			`rpc error: code = Internal desc = pq: duplicate key value violates unique constraint "accounts_email_key"`,
			available,
		},
		{
			"Test for empty string",
			&pblib.User{
				FirstName: "",
			},
			true,
			"rpc error: code = Internal desc = invalid User first name",
			available,
		},
		{
			"Test for missing password",
			&pblib.User{
				FirstName:    randomdata.FirstName(randomdata.Male),
				LastName:     randomdata.LastName(),
				Email:        validEmail,
				Password:     "",
				Organization: testOrg,
			},
			true,
			"rpc error: code = Internal desc = invalid User password",
			available,
		},
	}

	for _, c := range cases {
		serviceStateLocker.currentServiceState = c.serverState
		s := Service{}
		resp, err := s.CreateUser(context.TODO(), &pbsvc.AppGatewayServiceRequest{
			UserRequest: &pbuser.UserRequest{
				User: c.user,
			},
		})
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

func TestService_GetNewAuthToken(t *testing.T) {
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
	assert.Nil(t, err, "no error creating user")
	assert.NotNil(t, resp, "response is not nil for creating user")
	emailToken := resp.GetIdentification().GetToken()
	err = userSvc.verifyEmailToken(emailToken)
	assert.Nil(t, err, "no error verifying the email")

	resp, err = userSvc.authenticateUser(validEmail, validPassword)
	assert.Nil(t, err, "no error authenticating user")
	assert.NotNil(t, resp, "response is not nil for authenticating user")
	authToken := resp.GetIdentification().GetToken()
	// sleep to ensure a different auth token is generated
	time.Sleep(2 * time.Second)

	cases := []struct {
		desc        string
		authToken   string
		isExpErr    bool
		errStr      string
		serverState state
	}{
		{
			"Test for unavailable service",
			"",
			true,
			"rpc error: code = Unavailable desc = service is unavailable",
			unavailable,
		},
		{
			"Test for empty string",
			"",
			true,
			"rpc error: code = InvalidArgument desc = nil request",
			available,
		},
		{
			"Test for invalid token type",
			fakeAuthToken,
			true,
			"rpc error: code = DeadlineExceeded desc = no matching auth token were found with given token",
			available,
		},
		{
			"Test for expired token",
			expiredUserToken,
			true,
			"rpc error: code = DeadlineExceeded desc = no matching auth token were found with given token",
			available,
		},
		{
			"Test for valid token",
			authToken,
			false,
			"",
			available,
		},
	}

	for _, c := range cases {
		serviceStateLocker.currentServiceState = c.serverState
		s := Service{}
		resp, err := s.GetNewAuthToken(context.TODO(), &pbsvc.AppGatewayServiceRequest{
			UserRequest: &pbuser.UserRequest{
				Identification: &pblib.Identification{
					Token: c.authToken,
				},
			},
		})
		if c.isExpErr {
			assert.Nil(t, resp, c.desc)
			assert.EqualError(t, err, c.errStr, c.desc)
		} else {
			assert.Nil(t, err, c.desc)
			assert.NotNil(t, resp, c.desc)
			assert.NotZero(t, resp.GetToken(), c.desc)
			assert.NotEqual(t, c.authToken, resp.GetToken(), c.desc)
		}
	}
}
