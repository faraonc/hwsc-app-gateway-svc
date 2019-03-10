package service

import (
	pbsvc "github.com/hwsc-org/hwsc-api-blocks/int/hwsc-app-gateway-svc/app"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

func TestMain(m *testing.M) {

}


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
