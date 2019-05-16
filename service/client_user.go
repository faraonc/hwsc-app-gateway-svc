package service

import (
	"context"
	"flag"
	pbuser "github.com/hwsc-org/hwsc-api-blocks/protobuf/hwsc-user-svc/user"
	"github.com/hwsc-org/hwsc-api-blocks/protobuf/lib"
	"github.com/hwsc-org/hwsc-app-gateway-svc/conf"
	"github.com/hwsc-org/hwsc-app-gateway-svc/consts"
	"github.com/hwsc-org/hwsc-lib/auth"
	log "github.com/hwsc-org/hwsc-lib/logger"
	"google.golang.org/grpc"
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

	if err := userSvc.refreshCurrAuthSecret(); err != nil {
		log.Error(consts.UserClientTag, err.Error())
	}
	log.Info(consts.UserClientTag, "AuthSecret obtained")

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

// dial to user-svc.
// Returns an error if it exists.
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

// getConnection returns the grpc client connection.
func (svc *userService) getConnection() *grpc.ClientConn {
	return svc.userSvcConn
}

// getStatus enables the client to check the user-svc.
func (svc *userService) getStatus() (*pbuser.UserResponse, error) {
	if err := refreshConnection(svc, consts.UserClientTag); err != nil {
		return nil, err
	}
	// not guaranteed that we are connected, but return the error and try reconnecting again later
	resp, err := svc.client.GetStatus(context.TODO(), &pbuser.UserRequest{})
	if err != nil {
		log.Error(consts.UserClientTag, err.Error())
		return nil, err
	}
	return resp, nil
}

// createUser creates a user
// On success, returns user object with password set to empty for security reasons.
func (svc *userService) createUser(user *lib.User) (*pbuser.UserResponse, error) {
	if err := refreshConnection(svc, consts.UserClientTag); err != nil {
		return nil, err
	}
	// not guaranteed that we are connected, but return the error and try reconnecting again later
	resp, err := svc.client.CreateUser(context.TODO(), &pbuser.UserRequest{User: user})
	if err != nil {
		log.Error(consts.UserClientTag, err.Error())
		return nil, err
	}
	return resp, nil
}

// makeNewAuthSecret forces user-svc to replace its current AuthSecret.
// Returns an error from the user-svc.
func (svc *userService) makeNewAuthSecret() error {
	if err := refreshConnection(svc, consts.UserClientTag); err != nil {
		return err
	}
	// not guaranteed that we are connected, but return the error and try reconnecting again later
	_, err := svc.client.MakeNewAuthSecret(context.TODO(), &pbuser.UserRequest{})
	if err != nil {
		log.Error(consts.UserClientTag, err.Error())
		return err
	}
	return nil
}

// getAuthSecret gets the current AuthSecret in user-svc.
// Returns an error from the user-svc.
func (svc *userService) getAuthSecret() (*lib.Secret, error) {
	if err := refreshConnection(svc, consts.UserClientTag); err != nil {
		return nil, err
	}
	// not guaranteed that we are connected, but return the error and try reconnecting again later
	resp, err := svc.client.GetAuthSecret(context.TODO(), &pbuser.UserRequest{})
	if err != nil {
		log.Error(consts.UserClientTag, err.Error())
		return nil, err
	}
	newSecret := resp.GetIdentification().GetSecret()
	if err := auth.ValidateSecret(newSecret); err != nil {
		log.Error(consts.UserClientTag, err.Error())
		return nil, err
	}
	currAuthSecret = newSecret
	return resp.Identification.Secret, nil
}

// authenticateUser authenticates user with email and password.
// On success, returns a JWT and User.
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
	if err != nil {
		log.Error(consts.UserClientTag, err.Error())
		return nil, err
	}
	newSecret := resp.GetIdentification().GetSecret()
	if err := auth.ValidateSecret(newSecret); err != nil {
		log.Error(consts.UserClientTag, err.Error())
		return nil, err
	}
	currAuthSecret = newSecret
	return resp, nil
}

// verifyAuthToken checks if auth token or JWT is valid.
// On success, returns JWT and paired secret.
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
	if err != nil {
		log.Error(consts.UserClientTag, err.Error())
		return nil, err
	}
	newSecret := resp.GetIdentification().GetSecret()
	if err := auth.ValidateSecret(newSecret); err != nil {
		log.Error(consts.UserClientTag, err.Error())
		return nil, err
	}
	currAuthSecret = newSecret
	return resp, nil
}

// verifyEmailToken verifies a prospective email.
// Returns an error if email token verification fails.
func (svc *userService) verifyEmailToken(token string) error {
	if err := refreshConnection(svc, consts.UserClientTag); err != nil {
		return err
	}
	// not guaranteed that we are connected, but return the error and try reconnecting again later
	_, err := svc.client.VerifyEmailToken(
		context.TODO(),
		&pbuser.UserRequest{
			Identification: &lib.Identification{
				Token: token,
			},
		},
	)
	if err != nil {
		log.Error(consts.UserClientTag, err.Error())
		return err
	}
	return nil
}

// refreshCurrAuthSecret refreshes currAuthSecret if it is invalid.
// Returns an error from the user-svc.
func (svc *userService) refreshCurrAuthSecret() error {
	if err := auth.ValidateSecret(currAuthSecret); err != nil {
		err = nil
		currAuthSecret, err = userSvc.getAuthSecret()
		if err != nil {
			return consts.ErrUnableToUpdateAuthSecret
		}
	}
	return nil
}

// replaceCurrAuthSecret force to replace currAuthSecret even if still valid.
// Returns an error from the user-svc.
func (svc *userService) replaceCurrAuthSecret() error {
	newAuthSecret, err := userSvc.getAuthSecret()
	if err != nil {
		return consts.ErrUnableToUpdateAuthSecret
	}
	currAuthSecret = newAuthSecret
	return nil
}
