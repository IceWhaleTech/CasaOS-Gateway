package main

import (
	"os"

	interfaces "github.com/IceWhaleTech/CasaOS-Common"
	"gopkg.in/ini.v1"
)

const casaOSConfFilePath = "/etc/casaos.conf"

type updater03x struct{}

func (u *updater03x) IsMigrationNeeded() bool {
	// check if file exists
	if _, err := os.Stat(casaOSConfFilePath); os.IsNotExist(err) {
		_logger.Debug("could not find " + casaOSConfFilePath)
		return false
	}

	cfg, err := ini.Load(casaOSConfFilePath)
	if err != nil {
		_logger.Debug("could not load " + casaOSConfFilePath)
		return false
	}

	return false
}

func (u *updater03x) PreMigrate() error {
	return nil
}

func (u *updater03x) Migrate() error {
	return nil
}

func (u *updater03x) PostMigrate() error {
	return nil
}

func NewUpdater03x() interfaces.Updater {
	return &updater03x{}
}
