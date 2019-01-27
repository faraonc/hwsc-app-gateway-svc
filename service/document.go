package service

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/hwsc-org/hwsc-api-blocks/int/hwsc-app-gateway-svc/proto"
	log "github.com/hwsc-org/hwsc-logger/logger"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
)

var (
	documentSvcClient pb.DocumentServiceClient

	documentSvcAddr = flag.String("document_service_addr", ":50051",
		"The server manages user documents")

	documentSvcOpts    []grpc.DialOption
	documentSvcConn    *grpc.ClientConn
	documentSvcConnErr error
)

func init() {
	if err := refreshDocumentServiceClient(); err != nil {
		log.Fatal(documentClientTag, err.Error())
	}
	// Handle Terminate Signal(Ctrl + C)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		_ = documentSvcConn.Close()
		fmt.Println()
		log.Fatal(documentClientTag, "hwsc-app-gateway-svc terminated")
	}()
	resp, _ := documentSvcClient.GetStatus(context.TODO(), &pb.DocumentRequest{})
	log.Info(documentClientTag, resp.String())
}

func refreshDocumentServiceClient() error {
	// TODO
	//if *tls {
	//	if *caFile == "" {
	//		*caFile = testdata.Path("ca.pem")
	//	}
	//	creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
	//	if err != nil {
	//		log.Fatalf("Failed to create TLS credentials %v", err)
	//	}
	//	opts = append(opts, grpc.WithTransportCredentials(creds))
	//} else {
	//	opts = append(opts, grpc.WithInsecure())
	//}
	documentSvcOpts = append(documentSvcOpts, grpc.WithInsecure()) // TODO delete after implementing above TODO

	documentSvcConnErr = nil
	documentSvcConn, documentSvcConnErr = grpc.Dial(*documentSvcAddr, documentSvcOpts...)
	if documentSvcConnErr != nil {
		return documentSvcConnErr
	}

	documentSvcClient = pb.NewDocumentServiceClient(documentSvcConn)
	return nil
}
