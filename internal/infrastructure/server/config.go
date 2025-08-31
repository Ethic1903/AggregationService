package server

import (
	"fmt"
	"os"
)

const (
	_serverHost = "SERVER_HOST"
	_serverPort = "SERVER_PORT"
)

type IServerConfig interface {
	Address() string
}

type serverConfig struct {
	host string
	port string
}

func NewServerConfig() (IServerConfig, error) {
	host := os.Getenv(_serverHost)
	if len(host) == 0 {
		return nil, fmt.Errorf("env %s is empty", _serverHost)
	}

	port := os.Getenv(_serverPort)
	if len(port) == 0 {
		return nil, fmt.Errorf("env %s is empty", _serverPort)
	}

	return &serverConfig{
		host: host,
		port: port,
	}, nil
}

func (s serverConfig) Address() string {
	return s.host + ":" + s.port
}
