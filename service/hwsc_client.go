package service

import (
	log "github.com/hwsc-org/hwsc-logger/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type hwscClient interface {
	dial() error
	getConnection() *grpc.ClientConn
}

func disconnect(conn *grpc.ClientConn, tag string) error {
	if conn == nil {
		log.Info(tag, "Connection already closed")
		return nil
	}
	if err := conn.Close(); err != nil {
		log.Error(tag, "Connection error closed")
		return err
	}
	log.Info(tag, "Connection closed")
	return nil
}

func isHealthy(conn *grpc.ClientConn, tag string) bool {
	if conn == nil {
		return false
	}
	state := conn.GetState()
	log.Info(tag, state.String())
	if state == connectivity.Ready ||
		state == connectivity.Connecting ||
		state == connectivity.Idle {
		return true
	}
	return false
}

func refreshConnection(client hwscClient, tag string) error {
	conn := client.getConnection()
	if conn == nil || !isHealthy(conn, tag) {
		if err := client.dial(); err != nil {
			return err
		}
	}
	return nil
}
