package pkg

import (
	"errors"

	"github.com/IceWhaleTech/CasaOS-Gateway/common"
)

func GetGatewayPort() (int, error) {
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
