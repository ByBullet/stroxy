package test

import (
	"testing"

	client_config "github.com/ByBullet/stroxy/client/config"
	"github.com/ByBullet/stroxy/logger"
	server_config "github.com/ByBullet/stroxy/server/config"
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
