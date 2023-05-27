package test

import (
	client_config "stroxy/client/config"
	"stroxy/logger"
	server_config "stroxy/server/config"
	"testing"
)

func TestMain(m *testing.M) {
	logger.Init()
	m.Run()
}

func TestLoadServerConfigFile(t *testing.T) {
	server_config.Init()
}

func TestLoadClientConfigFile(t *testing.T) {
	client_config.Init()
}
