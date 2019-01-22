package service

import (
	pb "github.com/hwsc-org/hwsc-api-blocks/int/hwsc-app-gateway-svc/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
)

type Service struct{}

func (s *Service) GetStatus(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}
func (s *Service) CreateUser(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}
func (s *Service) DeleteUser(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}
func (s *Service) UpdateUser(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}
func (s *Service) AuthenticateUser(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}
func (s *Service) ListUsers(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}
func (s *Service) GetUser(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}
func (s *Service) ShareDocument(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}
func (s *Service) CreateDocument(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}
func (s *Service) ListUserDocumentCollection(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}
func (s *Service) UpdateDocument(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}
func (s *Service) DeleteDocument(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}
func (s *Service) AddFile(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}
func (s *Service) DeleteFile(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}
func (s *Service) ListDistinctFieldValues(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}
func (s *Service) QueryDocument(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}
