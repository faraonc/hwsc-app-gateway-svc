package service

import (
	"encoding/base64"
	"github.com/Pallinder/go-randomdata"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	pblib "github.com/hwsc-org/hwsc-api-blocks/protobuf/lib"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestAuthenticate(t *testing.T) {
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
	err := userSvc.verifyEmailToken(emailToken)
	assert.Nil(t, err, "TestAuthenticate create valid user - verify email token")
	validEnc := base64.StdEncoding.EncodeToString([]byte(validEmail + ":" + validPassword))

	header := metadata.New(map[string]string{
		"authorization": "Basic " + validEnc,
	})
	ctx := metadata.NewIncomingContext(context.Background(), header)
	_, expNoErr := Authenticate(ctx)
	assert.Nil(t, expNoErr, "TestAuthenticate create valid user - success")

	_, expNoRepeatErr := Authenticate(ctx)
	assert.Nil(t, expNoRepeatErr, "TestAuthenticate create valid user - repeat success")
	expNoErrNewSecret := userSvc.makeNewAuthSecret()
	assert.Nil(t, expNoErrNewSecret, "TestAuthenticate create valid user - new secret")
	_, expNoRepeatErrWithNewSecret := Authenticate(ctx)
	assert.Nil(t, expNoRepeatErrWithNewSecret, "TestAuthenticate create valid user - repeat success")
	// TODO rename variables
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
