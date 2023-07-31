package common

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/IceWhaleTech/CasaOS-Common/utils/constants"
)

const (
	ConfigKeyLogPath     = "gateway.LogPath"
	ConfigKeyLogSaveName = "gateway.LogSaveName"
	ConfigKeyLogFileExt  = "gateway.LogFileExt"
	ConfigKeyGatewayPort = "gateway.Port"
	ConfigKeyWWWPath     = "gateway.WWWPath"
	ConfigKeyRuntimePath = "common.RuntimePath"

	GatewayConfigName = "gateway"
	GatewayConfigType = "ini"
)

func LoadConfig() (*viper.Viper, error) {
	config := viper.New()

	config.SetDefault(ConfigKeyLogPath, constants.DefaultLogPath)
	config.SetDefault(ConfigKeyLogSaveName, "gateway")
	config.SetDefault(ConfigKeyLogFileExt, "log")

	config.SetDefault(ConfigKeyWWWPath, filepath.Join(constants.DefaultDataPath, "www"))
	config.SetDefault(ConfigKeyRuntimePath, constants.DefaultRuntimePath) // See https://refspecs.linuxfoundation.org/FHS_3.0/fhs/ch05s13.html

	config.SetConfigName(GatewayConfigName)
	config.SetConfigType(GatewayConfigType)

	if currentDirectory, err := os.Getwd(); err != nil {
		log.Println(err)
	} else {
		config.AddConfigPath(currentDirectory)
		config.AddConfigPath(filepath.Join(currentDirectory, "conf"))
	}

	if configPath, success := os.LookupEnv("CASAOS_CONFIG_PATH"); success {
		config.AddConfigPath(configPath)
	}

	config.AddConfigPath(constants.DefaultConfigPath)

	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}

	return config, nil
}
