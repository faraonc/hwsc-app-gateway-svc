package service

import (
	"encoding/base64"
	"github.com/Pallinder/go-randomdata"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	pblib "github.com/hwsc-org/hwsc-api-blocks/protobuf/lib"
	"github.com/hwsc-org/hwsc-app-gateway-svc/consts"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	/*
		Test for basic authentication
	*/
	validEmail := randomdata.Email()
	validPassword := "Abcd!123@"
	resp, errCreateUser := userSvc.createUser(
		&pblib.User{
			FirstName:    randomdata.FirstName(randomdata.Male),
			LastName:     randomdata.LastName(),
			Email:        validEmail,
			Password:     validPassword,
			Organization: "TestOrg",
		})
	assert.Nil(t, errCreateUser, "TestAuthenticate create valid user - no err")
	assert.NotNil(t, resp, "TestAuthenticate create valid user - resp not nil")
	emailToken := resp.GetIdentification().GetToken()
	errVerifyEmailToken := userSvc.verifyEmailToken(emailToken)
	assert.Nil(t, errVerifyEmailToken, "TestAuthenticate create valid user - verify email token")
	validEnc := base64.StdEncoding.EncodeToString([]byte(validEmail + ":" + validPassword))

	header := metadata.New(map[string]string{
		"authorization": "Basic " + validEnc,
	})
	validOutgoingCtx := metadata.NewIncomingContext(context.Background(), header)

	incomingCtx1, errAuthenticate := Authenticate(validOutgoingCtx)
	assert.Nil(t, errAuthenticate, "TestAuthenticate create valid user - success")
	incomingMd1, incomingMdOk1 := metadata.FromIncomingContext(incomingCtx1)
	assert.True(t, incomingMdOk1, "TestAuthenticate create valid user - incoming ctx1 OK")
	incomingAuthToken1, incomingAuthTokenOk1 := incomingMd1[consts.StrMdAuthToken]
	assert.True(t, incomingAuthTokenOk1, "TestAuthenticate create valid user - incoming auth token1 OK")
	assert.Equal(t, 1, len(incomingAuthToken1), "TestAuthenticate create valid user - incoming auth token1")

	// repeat Authenticate
	incomingCtx2, errAuthenticate := Authenticate(validOutgoingCtx)
	assert.Nil(t, errAuthenticate, "TestAuthenticate create valid user - repeat success")
	incomingMd2, incomingMdOk2 := metadata.FromIncomingContext(incomingCtx2)
	assert.True(t, incomingMdOk2, "TestAuthenticate create valid user - incoming ctx2 OK")
	incomingAuthToken2, incomingAuthTokenOk2 := incomingMd2[consts.StrMdAuthToken]
	assert.True(t, incomingAuthTokenOk2, "TestAuthenticate create valid user - incoming auth token2 OK")
	assert.Equal(t, 1, len(incomingAuthToken2), "TestAuthenticate create valid user - incoming auth token2")

	// repeat Authenticate with new AuthSecret
	errNewSecret1 := userSvc.makeNewAuthSecret()
	assert.Nil(t, errNewSecret1, "TestAuthenticate create valid user - new secret 1")
	incomingCtx3, errAuthenticate := Authenticate(validOutgoingCtx)
	assert.Nil(t, errAuthenticate, "TestAuthenticate create valid user - repeat success")
	incomingMd3, incomingMdOk3 := metadata.FromIncomingContext(incomingCtx3)
	assert.True(t, incomingMdOk3, "TestAuthenticate create valid user - incoming ctx3 OK")
	incomingAuthToken3, incomingAuthTokenOk3 := incomingMd3[consts.StrMdAuthToken]
	assert.True(t, incomingAuthTokenOk3, "TestAuthenticate create valid user - incoming auth token3 OK")
	assert.Equal(t, 1, len(incomingAuthToken3), "TestAuthenticate create valid user - incoming auth token3")

	/*
		Test for error cases
	*/
	cases := []struct {
		desc       string
		authHeader string
		authPrefix string
		mdKey      string
		input      string
		isExpErr   bool
		errStr     string
	}{
		{
			"test missing auth header",
			"",
			consts.StrBasicAuthPrefix,
			consts.StrMdAuthToken,
			randomdata.Email() + ":" + "Qwert!123@",
			true,
			consts.StatusUnauthenticated.Error(),
		},
		{
			"test wrong auth header",
			"wrong",
			consts.StrBasicAuthPrefix,
			consts.StrMdAuthToken,
			randomdata.Email() + ":" + "Qwert!123@",
			true,
			consts.StatusUnauthenticated.Error(),
		},
		{
			"test missing auth prefix",
			consts.StrMdBasicAuthHeader,
			"",
			consts.StrMdAuthToken,
			randomdata.Email() + ":" + "Qwert!123@",
			true,
			consts.StatusUnauthenticated.Error(),
		},
		{
			"test wrong auth prefix",
			consts.StrMdBasicAuthHeader,
			"wrong",
			consts.StrMdAuthToken,
			randomdata.Email() + ":" + "Qwert!123@",
			true,
			consts.StatusUnauthenticated.Error(),
		},
		{
			"test wrong password",
			consts.StrMdBasicAuthHeader,
			consts.StrBasicAuthPrefix,
			consts.StrMdAuthToken,
			validEmail + ":" + "invalidPassword!123",
			true,
			consts.StatusUnauthenticated.Error(),
		},
		// TODO more test cases
	}
	for _, c := range cases {
		enc := base64.StdEncoding.EncodeToString([]byte(c.input))

		header := metadata.New(map[string]string{
			c.authHeader: c.authPrefix + enc,
		})
		ctx := metadata.NewIncomingContext(context.Background(), header)
		respCtx, err := Authenticate(ctx)
		if c.isExpErr {
			assert.EqualError(t, err, c.errStr, c.desc)
		} else {
			assert.Nil(t, err, c.desc)
			md, ok := metadata.FromIncomingContext(respCtx)
			assert.True(t, ok, c.desc)
			token, ok := md[c.mdKey]
			assert.True(t, ok, c.desc)
			assert.Equal(t, 1, len(token), c.desc)
		}
	}

	// TODO more test cases

}

func TestTryEmailTokenVerification(t *testing.T) {

}

func TestTryTokenAuth(t *testing.T) {

}

func TestTryBasicAuth(t *testing.T) {

}

func TestFinalizeAuth(t *testing.T) {

}

func TestExtractContextHeader(t *testing.T) {

}

func TestPurgeContextHeader(t *testing.T) {

}
