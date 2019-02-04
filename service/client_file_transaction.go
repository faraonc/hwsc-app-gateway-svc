package service

import (
	"context"
	"flag"
	pb "github.com/hwsc-org/hwsc-api-blocks/int/hwsc-app-gateway-svc/proto"
	"github.com/hwsc-org/hwsc-app-gateway-svc/conf"
	"github.com/hwsc-org/hwsc-app-gateway-svc/consts"
	log "github.com/hwsc-org/hwsc-lib/logger"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
)

var (
	fileTransSvc     *fileTransactionService
	fileTransSvcAddr = flag.String("file_service_addr", conf.FileTransSvc.String(),
		"The server manages user files")
)

func init() {
	fileTransSvc = &fileTransactionService{}
	if err := refreshConnection(fileTransSvc, consts.FileTransactionClientTag); err != nil {
		// TODO once docker container is runnable
		log.Error(consts.FileTransactionClientTag, err.Error())
		//log.Fatal(consts.FileTransactionClientTag, err.Error())
	}
	// NOTE:
	// app-gateway-svc does not start if all the services are not ready
	// this is ONLY on app-gateway-svc startup
	resp, err := fileTransSvc.getStatus()
	if err != nil {
		// TODO once docker container is runnable
		log.Error(consts.FileTransactionClientTag, err.Error())
		//log.Fatal(consts.FileTransactionClientTag, err.Error())
	} else {
		log.Info(consts.FileTransactionClientTag, resp.String())
	}
	// Handle Terminate Signal(Ctrl + C)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if err := disconnect(fileTransSvc.fileTransSvcConn, consts.FileTransactionClientTag); err != nil {
			log.Error(consts.FileTransactionClientTag, err.Error())
		}
		log.Info(consts.FileTransactionClientTag, "hwsc-app-gateway-svc terminated")
		serviceWg.Done()
	}()

}

type fileTransactionService struct {
	client           pb.FileTransactionServiceClient
	fileTransSvcOpts []grpc.DialOption
	fileTransSvcConn *grpc.ClientConn
}

func (svc *fileTransactionService) dial() error {
	svc.fileTransSvcOpts = nil // set to nil for reconnect purposes
	// TODO
	//if *tls {
	//	if *caFile == "" {
	//		*caFile = testdata.Path("ca.pem")
	//	}
	//	creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
	//	if err != nil {
	//		log.Fatal(consts.FileTransactionClientTag, err.Error)
	//	}
	//	svc.fileTransSvcOpts  = append(svc.fileTransSvcOpts , grpc.WithTransportCredentials(creds))
	//} else {
	//	c.fileTransSvcOpts  = append(svc.fileTransSvcOpts , grpc.WithInsecure())
	//}
	svc.fileTransSvcOpts = append(svc.fileTransSvcOpts, grpc.WithInsecure()) // TODO delete after implementing above TODO
	var err error
	svc.fileTransSvcConn, err = grpc.Dial(*fileTransSvcAddr, svc.fileTransSvcOpts...)
	if err != nil {
		return err
	}
	svc.client = pb.NewFileTransactionServiceClient(svc.fileTransSvcConn)
	return nil
}

func (svc *fileTransactionService) getConnection() *grpc.ClientConn {
	return svc.fileTransSvcConn
}

func (svc *fileTransactionService) getStatus() (*pb.FileTransactionResponse, error) {
	if err := refreshConnection(svc, consts.FileTransactionClientTag); err != nil {
		return nil, err
	}
	// not guaranteed that we are connected, but return the error and try reconnecting again later
	resp, err := svc.client.GetStatus(context.TODO(), &pb.FileTransactionRequest{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
