package pkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/IceWhaleTech/CasaOS-Common/utils/constants"
	"github.com/IceWhaleTech/CasaOS-Gateway/common"
)

func GetGatewayPort() (string, error) {
	ConfigFilePath := filepath.Join(constants.DefaultConfigPath, common.GatewayName+"."+common.GatewayConfigType)
	if _, err := os.Stat(ConfigFilePath); os.IsNotExist(err) {
		return "", errors.New(fmt.Sprintf("config file %s not exist", ConfigFilePath))
	}

	config, err := common.LoadConfig()
	if err != nil {
		return "", err
	}
	if config != nil {
		return config.GetString(common.ConfigKeyGatewayPort), nil
	}
	return "", errors.New("config is nil")
}
