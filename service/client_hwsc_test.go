package service

import (
	"github.com/hwsc-org/hwsc-app-gateway-svc/consts"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDisconnect(t *testing.T) {

}

func TestIsHealthy(t *testing.T) {

}

func TestRefreshConnection(t *testing.T) {
	cases := []struct {
		client    hwscClient
		expErr    bool
		expErrMsg string
	}{
		{
			nil,
			true,
			consts.ErrNilHwscGrpcClient.Error(),
		},
		{
			userSvc,
			false,
			"",
		},
		{
			documentSvc,
			false,
			"",
		},
		{
			fileTransSvc,
			false,
			"",
		},
	}

	for _, c := range cases {
		err := refreshConnection(c.client, "placeholder")
		if c.expErr {
			assert.EqualError(t, err, c.expErrMsg)
		} else {
			assert.Nil(t, err)
		}
	}

	assert.Nil(t, userSvc.userSvcConn.Close(), "test close userSvc client connection")
	err := refreshConnection(userSvc, "placeholder")
	assert.Nil(t, err, "test setting userSvc to nil - err nil")
	assert.NotNil(t, userSvc, "test setting userSvc to nil - userSvc not nil")
}
