package service

import (
	pb "github.com/hwsc-org/hwsc-api-blocks/int/hwsc-app-gateway-svc/proto"
	"github.com/hwsc-org/hwsc-app-gateway-svc/consts"
	log "github.com/hwsc-org/hwsc-logger/logger"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// TODO state when can we provide services

const (
	numServices = 3
)

var (
	serviceWg sync.WaitGroup
)

func init() {
	serviceWg.Add(numServices)
	// Handle Terminate Signal(Ctrl + C)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		serviceWg.Wait()
		log.Fatal(consts.AppGatewayServiceTag, "hwsc-app-gateway-svc terminated")
	}()

}

// Service struct type, implements the generated (pb file) AppGatewayServiceServer interface
type Service struct{}

// GetStatus gets the current status of the application gateway
func (s *Service) GetStatus(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// CreateUser creates a user
// Returns the user with password field set to empty string
func (s *Service) CreateUser(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// DeleteUser deletes a user
// Returns the deleted user (TODO decide if we really need to return this to chrome)
func (s *Service) DeleteUser(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// UpdateUser updates a user
// Returns the updated user
func (s *Service) UpdateUser(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// AuthenticateUser looks through users and perform email and password match
// Returns matched user
func (s *Service) AuthenticateUser(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// ListUsers retrieves all the users
// Returns a collection of users
func (s *Service) ListUsers(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// GetUser retrieves a user given UUID
// Return the matched user
func (s *Service) GetUser(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// ShareDocument shares a user's document to another user
func (s *Service) ShareDocument(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// CreateDocument creates a document
// Returns the created document
func (s *Service) CreateDocument(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// ListUserDocumentCollection retrieves all documents for a specific user with the given UUID
// Returns all documents for a specific user with the given UUID
func (s *Service) ListUserDocumentCollection(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// UpdateDocument updates a document using DUID
// Returns the updated document
func (s *Service) UpdateDocument(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// DeleteDocument deletes a document using DUID
// Returns the deleted document
func (s *Service) DeleteDocument(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// AddFile upload a new file
// Returns the updated document
func (s *Service) AddFile(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// DeleteFile upload a new file
// Returns the updated document
func (s *Service) DeleteFile(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// ListDistinctFieldValues retrieves all the unique fields values required for the front-end drop-down filter
// Returns the query transaction
func (s *Service) ListDistinctFieldValues(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// QueryDocument queries the document service with the given query parameters
// Returns a collection of documents
func (s *Service) QueryDocument(ctx context.Context, req *pb.AppGatewayServiceRequest) (*pb.AppGatewayServiceResponse, error) {
	// TODO
	return &pb.AppGatewayServiceResponse{
		Status:  &pb.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}
