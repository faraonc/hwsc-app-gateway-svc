package service

import (
	"fmt"
	pb "github.com/faraonc/hwsc-api-blocks/int/hwsc-app-gateway-svc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
)

type Service struct{}

func (s Service) GetStatus(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	fmt.Println(req.GetMessage())
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}
