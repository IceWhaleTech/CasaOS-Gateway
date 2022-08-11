package main

import (
	interfaces "github.com/IceWhaleTech/CasaOS-Common"
	"github.com/IceWhaleTech/CasaOS-Common/utils/systemctl"
	"github.com/IceWhaleTech/CasaOS-Common/utils/version"
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
	_logger.Info("Executing steps before migration for CasaoS version between 0.3.3 to 0.3.5...")

	// disable legacy CasaOS service
	err := systemctl.DisableService(version.LegacyCasaOSServiceName)
	if err != nil {
		return err
	}

	// setup new gateway config file
	gatewayConfigSampleFilePath := "/etc/casaos/gateway.ini.sample"
	gatewayConfigFilePath := "/etc/casaos/gateway.ini"

	return nil
}

func (u *updater033to035) Migrate() error {
	_logger.Info("Executing migration steps for CasaoS version between 0.3.3 to 0.3.5...")

	// load config file
	configFile, err := ini.Load(version.LegacyCasaOSConfigFilePath)
	if err != nil {
		return err
	}

	key, err := configFile.Section("server").GetKey("HttpPort")
	if err != nil {
		return err
	}

	httpPort := key.Value()

	return nil
}

func (u *updater033to035) PostMigrate() error {
	_logger.Info("Executing steps after migration for CasaoS version between 0.3.3 to 0.3.5...")
	return nil
}

func NewUpdater033to035() interfaces.Updater {
	return &updater033to035{}
}
