package service

import (
	pb "github.com/hwsc-org/hwsc-api-blocks/int/hwsc-app-gateway-svc/proto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

func TestGetStatus(t *testing.T) {
	cases := []struct {
		req         *pb.AppGatewayServiceRequest
		serverState state
		expMsg      string
	}{
		{&pb.AppGatewayServiceRequest{}, available, "OK"},
		{&pb.AppGatewayServiceRequest{}, unavailable, "Unavailable"},
	}

	for _, c := range cases {
		serviceStateLocker.currentServiceState = c.serverState
		s := Service{}
		res, _ := s.GetStatus(context.TODO(), c.req)
		assert.Equal(t, c.expMsg, res.GetMessage())
	}
}
