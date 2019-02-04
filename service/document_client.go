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
	"sync"
	"syscall"
)

var (
	documentSvc     *documentService
	documentSvcAddr = flag.String("document_service_addr", conf.DocumentSvc.String(),
		"The server manages user documents")
)

func init() {
	documentSvc = &documentService{}
	if err := refreshConnection(documentSvc, consts.DocumentClientTag); err != nil {
		log.Fatal(consts.DocumentClientTag, err.Error())
	}
	// NOTE:
	// app-gateway-svc does not start if all the services are not ready
	// this is ONLY on app-gateway-svc startup
	resp, err := documentSvc.getStatus()
	if err != nil {
		log.Fatal(consts.DocumentClientTag, err.Error())
	} else {
		log.Info(consts.DocumentClientTag, resp.String())
	}
	// Handle Terminate Signal(Ctrl + C)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if err := disconnect(documentSvc.documentSvcConn, consts.DocumentClientTag); err != nil {
			log.Error(consts.DocumentClientTag, err.Error())
		}
		log.Error(consts.DocumentClientTag, "hwsc-app-gateway-svc terminated")
		serviceWg.Done()
	}()

}

type documentService struct {
	client          pb.DocumentServiceClient
	lock            sync.RWMutex //TODO implement locks
	documentSvcOpts []grpc.DialOption
	documentSvcConn *grpc.ClientConn
}

func (svc *documentService) dial() error {
	svc.documentSvcOpts = nil // set to nil for reconnect purposes
	// TODO
	//if *tls {
	//	if *caFile == "" {
	//		*caFile = testdata.Path("ca.pem")
	//	}
	//	creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
	//	if err != nil {
	//		log.Fatal(documentClientTag, err.Error)
	//	}
	//	svc.documentSvcOpts  = append(svc.documentSvcOpts , grpc.WithTransportCredentials(creds))
	//} else {
	//	svc.documentSvcOpts  = append(svc.documentSvcOpts , grpc.WithInsecure())
	//}
	svc.documentSvcOpts = append(svc.documentSvcOpts, grpc.WithInsecure()) // TODO delete after implementing above TODO
	var err error
	svc.documentSvcConn, err = grpc.Dial(*documentSvcAddr, svc.documentSvcOpts...)
	if err != nil {
		return err
	}
	svc.client = pb.NewDocumentServiceClient(svc.documentSvcConn)
	return nil
}

func (svc *documentService) getConnection() *grpc.ClientConn {
	return svc.documentSvcConn
}

func (svc *documentService) getStatus() (*pb.DocumentResponse, error) {
	if err := refreshConnection(svc, consts.DocumentClientTag); err != nil {
		return nil, err
	}
	// not guaranteed that we are connected, but return the error and try reconnecting again later
	resp, err := svc.client.GetStatus(context.TODO(), &pb.DocumentRequest{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
