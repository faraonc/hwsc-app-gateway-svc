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
	baseAuthCase := "test authenticating valid user"
	validEmail := randomdata.Email()
	validPassword := "Abcd!123@"
	resp, errCreateUser := userSvc.createUser(
		&pblib.User{
			FirstName:    randomdata.FirstName(randomdata.Male),
			LastName:     randomdata.LastName(),
			Email:        validEmail,
			Password:     validPassword,
			Organization: testOrg,
		})
	assert.Nil(t, errCreateUser, baseAuthCase)
	assert.NotNil(t, resp, baseAuthCase)
	emailToken := resp.GetIdentification().GetToken()
	errVerifyEmailToken := userSvc.verifyEmailToken(emailToken)
	assert.Nil(t, errVerifyEmailToken, baseAuthCase)
	validEnc := base64.StdEncoding.EncodeToString([]byte(validEmail + ":" + validPassword))

	header := metadata.New(map[string]string{
		consts.StrMdBasicAuthHeader: consts.StrBasicAuthPrefix + validEnc,
	})
	validIncomingCtx := metadata.NewIncomingContext(context.Background(), header)

	incomingCtx1, errAuthenticate := Authenticate(validIncomingCtx)
	assert.Nil(t, errAuthenticate, baseAuthCase)
	incomingMd1, incomingMdOk1 := metadata.FromIncomingContext(incomingCtx1)
	assert.True(t, incomingMdOk1, baseAuthCase)
	incomingAuthToken1, incomingAuthTokenOk1 := incomingMd1[consts.StrMdAuthToken]
	assert.True(t, incomingAuthTokenOk1, baseAuthCase)
	assert.Equal(t, 1, len(incomingAuthToken1), baseAuthCase)

	// repeat Authenticate
	baseAuthRepeatCase := "test with repeated authentication"
	incomingCtx2, errAuthenticate := Authenticate(validIncomingCtx)
	assert.Nil(t, errAuthenticate, baseAuthRepeatCase)
	incomingMd2, incomingMdOk2 := metadata.FromIncomingContext(incomingCtx2)
	assert.True(t, incomingMdOk2, baseAuthRepeatCase)
	incomingAuthToken2, incomingAuthTokenOk2 := incomingMd2[consts.StrMdAuthToken]
	assert.True(t, incomingAuthTokenOk2, baseAuthRepeatCase)
	assert.Equal(t, 1, len(incomingAuthToken2), baseAuthRepeatCase)

	// repeat Authenticate with new AuthSecret
	baseAuthNewAuthSecretCase := "test authenticate with new auth secret"
	errNewSecret1 := userSvc.makeNewAuthSecret()
	assert.Nil(t, errNewSecret1, baseAuthNewAuthSecretCase)
	incomingCtx3, errAuthenticate := Authenticate(validIncomingCtx)
	assert.Nil(t, errAuthenticate, baseAuthNewAuthSecretCase)
	incomingMd3, incomingMdOk3 := metadata.FromIncomingContext(incomingCtx3)
	assert.True(t, incomingMdOk3, baseAuthNewAuthSecretCase)
	incomingAuthToken3, incomingAuthTokenOk3 := incomingMd3[consts.StrMdAuthToken]
	assert.True(t, incomingAuthTokenOk3, baseAuthNewAuthSecretCase)
	assert.Equal(t, 1, len(incomingAuthToken3), baseAuthNewAuthSecretCase)

	/*
		Test for auth token
	*/
	authHeader := metadata.New(map[string]string{
		consts.StrMdBasicAuthHeader: consts.StrTokenAuthPrefix + incomingAuthToken3[0],
	})
	authCtx := metadata.NewIncomingContext(context.Background(), authHeader)
	authIncomingCtx, errAuth := Authenticate(authCtx)
	assert.Nil(t, errAuth, "TestAuthenticate auth token - no err")
	authMd, authOk := metadata.FromIncomingContext(authIncomingCtx)
	assert.True(t, authOk, "TestAuthenticate auth token - incoming ctx OK")
	authToken, authTokenOk := authMd[consts.StrMdAuthToken]
	assert.True(t, authTokenOk, "TestAuthenticate create auth token - incoming auth token OK")
	assert.Equal(t, 1, len(authToken), "TestAuthenticate auth token - incoming auth token")
	assert.Equal(t, authToken[0], incomingAuthToken3[0], "TestAuthenticate auth token - same token")

	/*
		Test for email token verification
	*/
	tokenAuthCase := "test authenticating with valid auth token"
	newEmailVerificationUserEmail := randomdata.Email()
	newEmailVerificationUserPassword := "Abcd!123@"
	newEmailVerificationResp, errNewEmailVerificationCreateUser := userSvc.createUser(
		&pblib.User{
			FirstName:    randomdata.FirstName(randomdata.Male),
			LastName:     randomdata.LastName(),
			Email:        newEmailVerificationUserEmail,
			Password:     newEmailVerificationUserPassword,
			Organization: testOrg,
		})
	assert.Nil(t, errNewEmailVerificationCreateUser, tokenAuthCase)
	assert.NotNil(t, newEmailVerificationResp, tokenAuthCase)
	newEmailVerificationToken := newEmailVerificationResp.GetIdentification().GetToken()
	newEmailVerificationHeader := metadata.New(map[string]string{
		consts.StrMdBasicAuthHeader: consts.StrEmailTokenVerificationPrefix + newEmailVerificationToken,
	})
	newEmailVerificationOutgoingCtx := metadata.NewIncomingContext(context.Background(), newEmailVerificationHeader)
	_, errNewEmailVerification := Authenticate(newEmailVerificationOutgoingCtx)
	assert.Nil(t, errNewEmailVerification, tokenAuthCase)

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
		{
			"test fake auth token",
			consts.StrMdBasicAuthHeader,
			consts.StrTokenAuthPrefix,
			consts.StrMdAuthToken,
			fakeAuthToken,
			true,
			consts.StatusUnauthenticated.Error(),
		},
		{
			"test fake email token",
			consts.StrMdBasicAuthHeader,
			consts.StrEmailTokenVerificationPrefix,
			consts.StrTokenAuthPrefix,
			fakeAuthToken,
			true,
			consts.StatusUnauthenticated.Error(),
		},
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
			assert.NotNil(t, respCtx, c.desc)
			md, ok := metadata.FromIncomingContext(respCtx)
			assert.True(t, ok, c.desc)
			token, ok := md[c.mdKey]
			assert.True(t, ok, c.desc)
			assert.Equal(t, 1, len(token), c.desc)
		}
	}

	nilIncomingCtxCase := "test for nil incoming context"
	ctxNil, errNil := Authenticate(nil)
	assert.EqualError(t, errNil, consts.ErrNilContext.Error(), nilIncomingCtxCase)
	assert.Equal(t, context.TODO(), ctxNil, nilIncomingCtxCase)
}

func TestTryEmailTokenVerification(t *testing.T) {
	nilIncomingCtxCase := "test for nil incoming context"
	ctxNil, errNil := tryEmailTokenVerification(nil)
	assert.EqualError(t, errNil, consts.ErrNilContext.Error(), nilIncomingCtxCase)
	assert.Equal(t, context.TODO(), ctxNil, nilIncomingCtxCase)

	validCase := "test verify email token"
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

	cases := []struct {
		desc       string
		authHeader string
		authPrefix string
		input      string
		isExpErr   bool
		errStr     string
	}{
		{
			"test missing auth header",
			"",
			consts.StrEmailTokenVerificationPrefix,
			"",
			true,
			"rpc error: code = Unauthenticated desc = no \"authorization\" header in request",
		},
		{
			"test wrong auth header",
			"wrong",
			consts.StrEmailTokenVerificationPrefix,
			"",
			true,
			"rpc error: code = Unauthenticated desc = no \"authorization\" header in request",
		},
		{
			"test missing auth prefix",
			consts.StrMdBasicAuthHeader,
			"",
			"",
			true,
			"rpc error: code = Unauthenticated desc = missing \"Email Token \" prefix in \"authorization\" header",
		},
		{
			"test wrong auth prefix",
			consts.StrMdBasicAuthHeader,
			"wrong",
			"",
			true,
			"rpc error: code = Unauthenticated desc = missing \"Email Token \" prefix in \"authorization\" header",
		},
		{
			"test fake email token",
			consts.StrMdBasicAuthHeader,
			consts.StrEmailTokenVerificationPrefix,
			fakeAuthToken,
			true,
			"rpc error: code = Internal desc = no matching email token were found with given token",
		},
		{
			validCase,
			consts.StrMdBasicAuthHeader,
			consts.StrEmailTokenVerificationPrefix,
			emailToken,
			false,
			"",
		},
	}
	for _, c := range cases {
		header := metadata.New(map[string]string{
			c.authHeader: c.authPrefix + c.input,
		})
		ctx := metadata.NewIncomingContext(context.Background(), header)
		respCtx, err := tryEmailTokenVerification(ctx)
		if c.isExpErr {
			assert.EqualError(t, err, c.errStr, c.desc)
		} else {
			assert.Nil(t, err, c.desc)
			assert.NotNil(t, respCtx, c.desc)
		}
	}

}

func TestTryTokenAuth(t *testing.T) {
	nilIncomingCtxCase := "test for nil incoming context"
	ctxNil, errNil := tryTokenAuth(nil)
	assert.EqualError(t, errNil, consts.ErrNilContext.Error(), nilIncomingCtxCase)
	assert.Equal(t, context.TODO(), ctxNil, nilIncomingCtxCase)

	validCase := "test valid auth token"
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
	errVerifyEmailToken := userSvc.verifyEmailToken(emailToken)
	assert.Nil(t, errVerifyEmailToken, validCase)
	validEnc := base64.StdEncoding.EncodeToString([]byte(validEmail + ":" + validPassword))
	header := metadata.New(map[string]string{
		consts.StrMdBasicAuthHeader: consts.StrBasicAuthPrefix + validEnc,
	})
	validIncomingCtx := metadata.NewIncomingContext(context.Background(), header)
	incomingCtx, errAuthenticate := Authenticate(validIncomingCtx)
	assert.Nil(t, errAuthenticate, validCase)
	incomingMd, incomingMdOk := metadata.FromIncomingContext(incomingCtx)
	assert.True(t, incomingMdOk, validCase)
	incomingAuthToken, incomingAuthTokenOk := incomingMd[consts.StrMdAuthToken]
	assert.True(t, incomingAuthTokenOk, validCase)
	assert.Equal(t, 1, len(incomingAuthToken), validCase)
	authHeader := metadata.New(map[string]string{
		consts.StrMdBasicAuthHeader: consts.StrTokenAuthPrefix + incomingAuthToken[0],
	})
	authCtx := metadata.NewIncomingContext(context.Background(), authHeader)
	authIncomingCtx, errAuth := Authenticate(authCtx)
	assert.Nil(t, errAuth, validCase)
	authMd, authOk := metadata.FromIncomingContext(authIncomingCtx)
	assert.True(t, authOk, validCase)
	authToken, authTokenOk := authMd[consts.StrMdAuthToken]
	assert.True(t, authTokenOk, validCase)
	assert.Equal(t, 1, len(authToken), validCase)
	assert.Equal(t, authToken[0], incomingAuthToken[0], validCase)

	cases := []struct {
		desc       string
		authHeader string
		authPrefix string
		input      string
		isExpErr   bool
		errStr     string
	}{
		{
			"test missing auth header",
			"",
			consts.StrTokenAuthPrefix,
			"",
			true,
			"rpc error: code = Unauthenticated desc = no \"authorization\" header in request",
		},
		{
			"test wrong auth header",
			"wrong",
			consts.StrTokenAuthPrefix,
			"",
			true,
			"rpc error: code = Unauthenticated desc = no \"authorization\" header in request",
		},
		{
			"test missing auth prefix",
			consts.StrMdBasicAuthHeader,
			"",
			"",
			true,
			"rpc error: code = Unauthenticated desc = missing \"Auth Token \" prefix in \"authorization\" header",
		},
		{
			"test wrong auth prefix",
			consts.StrMdBasicAuthHeader,
			"wrong",
			"",
			true,
			"rpc error: code = Unauthenticated desc = missing \"Auth Token \" prefix in \"authorization\" header",
		},
		{
			"test fake auth token",
			consts.StrMdBasicAuthHeader,
			consts.StrTokenAuthPrefix,
			fakeAuthToken,
			true,
			"rpc error: code = Unauthenticated desc = no matching auth token were found with given token",
		},
	}
	for _, c := range cases {
		header := metadata.New(map[string]string{
			c.authHeader: c.authPrefix + c.input,
		})
		ctx := metadata.NewIncomingContext(context.Background(), header)
		respCtx, err := tryTokenAuth(ctx)
		assert.EqualError(t, err, c.errStr, c.desc)
		assert.Equal(t, ctx, respCtx, c.desc)
	}

}

func TestTryBasicAuth(t *testing.T) {
	nilIncomingCtxCase := "test for nil incoming context"
	ctxNil, errNil := tryBasicAuth(nil)
	assert.EqualError(t, errNil, consts.ErrNilContext.Error(), nilIncomingCtxCase)
	assert.Equal(t, context.TODO(), ctxNil, nilIncomingCtxCase)

	validCase := "test valid basic auth"
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
	errVerifyEmailToken := userSvc.verifyEmailToken(emailToken)
	assert.Nil(t, errVerifyEmailToken, validCase)

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
			"rpc error: code = Unauthenticated desc = no \"authorization\" header in request",
		},
		{
			"test wrong auth header",
			"wrong",
			consts.StrBasicAuthPrefix,
			consts.StrMdAuthToken,
			randomdata.Email() + ":" + "Qwert!123@",
			true,
			"rpc error: code = Unauthenticated desc = no \"authorization\" header in request",
		},
		{
			"test missing auth prefix",
			consts.StrMdBasicAuthHeader,
			"",
			consts.StrMdAuthToken,
			randomdata.Email() + ":" + "Qwert!123@",
			true,
			"rpc error: code = Unauthenticated desc = missing \"Basic \" prefix in \"authorization\" header",
		},
		{
			"test wrong auth prefix",
			consts.StrMdBasicAuthHeader,
			"wrong",
			consts.StrMdAuthToken,
			randomdata.Email() + ":" + "Qwert!123@",
			true,
			"rpc error: code = Unauthenticated desc = missing \"Basic \" prefix in \"authorization\" header",
		},
		{
			"test with missing input",
			consts.StrMdBasicAuthHeader,
			consts.StrBasicAuthPrefix,
			consts.StrMdAuthToken,
			"",
			true,
			"rpc error: code = Unauthenticated desc = invalid basic auth format",
		},
		{
			"test with missing colon",
			consts.StrMdBasicAuthHeader,
			consts.StrBasicAuthPrefix,
			consts.StrMdAuthToken,
			validEmail,
			true,
			"rpc error: code = Unauthenticated desc = invalid basic auth format",
		},
		{
			"test with missing password",
			consts.StrMdBasicAuthHeader,
			consts.StrBasicAuthPrefix,
			consts.StrMdAuthToken,
			validEmail + ":",
			true,
			"rpc error: code = Unauthenticated desc = invalid basic auth format",
		},
		{
			"test wrong password",
			consts.StrMdBasicAuthHeader,
			consts.StrBasicAuthPrefix,
			consts.StrMdAuthToken,
			validEmail + ":" + "invalidPassword!123",
			true,
			"rpc error: code = Unauthenticated desc = crypto/bcrypt: hashedPassword is not the hash of the given password",
		},
		{
			validCase,
			consts.StrMdBasicAuthHeader,
			consts.StrBasicAuthPrefix,
			consts.StrMdAuthToken,
			validEmail + ":" + validPassword,
			false,
			"",
		},
	}
	for _, c := range cases {
		enc := base64.StdEncoding.EncodeToString([]byte(c.input))

		header := metadata.New(map[string]string{
			c.authHeader: c.authPrefix + enc,
		})
		ctx := metadata.NewIncomingContext(context.Background(), header)
		respCtx, err := tryBasicAuth(ctx)
		if c.isExpErr {
			assert.EqualError(t, err, c.errStr, c.desc)
		} else {
			assert.Nil(t, err, c.desc)
			assert.NotNil(t, respCtx, c.desc)
			md, ok := metadata.FromIncomingContext(respCtx)
			assert.True(t, ok, c.desc)
			token, ok := md[c.mdKey]
			assert.True(t, ok, c.desc)
			assert.Equal(t, 1, len(token), c.desc)
		}
	}
	badEncodingCase := "test bad base64 encoding"
	header := metadata.New(map[string]string{
		consts.StrMdBasicAuthHeader: consts.StrBasicAuthPrefix + placeholder,
	})
	ctx := metadata.NewIncomingContext(context.Background(), header)
	respCtx, err := tryBasicAuth(ctx)
	assert.Equal(t, ctx, respCtx, badEncodingCase)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = invalid base64 in header", badEncodingCase)

}

func TestFinalizeAuth(t *testing.T) {
	nilIncomingCtxCase := "test for nil incoming context"
	actOutput, err := finalizeAuth(nil, placeholder, nil)
	assert.EqualError(t, err, consts.ErrNilContext.Error(), nilIncomingCtxCase)
	assert.Equal(t, context.TODO(), actOutput, nilIncomingCtxCase)

	missingTagCase := "test for missing tag"
	actOutput, err = finalizeAuth(context.TODO(), "", nil)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = missing tag", missingTagCase)
	assert.Equal(t, context.TODO(), actOutput, missingTagCase)

	invalidIdentificationCase := "test for invalid identification"
	actOutput, err = finalizeAuth(context.TODO(), consts.BasicAuthTag, nil)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = nil identification", invalidIdentificationCase)
	assert.Equal(t, context.TODO(), actOutput, invalidIdentificationCase)

	validCase := "test for valid ctx"
	header := metadata.New(map[string]string{
		consts.StrMdBasicAuthHeader: consts.StrBasicAuthPrefix + placeholder,
	})
	ctx := metadata.NewIncomingContext(context.Background(), header)
	actOutput, err = finalizeAuth(ctx, consts.BasicAuthTag, validIdentification)
	assert.Nil(t, err, "", validCase)
	md, ok := metadata.FromIncomingContext(actOutput)
	assert.True(t, ok, validCase)
	authHeader, ok := md[consts.StrMdAuthToken]
	assert.True(t, ok, validCase)
	assert.Equal(t, validIdentification.GetToken(), authHeader[0], validCase)
}

func TestExtractContextHeader(t *testing.T) {
	nilIncomingCtxCase := "test for nil incoming context"
	actOutput, errNil := extractContextHeader(nil, "")
	assert.EqualError(t, errNil, consts.ErrNilContext.Error(), nilIncomingCtxCase)
	assert.Zero(t, actOutput, nilIncomingCtxCase)

	// test for wrong context by using OutgoingCtx
	invalidOutgoingCtxCase := "test for invalid context using OutgoingCtx instead of IncomingCtx"
	validEmail := randomdata.Email()
	validPassword := "Abcd!123@"
	resp, errCreateUser := userSvc.createUser(
		&pblib.User{
			FirstName:    randomdata.FirstName(randomdata.Male),
			LastName:     randomdata.LastName(),
			Email:        validEmail,
			Password:     validPassword,
			Organization: testOrg,
		})
	assert.Nil(t, errCreateUser, invalidOutgoingCtxCase)
	assert.NotNil(t, resp, invalidOutgoingCtxCase)
	emailToken := resp.GetIdentification().GetToken()
	errVerifyEmailToken := userSvc.verifyEmailToken(emailToken)
	assert.Nil(t, errVerifyEmailToken, invalidOutgoingCtxCase)
	validEnc := base64.StdEncoding.EncodeToString([]byte(validEmail + ":" + validPassword))

	header := metadata.New(map[string]string{
		consts.StrMdBasicAuthHeader: consts.StrBasicAuthPrefix + validEnc,
	})
	invalidOutgoingCtx := metadata.NewOutgoingContext(context.Background(), header)
	actualOutput1, err := extractContextHeader(invalidOutgoingCtx, consts.StrMdBasicAuthHeader)
	assert.Zero(t, actualOutput1, invalidOutgoingCtxCase)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = no headers in request", invalidOutgoingCtxCase)

	invalidHeaderCase := "test for authorization invalid header"
	validIncomingCtx1 := metadata.NewIncomingContext(context.Background(), header)
	actualOutput2, err := extractContextHeader(validIncomingCtx1, placeholder)
	assert.Zero(t, actualOutput2, invalidHeaderCase)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = no \"authorization\" header in request",
		invalidHeaderCase)

	multiHeaderCase := "test for multi authorization headers"
	md1 := metadata.New(map[string]string{
		consts.StrMdBasicAuthHeader: consts.StrBasicAuthPrefix + validEnc,
	})
	md2 := metadata.New(map[string]string{
		consts.StrMdBasicAuthHeader: consts.StrBasicAuthPrefix + validEnc,
	})
	multiHeader := metadata.Join(md1, md2)
	invalidIncomingCtx2 := metadata.NewIncomingContext(context.Background(), multiHeader)
	actualOutput3, err := extractContextHeader(invalidIncomingCtx2, consts.StrMdBasicAuthHeader)
	assert.Zero(t, actualOutput3, multiHeaderCase)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = more than 1 header in request", multiHeaderCase)

	cases := []struct {
		desc       string
		authHeader string
		authPrefix string
		input      string
		isExpErr   bool
		errStr     string
	}{
		{
			"test missing auth header",
			"",
			consts.StrBasicAuthPrefix,
			randomdata.Email() + ":" + "Qwert!123@",
			true,
			"rpc error: code = Unauthenticated desc = missing header",
		},
		{
			"test for valid context",
			consts.StrMdBasicAuthHeader,
			consts.StrBasicAuthPrefix,
			randomdata.Email() + ":" + "Qwert!123@",
			false,
			"",
		},
	}
	for _, c := range cases {
		enc := base64.StdEncoding.EncodeToString([]byte(c.input))

		header := metadata.New(map[string]string{
			c.authHeader: c.authPrefix + enc,
		})
		ctx := metadata.NewIncomingContext(context.Background(), header)
		actOutput, err := extractContextHeader(ctx, c.authHeader)
		if c.isExpErr {
			assert.EqualError(t, err, c.errStr, c.desc)
		} else {
			assert.Nil(t, err, c.desc)
			assert.Equal(t, c.authPrefix+enc, actOutput, c.desc)
		}
	}
}

func TestPurgeContextHeader(t *testing.T) {
	nilIncomingCtxCase := "test for nil incoming context"
	actOutput, err := purgeContextHeader(nil, "")
	assert.EqualError(t, err, consts.ErrNilContext.Error(), nilIncomingCtxCase)
	assert.Nil(t, actOutput, nilIncomingCtxCase)

	missingHeaderCase := "test for missing header"
	actOutput, err = purgeContextHeader(context.TODO(), "")
	assert.EqualError(t, err, consts.ErrMissingHeader.Error(), missingHeaderCase)
	assert.Nil(t, actOutput, missingHeaderCase)

	validCase := "test for valid context"
	header := metadata.New(map[string]string{
		consts.StrMdBasicAuthHeader: consts.StrBasicAuthPrefix + placeholder,
	})
	ctx := metadata.NewIncomingContext(context.Background(), header)
	newCtx, err := purgeContextHeader(ctx, consts.StrMdBasicAuthHeader)
	assert.Nil(t, err, validCase)
	md, ok := metadata.FromIncomingContext(newCtx)
	assert.True(t, ok, validCase)
	authHeader, ok := md[consts.StrMdBasicAuthHeader]
	assert.True(t, ok, validCase)
	assert.Nil(t, authHeader, validCase)
}
