package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	configKeyGatewayPort = "gateway.Port"
	configKeyRuntimePath = "common.RuntimePath"
)

type Config struct {
	gatewayPort string
	runtimePath string
	onChange    []func(*Config) error
}

func NewConfig() *Config {
	return &Config{
		gatewayPort: "",
		runtimePath: "",
		onChange:    make([]func(*Config) error, 0),
	}
}

func (c *Config) SetGatewayPort(port string) error {
	c.gatewayPort = port
	return c.change()
}

func (c *Config) GetGatewayPort() string {
	return c.gatewayPort
}

func (c *Config) SetRuntimePath(path string) error {
	c.runtimePath = path
	return c.change()
}

func (c *Config) GetRuntimePath() string {
	return c.runtimePath
}

func (c *Config) OnChange(f func(*Config) error) {
	c.onChange = append(c.onChange, f)
}

func (c *Config) change() error {
	for _, f := range c.onChange {
		if err := f(c); err != nil {
			return err
		}
	}

	return nil
}

func Load(config *Config) error {
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

	if err := config.SetRuntimePath(viper.GetString(configKeyRuntimePath)); err != nil {
		return err
	}

	if err := config.SetGatewayPort(viper.GetString(configKeyGatewayPort)); err != nil {
		return err
	}

	return nil
}
