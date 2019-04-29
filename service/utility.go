package service

import (
	"github.com/hwsc-org/hwsc-app-gateway-svc/consts"
	"github.com/hwsc-org/hwsc-lib/logger"
)

func isStateAvailable() bool {
	// Lock the state for reading
	serviceStateLocker.lock.RLock()
	// Unlock the state before function exits
	defer serviceStateLocker.lock.RUnlock()

	logger.Info(consts.AppGatewayServiceTag, serviceStateLocker.currentServiceState.String())
	if serviceStateLocker.currentServiceState != available {
		return false
	}

	return true
}
