package service

import (
	"context"
	"flag"
	pbuser "github.com/hwsc-org/hwsc-api-blocks/protobuf/hwsc-user-svc/user"
	"github.com/hwsc-org/hwsc-api-blocks/protobuf/lib"
	"github.com/hwsc-org/hwsc-app-gateway-svc/conf"
	"github.com/hwsc-org/hwsc-app-gateway-svc/consts"
	log "github.com/hwsc-org/hwsc-lib/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"os/signal"
	"syscall"
)

var (
	userSvc     *userService
	userSvcAddr = flag.String("user_service_addr", conf.UserSvc.String(),
		"The server manages user accounts")
)

func init() {
	userSvc = &userService{}
	if err := refreshConnection(userSvc, consts.UserClientTag); err != nil {
		// TODO once docker container is runnable
		log.Error(consts.UserClientTag, err.Error())
		//log.Fatal(consts.UserClientTag, err.Error())
	}
	// NOTE:
	// app-gateway-svc does not start if all the services are not ready
	// this is ONLY on app-gateway-svc startup
	resp, err := userSvc.getStatus()
	if err != nil {
		// TODO once docker container is runnable
		log.Error(consts.UserClientTag, err.Error())
		//log.Fatal(consts.UserClientTag, err.Error())
	} else {
		log.Info(consts.UserClientTag, resp.String())
	}
	// Handle Terminate Signal(Ctrl + C)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if err := disconnect(userSvc.userSvcConn, consts.UserClientTag); err != nil {
			log.Error(consts.UserClientTag, err.Error())
		}
		log.Info(consts.UserClientTag, "hwsc-app-gateway-svc terminated")
		serviceWg.Done()
	}()

}

type userService struct {
	client      pbuser.UserServiceClient
	userSvcOpts []grpc.DialOption
	userSvcConn *grpc.ClientConn
}

func (svc *userService) dial() error {
	svc.userSvcOpts = nil // set to nil for reconnect purposes
	// TODO
	//if *tls {
	//	if *caFile == "" {
	//		*caFile = testdata.Path("ca.pem")
	//	}
	//	creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
	//	if err != nil {
	//		log.Fatal(consts.UserClientTag, err.Error)
	//	}
	//	svc.userSvcOpts  = append(c.userSvcOpts , grpc.WithTransportCredentials(creds))
	//} else {
	//	svc.userSvcOpts  = append(c.userSvcOpts , grpc.WithInsecure())
	//}
	svc.userSvcOpts = append(svc.userSvcOpts, grpc.WithInsecure()) // TODO delete after implementing above TODO
	var err error
	svc.userSvcConn, err = grpc.Dial(*userSvcAddr, svc.userSvcOpts...)
	if err != nil {
		return err
	}
	svc.client = pbuser.NewUserServiceClient(svc.userSvcConn)
	return nil
}

func (svc *userService) getConnection() *grpc.ClientConn {
	return svc.userSvcConn
}

func (svc *userService) getStatus() (*pbuser.UserResponse, error) {
	if err := refreshConnection(svc, consts.UserClientTag); err != nil {
		return nil, err
	}
	// not guaranteed that we are connected, but return the error and try reconnecting again later
	resp, err := svc.client.GetStatus(context.TODO(), &pbuser.UserRequest{})
	st, ok := status.FromError(err)
	if !ok {
		log.Error(consts.UserClientTag, st.Message())
		return nil, status.Error(st.Code(), st.Message())
	}
	return resp, nil
}

func (svc *userService) authenticateUser(email string, password string) (*pbuser.UserResponse, error) {
	if err := refreshConnection(svc, consts.UserClientTag); err != nil {
		return nil, err
	}
	// not guaranteed that we are connected, but return the error and try reconnecting again later
	resp, err := svc.client.AuthenticateUser(
		context.TODO(),
		&pbuser.UserRequest{
			User: &lib.User{
				Email:    email,
				Password: password,
			},
		},
	)
	st, ok := status.FromError(err)
	if !ok {
		log.Error(consts.UserClientTag, st.Message())
		return nil, status.Error(codes.Unauthenticated, st.Message())
	}
	return resp, nil
}

func (svc *userService) verifyAuthToken(token string) (*pbuser.UserResponse, error) {
	if err := refreshConnection(svc, consts.UserClientTag); err != nil {
		return nil, err
	}
	// not guaranteed that we are connected, but return the error and try reconnecting again later
	resp, err := svc.client.VerifyAuthToken(
		context.TODO(),
		&pbuser.UserRequest{
			Identification: &lib.Identification{
				Token: token,
			},
		},
	)
	st, ok := status.FromError(err)
	if !ok {
		log.Error(consts.UserClientTag, st.Message())
		return nil, status.Error(codes.Unauthenticated, st.Message())
	}
	return resp, nil
}

func (svc *userService) verifyEmailToken(token string) (*pbuser.UserResponse, error) {
	if err := refreshConnection(svc, consts.UserClientTag); err != nil {
		return nil, err
	}
	// not guaranteed that we are connected, but return the error and try reconnecting again later
	resp, err := svc.client.VerifyEmailToken(
		context.TODO(),
		&pbuser.UserRequest{
			Identification: &lib.Identification{
				Token: token,
			},
		},
	)
	st, ok := status.FromError(err)
	if !ok {
		log.Error(consts.UserClientTag, st.Message())
		return nil, status.Error(st.Code(), st.Message())
	}
	return resp, nil
}
