package main

import (
	interfaces "github.com/IceWhaleTech/CasaOS-Common"
	"github.com/IceWhaleTech/CasaOS-Common/utils/file"
	"github.com/IceWhaleTech/CasaOS-Common/utils/version"
	"github.com/IceWhaleTech/CasaOS-Gateway/common"
	"gopkg.in/ini.v1"
)

type migrationTool struct{}

func (u *migrationTool) IsMigrationNeeded() (bool, error) {
	_logger.Info("Checking if migration is needed for CasaoS version between 0.3.3 and 0.3.5...")

	minorVersion, err := version.DetectMinorVersion()
	if err != nil {
		return false, err
	}

	if minorVersion != 3 {
		return false, nil
	}

	// this is the best way to tell if CasaOS version is between 0.3.3 and 0.3.5
	isUserDataInDatabase, err := version.IsUserDataInDatabase()
	if err != nil {
		return false, err
	}

	if !isUserDataInDatabase {
		return false, nil
	}

	return true, nil
}

func (u *migrationTool) PreMigrate() error {
	_logger.Info("Copying %s to %s if it doesn't exist...", gatewayConfigSampleFilePath, gatewayConfigFilePath)
	if err := file.CopySingleFile(gatewayConfigSampleFilePath, gatewayConfigFilePath, "skip"); err != nil {
		return err
	}
	return nil
}

func (u *migrationTool) Migrate() error {
	_logger.Info("Loading legacy %s...", version.LegacyCasaOSConfigFilePath)
	legacyConfigFile, err := ini.Load(version.LegacyCasaOSConfigFilePath)
	if err != nil {
		return err
	}

	key, err := legacyConfigFile.Section("server").GetKey("HttpPort")
	if err != nil {
		return err
	}

	httpPort := key.Value()

	newConfigFile, err := common.LoadConfig()
	if err != nil {
		return err
	}

	_logger.Info("Updating %s to be '%s' in %s...", common.ConfigKeyGatewayPort, httpPort, gatewayConfigFilePath)
	newConfigFile.Set(common.ConfigKeyGatewayPort, httpPort)
	return newConfigFile.WriteConfig()
}

func (u *migrationTool) PostMigrate() error {
	return nil
}

func NewMigrationToolFor033_034_035() interfaces.MigrationTool {
	return &migrationTool{}
}
