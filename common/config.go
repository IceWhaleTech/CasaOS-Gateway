package common

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	ConfigKeyGatewayPort = "gateway.Port"
	ConfigKeyWWWPath     = "gateway.WWWPath"
	ConfigKeyRuntimePath = "common.RuntimePath"

	DefaultGatewayPort = "80"
)

func LoadConfig() (*viper.Viper, error) {
	config := viper.New()

	config.SetDefault(ConfigKeyGatewayPort, DefaultGatewayPort)
	config.SetDefault(ConfigKeyWWWPath, "/var/lib/casaos/www")
	config.SetDefault(ConfigKeyRuntimePath, "/var/run/casaos") // See https://refspecs.linuxfoundation.org/FHS_3.0/fhs/ch05s13.html

	config.SetConfigName("gateway")
	config.SetConfigType("ini")

	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	if configPath, success := os.LookupEnv("CASAOS_CONFIG_PATH"); success {
		config.AddConfigPath(configPath)
	}

	config.AddConfigPath(currentDirectory)
	config.AddConfigPath(filepath.Join(currentDirectory, "conf"))
	config.AddConfigPath(filepath.Join("/", "etc", "casaos"))

	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}

	return config, nil
}
