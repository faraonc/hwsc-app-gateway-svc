package main

import (
	pb "github.com/hwsc-org/hwsc-api-blocks/int/hwsc-app-gateway-svc/proto"
	svc "github.com/hwsc-org/hwsc-app-gateway-svc/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	log.Println("[INFO] hwsc-app-gateway-svc initiating...")

	// make TCP listener
	lis, err := net.Listen("tcp", "localhost:50055")
	if err != nil {
		// log.Fatalf will print message to console, then crashes the program
		// %v is the value in a default format
		log.Fatalf("[FATAL] Failed to initialize TCP listener %v\n", err)
	}

	// make gRPC server
	s := grpc.NewServer()

	// implement services in /service/service.go
	pb.RegisterAppGatewayServiceServer(s, &svc.Service{})
	log.Println("[INFO] hwsc-app-gateway-svc at localhost: 50055...")

	// start gRPC server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("[FATAL] Failed to serve %v\n", err)
	}
}
