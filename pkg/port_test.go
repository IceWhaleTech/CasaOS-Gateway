package pkg_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/IceWhaleTech/CasaOS-Common/utils/constants"
	"github.com/IceWhaleTech/CasaOS-Gateway/common"
	"github.com/IceWhaleTech/CasaOS-Gateway/pkg"
	"github.com/stretchr/testify/assert"
)

const _confSample = `[common]
runtimepath=/var/run/casaos

[gateway]
port=80`

func setupGatewayConfig(t *testing.T) {
	// the setup should only run in CICD environment

	ConfigFilePath := filepath.Join(constants.DefaultConfigPath, common.GatewayName+"."+common.GatewayConfigType)
	if _, err := os.Stat(ConfigFilePath); os.IsNotExist(err) {
		// create config file
		file, err := os.Create(ConfigFilePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// write default config
		_, err = file.WriteString(_confSample)
		assert.NoError(t, err)
	}
}

func TestGetPort(t *testing.T) {
	setupGatewayConfig(t)
	port, err := pkg.GetGatewayPort()
	assert.NoError(t, err)
	assert.Equal(t, "80", port)
}
