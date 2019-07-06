package service

import (
	pbsvc "github.com/hwsc-org/hwsc-api-blocks/protobuf/hwsc-app-gateway-svc/app"
	"github.com/hwsc-org/hwsc-api-blocks/protobuf/lib"
	"github.com/hwsc-org/hwsc-app-gateway-svc/consts"
	log "github.com/hwsc-org/hwsc-lib/logger"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// state of the service
type state uint32

// stateLocker synchronizes the state of the service
type stateLocker struct {
	lock                sync.RWMutex
	currentServiceState state
}

// Service struct type, implements the generated (pb file) AppGatewayServiceServer interface
type Service struct{}

const (
	// TODO must be the same number of services
	numServices = 3

	// available - Service is ready and available
	available state = 0

	// unavailable - Service is unavailable. Example: Provisioning something
	unavailable state = 1
)

var (
	serviceWg          sync.WaitGroup
	serviceStateLocker stateLocker
	serviceStateMap    map[state]string
	currAuthSecret     *lib.Secret
)

func init() {
	// This ensures that all services are disconnected before exit during development mode
	serviceWg.Add(numServices)
	// Handle Terminate Signal(Ctrl + C)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		serviceWg.Wait()
		log.Info(consts.AppGatewayServiceTag, "hwsc-app-gateway-svc terminated")
		os.Exit(0)
	}()

	serviceStateMap = map[state]string{
		available:   "Available",
		unavailable: "Unavailable",
	}
}

func (s state) String() string {
	return serviceStateMap[s]
}

// GetStatus gets the current status of the application gateway
// TODO integration test
func (s *Service) GetStatus(ctx context.Context, req *pbsvc.AppGatewayServiceRequest) (*pbsvc.AppGatewayServiceResponse, error) {
	log.RequestService("GetStatus")

	// Lock the state for reading
	serviceStateLocker.lock.RLock()
	// Unlock the state before function exits
	defer serviceStateLocker.lock.RUnlock()

	log.Info(consts.AppGatewayServiceTag, "Service State:", serviceStateLocker.currentServiceState.String())
	if serviceStateLocker.currentServiceState == unavailable {
		return &pbsvc.AppGatewayServiceResponse{
			Status:  &pbsvc.AppGatewayServiceResponse_Code{Code: uint32(codes.Unavailable)},
			Message: codes.Unavailable.String(),
		}, nil
	}

	return &pbsvc.AppGatewayServiceResponse{
		Status:  &pbsvc.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// CreateUser creates a user
// Returns the user with password field set to empty string
func (s *Service) CreateUser(ctx context.Context, req *pbsvc.AppGatewayServiceRequest) (*pbsvc.AppGatewayServiceResponse, error) {
	log.RequestService("CreateUser")

	if ok := isStateAvailable(); !ok {
		log.Info(consts.AppGatewayServiceTag, consts.ErrServiceUnavailable.Error())
		return nil, status.Error(codes.Unavailable, consts.ErrServiceUnavailable.Error())
	}

	if req == nil || req.GetUserRequest() == nil || req.GetUserRequest().GetUser() == nil {
		log.Error(consts.AppGatewayServiceTag, consts.ErrNilRequest.Error())
		return nil, status.Error(codes.InvalidArgument, consts.ErrNilRequest.Error())
	}
	resp, err := userSvc.createUser(req.GetUserRequest().GetUser())
	if err != nil {
		log.Error(consts.AppGatewayServiceTag, err.Error())
		return nil, err
	}
	return &pbsvc.AppGatewayServiceResponse{
		Status:  &pbsvc.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
		User:    resp.GetUser(),
	}, nil
}

// DeleteUser deletes a user
// Returns the deleted user (TODO decide if we really need to return this to chrome)
func (s *Service) DeleteUser(ctx context.Context, req *pbsvc.AppGatewayServiceRequest) (*pbsvc.AppGatewayServiceResponse, error) {
	// TODO
	return &pbsvc.AppGatewayServiceResponse{
		Status:  &pbsvc.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// UpdateUser updates a user
// Returns the updated user
func (s *Service) UpdateUser(ctx context.Context, req *pbsvc.AppGatewayServiceRequest) (*pbsvc.AppGatewayServiceResponse, error) {
	// TODO
	return &pbsvc.AppGatewayServiceResponse{
		Status:  &pbsvc.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// GetNewAuthToken retrieves a new auth token
// Returns a new auth token string
func (s *Service) GetNewAuthToken(ctx context.Context, req *pbsvc.AppGatewayServiceRequest) (*pbsvc.AppGatewayServiceResponse, error) {
	// TODO
	return &pbsvc.AppGatewayServiceResponse{
		Status:  &pbsvc.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// ListUsers retrieves all the users
// Returns a collection of users
func (s *Service) ListUsers(ctx context.Context, req *pbsvc.AppGatewayServiceRequest) (*pbsvc.AppGatewayServiceResponse, error) {
	// TODO
	return &pbsvc.AppGatewayServiceResponse{
		Status:  &pbsvc.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// GetUser retrieves a user given UUID
// Return the matched user
func (s *Service) GetUser(ctx context.Context, req *pbsvc.AppGatewayServiceRequest) (*pbsvc.AppGatewayServiceResponse, error) {
	// TODO
	return &pbsvc.AppGatewayServiceResponse{
		Status:  &pbsvc.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// ShareDocument shares a user's document to another user
func (s *Service) ShareDocument(ctx context.Context, req *pbsvc.AppGatewayServiceRequest) (*pbsvc.AppGatewayServiceResponse, error) {
	// TODO
	return &pbsvc.AppGatewayServiceResponse{
		Status:  &pbsvc.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// CreateDocument creates a document
// Returns the created document
func (s *Service) CreateDocument(ctx context.Context, req *pbsvc.AppGatewayServiceRequest) (*pbsvc.AppGatewayServiceResponse, error) {
	// TODO
	return &pbsvc.AppGatewayServiceResponse{
		Status:  &pbsvc.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// ListUserDocumentCollection retrieves all documents for a specific user with the given UUID
// Returns all documents for a specific user with the given UUID
func (s *Service) ListUserDocumentCollection(ctx context.Context, req *pbsvc.AppGatewayServiceRequest) (*pbsvc.AppGatewayServiceResponse, error) {
	// TODO
	return &pbsvc.AppGatewayServiceResponse{
		Status:  &pbsvc.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// UpdateDocument updates a document using DUID
// Returns the updated document
func (s *Service) UpdateDocument(ctx context.Context, req *pbsvc.AppGatewayServiceRequest) (*pbsvc.AppGatewayServiceResponse, error) {
	// TODO
	return &pbsvc.AppGatewayServiceResponse{
		Status:  &pbsvc.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// DeleteDocument deletes a document using DUID
// Returns the deleted document
func (s *Service) DeleteDocument(ctx context.Context, req *pbsvc.AppGatewayServiceRequest) (*pbsvc.AppGatewayServiceResponse, error) {
	// TODO
	return &pbsvc.AppGatewayServiceResponse{
		Status:  &pbsvc.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// AddFile upload a new file
// Returns the updated document
func (s *Service) AddFile(ctx context.Context, req *pbsvc.AppGatewayServiceRequest) (*pbsvc.AppGatewayServiceResponse, error) {
	// TODO
	return &pbsvc.AppGatewayServiceResponse{
		Status:  &pbsvc.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// DeleteFile upload a new file
// Returns the updated document
func (s *Service) DeleteFile(ctx context.Context, req *pbsvc.AppGatewayServiceRequest) (*pbsvc.AppGatewayServiceResponse, error) {
	// TODO
	return &pbsvc.AppGatewayServiceResponse{
		Status:  &pbsvc.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// ListDistinctFieldValues retrieves all the unique fields values required for the front-end drop-down filter
// Returns the query transaction
func (s *Service) ListDistinctFieldValues(ctx context.Context, req *pbsvc.AppGatewayServiceRequest) (*pbsvc.AppGatewayServiceResponse, error) {
	// TODO
	return &pbsvc.AppGatewayServiceResponse{
		Status:  &pbsvc.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}

// QueryDocument queries the document service with the given query parameters
// Returns a collection of documents
func (s *Service) QueryDocument(ctx context.Context, req *pbsvc.AppGatewayServiceRequest) (*pbsvc.AppGatewayServiceResponse, error) {
	// TODO
	return &pbsvc.AppGatewayServiceResponse{
		Status:  &pbsvc.AppGatewayServiceResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil
}
