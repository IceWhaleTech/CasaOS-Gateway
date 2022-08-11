package main

import (
	interfaces "github.com/IceWhaleTech/CasaOS-Common"
	"github.com/IceWhaleTech/CasaOS-Common/utils/file"
	"github.com/IceWhaleTech/CasaOS-Common/utils/systemctl"
	"github.com/IceWhaleTech/CasaOS-Common/utils/version"
	"github.com/IceWhaleTech/CasaOS-Gateway/common"
	"gopkg.in/ini.v1"
)

type updater033to035 struct{}

func (u *updater033to035) IsMigrationNeeded() (bool, error) {
	_logger.Info("Checking if migration is needed for CasaoS version between 0.3.3 and 0.3.5...")

	minorVersion, err := version.DetectMinorVersion()
	if err != nil {
		return false, err
	}

	if minorVersion != 3 {
		return false, nil
	}

	isUserDataInDatabase, err := version.IsUserDataInDatabase()
	if err != nil {
		return false, err
	}

	if !isUserDataInDatabase {
		return false, nil
	}

	return true, nil
}

func (u *updater033to035) PreMigrate() error {
	_logger.Info("Executing steps before migration for CasaoS version between 0.3.3 and 0.3.5...")

	_logger.Info("Disabling legacy %s...", version.LegacyCasaOSServiceName)
	err := systemctl.DisableService(version.LegacyCasaOSServiceName)
	if err != nil {
		return err
	}

	_logger.Info("Copying %s to %s if it doesn't exist...", gatewayConfigSampleFilePath, gatewayConfigFilePath)
	if err := file.CopySingleFile(gatewayConfigSampleFilePath, gatewayConfigFilePath, "skip"); err != nil {
		return err
	}

	return nil
}

func (u *updater033to035) Migrate() error {
	_logger.Info("Executing migration steps for CasaoS version between 0.3.3 and 0.3.5...")

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

func (u *updater033to035) PostMigrate() error {
	_logger.Info("Executing steps after migration for CasaoS version between 0.3.3 and 0.3.5...")

	// enable new gateway service
	_logger.Info("Enabling %s...", gatewayServiceName)
	if err := systemctl.EnableService(gatewayServiceName); err != nil {
		return err
	}

	return nil
}

func NewUpdater033to035() interfaces.Updater {
	return &updater033to035{}
}
