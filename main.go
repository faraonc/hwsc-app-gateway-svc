package main

import (
	pbsvc "github.com/hwsc-org/hwsc-api-blocks/int/hwsc-app-gateway-svc/app"
	"github.com/hwsc-org/hwsc-app-gateway-svc/conf"
	"github.com/hwsc-org/hwsc-app-gateway-svc/consts"
	svc "github.com/hwsc-org/hwsc-app-gateway-svc/service"
	log "github.com/hwsc-org/hwsc-lib/logger"
	"google.golang.org/grpc"
	"net"
)

func main() {
	log.Info(consts.AppGatewayServiceTag, "hwsc-app-gateway-svc initiating...")

	// make TCP listener
	lis, err := net.Listen(conf.AppGateWaySvc.Network, conf.AppGateWaySvc.String())
	if err != nil {
		log.Fatal(consts.AppGatewayServiceTag, "Failed to intialize TCP listener:", err.Error())
	}

	// make gRPC server
	s := grpc.NewServer()

	// implement services in /service/service.go
	pbsvc.RegisterAppGatewayServiceServer(s, &svc.Service{})
	log.Info("hwsc-app-gateway-svc at:", conf.AppGateWaySvc.String())

	// start gRPC server
	if err := s.Serve(lis); err != nil {
		log.Fatal(consts.AppGatewayServiceTag, "Failed to serve:", err.Error())
	}
}
