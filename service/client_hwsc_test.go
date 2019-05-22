package service

import (
	"github.com/hwsc-org/hwsc-app-gateway-svc/consts"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
)

func TestDisconnect(t *testing.T) {
	assert.Nil(t, disconnect(nil, placeholder), "test nil input")
	assert.Nil(t, disconnect(userSvc.getConnection(), placeholder), "test disconnect user svc client")
	assert.NotNil(t, disconnect(userSvc.getConnection(), placeholder), "test disconnect again")
	assert.Nil(t, refreshConnection(userSvc, placeholder), "reconnect again no error")

	cases := []struct {
		desc      string
		client    hwscClient
		expErr    bool
		expErrMsg string
	}{
		{
			"test user svc client connection",
			userSvc,
			false,
			"",
		},
		{
			"test document svc client connection",
			documentSvc,
			false,
			"",
		},
		{
			"test file trnsaction svc client client connection",
			fileTransSvc,
			false,
			"",
		},
	}

	for _, c := range cases {
		err := disconnect(c.client.getConnection(), placeholder)
		if c.expErr {
			assert.EqualError(t, err, c.expErrMsg, c.desc)
		} else {
			assert.Nil(t, err, c.desc)
			assert.Nil(t, refreshConnection(c.client, placeholder), c.desc)

		}
	}
}

func TestIsHealthy(t *testing.T) {
	cases := []struct {
		desc      string
		client    *grpc.ClientConn
		expOutput bool
	}{
		{
			"test nil client connection",
			nil,
			false,
		},
		{
			"test user svc client connection",
			userSvc.getConnection(),
			true,
		},
		// TODO https://github.com/hwsc-org/hwsc-app-gateway-svc/issues/54
		//{
		//	"test document svc client connection",
		//	documentSvc.getConnection(),
		//	true,
		//},
		//{
		//	"test file transaction svc client connection",
		//	fileTransSvc.getConnection(),
		//	true,
		//},
	}

	for _, c := range cases {
		actOutput := isHealthy(c.client, placeholder)
		assert.Equal(t, c.expOutput, actOutput, c.desc)
	}

	assert.Nil(t, userSvc.userSvcConn.Close(), "test close userSvc client connection")
	output := isHealthy(userSvc.getConnection(), placeholder)
	assert.Equal(t, false, output, "test setting userSvc to nil - err nil")
	err := refreshConnection(userSvc, placeholder)
	assert.Nil(t, err, "test close userSvc client connection - err nil")
}

func TestRefreshConnection(t *testing.T) {
	cases := []struct {
		desc      string
		client    hwscClient
		expErr    bool
		expErrMsg string
	}{
		{
			"test nil client connection",
			nil,
			true,
			consts.ErrNilHwscGrpcClient.Error(),
		},
		{
			"test user svc client connection",
			userSvc,
			false,
			"",
		},
		{
			"test document svc client connection",
			documentSvc,
			false,
			"",
		},
		{
			"test file trnsaction svc client client connection",
			fileTransSvc,
			false,
			"",
		},
	}

	for _, c := range cases {
		err := refreshConnection(c.client, placeholder)
		if c.expErr {
			assert.EqualError(t, err, c.expErrMsg, c.desc)
		} else {
			assert.Nil(t, err, c.desc)
		}
	}

	assert.Nil(t, userSvc.userSvcConn.Close(), "test close userSvc client connection")
	err := refreshConnection(userSvc, placeholder)
	assert.Nil(t, err, "test close userSvc client connection - err nil")
	assert.NotNil(t, userSvc, "test close userSvc client connection - userSvc not nil")
}
