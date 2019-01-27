package conf

import (
	"fmt"
	log "github.com/hwsc-org/hwsc-logger/logger"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/env"
	"github.com/micro/go-config/source/file"
)

var (
	// AppGateWaySvc the GRPC host for this service
	AppGateWaySvc Host
)

func init() {
	// Create new config
	conf := config.NewConfig()
	if err := conf.Load(file.NewSource(file.WithPath("conf/json/config.dev.json"))); err != nil {
		// TODO - This is a hacky solution for the unit test, because of a weird path issue with GoLang Unit Test
		if err := conf.Load(file.NewSource(file.WithPath("../conf/json/config.dev.json"))); err != nil {
			log.Info("Failed to initialize configuration file", err.Error())
			log.Info("Reading ENV variables")
			src := env.NewSource(
				env.WithPrefix("hosts"),
			)
			if err := conf.Load(src); err != nil {
				log.Fatal("Failed to initialize configuration %v\n", err.Error())
			}
		}
	}

	if err := conf.Get("hosts", "app").Scan(&AppGateWaySvc); err != nil {
		log.Fatal("Failed to get configuration", err.Error())
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
