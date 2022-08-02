package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/spf13/viper"
)

const (
	configKeyGatewayPort = "gateway.Port"
	configKeyRuntimePath = "common.RuntimePath"
)

func Load(state *service.State) error {
	viper.SetDefault(configKeyGatewayPort, "80")
	viper.SetDefault(configKeyRuntimePath, "/var/run/casaos") // See https://refspecs.linuxfoundation.org/FHS_3.0/fhs/ch05s13.html

	viper.SetConfigName("gateway")
	viper.SetConfigType("ini")

	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	if configPath, success := os.LookupEnv("CASAOS_CONFIG_PATH"); success {
		viper.AddConfigPath(configPath)
	}

	viper.AddConfigPath(currentDirectory)
	viper.AddConfigPath(filepath.Join(currentDirectory, "conf"))
	viper.AddConfigPath(filepath.Join("/", "etc", "casaos"))

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := state.SetRuntimePath(viper.GetString(configKeyRuntimePath)); err != nil {
		return err
	}

	if err := state.SetGatewayPort(viper.GetString(configKeyGatewayPort)); err != nil {
		return err
	}

	return nil
}

func Save(state *service.State) error {
	viper.Set(configKeyGatewayPort, state.GetGatewayPort())
	viper.Set(configKeyRuntimePath, state.GetRuntimePath())

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}
