package main

import (
	interfaces "github.com/IceWhaleTech/CasaOS-Common"
	"github.com/IceWhaleTech/CasaOS-Common/utils/version"
)

type updater033to035 struct{}

func (u *updater033to035) IsMigrationNeeded() (bool, error) {
	minorVersion, err := version.DetectMinorVersion()
	if err != nil {
		return false, err
	}

	if minorVersion != 2 {
		return false, nil
	}

	return true, nil
}

func (u *updater033to035) PreMigrate() error {
	return nil
}

func (u *updater033to035) Migrate() error {
	return nil
}

func (u *updater033to035) PostMigrate() error {
	return nil
}

func NewUpdater033to035() interfaces.Updater {
	return &updater033to035{}
}
