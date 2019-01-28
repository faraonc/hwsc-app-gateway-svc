package conf

import (
	"fmt"
	"github.com/hwsc-org/hwsc-app-gateway-svc/consts"
	log "github.com/hwsc-org/hwsc-logger/logger"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/env"
	"github.com/micro/go-config/source/file"
)

var (
	// AppGateWaySvc the GRPC host for this service
	AppGateWaySvc Host
	// DocumentSvc the GRPC host for document service
	DocumentSvc Host
	// UserSvc the GRPC host for user service
	UserSvc Host
	// FileTransSvc the GRPC host for file transaction service
	FileTransSvc Host
)

func init() {
	// Create new config
	conf := config.NewConfig()
	if err := conf.Load(file.NewSource(file.WithPath("conf/json/config.dev.json"))); err != nil {
		// TODO - This is a hacky solution for the unit test, because of a weird path issue with GoLang Unit Test
		if err := conf.Load(file.NewSource(file.WithPath("../conf/json/config.dev.json"))); err != nil {
			log.Info(consts.AppGatewayServiceTag, "Failed to initialize configuration file", err.Error())
			log.Info(consts.AppGatewayServiceTag, "Reading ENV variables")
			src := env.NewSource(
				env.WithPrefix("hosts"),
			)
			if err := conf.Load(src); err != nil {
				log.Fatal(consts.AppGatewayServiceTag, "Failed to initialize configuration %v\n", err.Error())
			}
		}
	}

	if err := conf.Get("hosts", "app").Scan(&AppGateWaySvc); err != nil {
		log.Fatal(consts.AppGatewayServiceTag, "Failed to get app-gateway-svc configuration", err.Error())
	}
	if err := conf.Get("hosts", "document").Scan(&DocumentSvc); err != nil {
		log.Fatal(consts.AppGatewayServiceTag, "Failed to get document-svc configuration", err.Error())
	}
	if err := conf.Get("hosts", "user").Scan(&UserSvc); err != nil {
		log.Fatal(consts.AppGatewayServiceTag, "Failed to get user-svc configuration", err.Error())
	}
	if err := conf.Get("hosts", "file").Scan(&FileTransSvc); err != nil {
		log.Fatal(consts.AppGatewayServiceTag, "Failed to get file-transaction-svc configuration", err.Error())
	}
}

// Host represents a server.
type Host struct {
	Address string `json:"address"`
	Port    string `json:"port"`
	Network string `json:"network"`
}

func (h *Host) String() string {
	return fmt.Sprintf("%s:%s", h.Address, h.Port)
}
