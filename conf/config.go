package conf

import (
	"github.com/hwsc-org/hwsc-app-gateway-svc/consts"
	"github.com/hwsc-org/hwsc-lib/hosts"
	"github.com/hwsc-org/hwsc-lib/logger"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/env"
)

var (
	// AppGateWaySvc the GRPC host for this service
	AppGateWaySvc hosts.Host
	// DocumentSvc the GRPC host for document service
	DocumentSvc hosts.Host
	// UserSvc the GRPC host for user service
	UserSvc hosts.Host
	// FileTransSvc the GRPC host for file transaction service
	FileTransSvc hosts.Host
)

func init() {
	// Create new config
	conf := config.NewConfig()

	logger.Info(consts.AppGatewayServiceTag, "Reading ENV variables")
	src := env.NewSource(
		env.WithPrefix("hosts"),
	)
	if err := conf.Load(src); err != nil {
		logger.Fatal(consts.AppGatewayServiceTag, "Failed to initialize configuration %v\n", err.Error())
	}
	if err := conf.Get("hosts", "app").Scan(&AppGateWaySvc); err != nil {
		logger.Fatal(consts.AppGatewayServiceTag, "Failed to get app-gateway-svc configuration", err.Error())
	}
	if err := conf.Get("hosts", "document").Scan(&DocumentSvc); err != nil {
		logger.Fatal(consts.AppGatewayServiceTag, "Failed to get document-svc configuration", err.Error())
	}
	if err := conf.Get("hosts", "user").Scan(&UserSvc); err != nil {
		logger.Fatal(consts.AppGatewayServiceTag, "Failed to get user-svc configuration", err.Error())
	}
	if err := conf.Get("hosts", "file").Scan(&FileTransSvc); err != nil {
		logger.Fatal(consts.AppGatewayServiceTag, "Failed to get file-transaction-svc configuration", err.Error())
	}
}
