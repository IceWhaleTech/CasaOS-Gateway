package common

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	ConfigKeyLogPath     = "gateway.LogPath"
	ConfigKeyLogSaveName = "gateway.LogSaveName"
	ConfigKeyLogFileExt  = "gateway.LogFileExt"
	ConfigKeyGatewayPort = "gateway.Port"
	ConfigKeyWWWPath     = "gateway.WWWPath"
	ConfigKeyRuntimePath = "common.RuntimePath"
	ConfigKeyNgrokToken  = "ngrok.Token" // nolint: gosec
)

func LoadConfig() (*viper.Viper, error) {
	config := viper.New()

	config.SetDefault(ConfigKeyLogPath, "/var/log/casaos")
	config.SetDefault(ConfigKeyLogSaveName, "gateway")
	config.SetDefault(ConfigKeyLogFileExt, "log")

	config.SetDefault(ConfigKeyWWWPath, "/var/lib/casaos/www")
	config.SetDefault(ConfigKeyRuntimePath, "/var/run/casaos") // See https://refspecs.linuxfoundation.org/FHS_3.0/fhs/ch05s13.html

	config.SetConfigName("gateway")
	config.SetConfigType("ini")

	if currentDirectory, err := os.Getwd(); err != nil {
		log.Println(err)
	} else {
		config.AddConfigPath(currentDirectory)
		config.AddConfigPath(filepath.Join(currentDirectory, "conf"))
	}

	if configPath, success := os.LookupEnv("CASAOS_CONFIG_PATH"); success {
		config.AddConfigPath(configPath)
	}

	config.AddConfigPath(filepath.Join("/", "etc", "casaos"))

	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}

	return config, nil
}
