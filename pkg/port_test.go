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

const _incorrectConfSample = `[common]
runtimepath=/var/run/casaos

[gateway]
port=`

func setupGatewayConfig(t *testing.T) func() {
	// the setup should only run in CICD environment

	testInCICD := false

	ConfigFilePath := filepath.Join(constants.DefaultConfigPath, common.GatewayName+"."+common.GatewayConfigType)
	if _, err := os.Stat(ConfigFilePath); os.IsNotExist(err) {
		// create config file
		os.MkdirAll(constants.DefaultConfigPath, os.ModePerm)
		file, err := os.Create(ConfigFilePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// write default config
		_, err = file.WriteString(_confSample)
		assert.NoError(t, err)
		testInCICD = true
	}

	return func() {
		if testInCICD {
			// remove config file
			err := os.Remove(ConfigFilePath)
			assert.NoError(t, err)
		}
	}
}

func setupIncorrectGatewayConfig(t *testing.T) {
	// the setup should only run in CICD environment

	ConfigFilePath := filepath.Join(constants.DefaultConfigPath, common.GatewayName+"."+common.GatewayConfigType)
	if _, err := os.Stat(ConfigFilePath); os.IsNotExist(err) {
		// create config file
		os.MkdirAll(constants.DefaultConfigPath, os.ModePerm)
		file, err := os.Create(ConfigFilePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// write default config
		_, err = file.WriteString(_incorrectConfSample)
		assert.NoError(t, err)
	}
}

func TestGetPort(t *testing.T) {
	defer setupGatewayConfig(t)()
	port, err := pkg.GetGatewayPort()
	assert.NoError(t, err)
	assert.Equal(t, 80, port)
}

func TestGetBlankPort(t *testing.T) {
	ConfigFilePath := filepath.Join(constants.DefaultConfigPath, common.GatewayName+"."+common.GatewayConfigType)
	// only run in CICD environment
	if _, err := os.Stat(ConfigFilePath); !os.IsNotExist(err) {
		t.Skip("the test only run in CICD environment to avoid overwrite the config file")
	}

	setupIncorrectGatewayConfig(t)
	port, err := pkg.GetGatewayPort()
	assert.NotNil(t, err)
	assert.Equal(t, 0, port)
}
