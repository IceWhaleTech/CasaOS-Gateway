package pkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/IceWhaleTech/CasaOS-Common/utils/constants"
	"github.com/IceWhaleTech/CasaOS-Gateway/common"
)

func GetGatewayPort() (int, error) {
	ConfigFilePath := filepath.Join(constants.DefaultConfigPath, common.GatewayName+"."+common.GatewayConfigType)
	if _, err := os.Stat(ConfigFilePath); os.IsNotExist(err) {
		return 0, errors.New(fmt.Sprintf("config file %s not exist", ConfigFilePath))
	}

	config, err := common.LoadConfig()
	if err != nil {
		return 0, err
	}
	if config != nil {
		// convert port to int
		port := config.GetInt(common.ConfigKeyGatewayPort)

		return port, nil
	}
	return 0, errors.New("config is nil")
}
