package main

import (
	"os"
	"time"

	interfaces "github.com/IceWhaleTech/CasaOS-Common"
	"github.com/IceWhaleTech/CasaOS-Common/utils/file"
	"github.com/IceWhaleTech/CasaOS-Common/utils/version"
	"github.com/IceWhaleTech/CasaOS-Gateway/common"
	"github.com/spf13/viper"
	"gopkg.in/ini.v1"
)

type migrationTool struct{}

var _status *version.GlobalMigrationStatus

func (u *migrationTool) IsMigrationNeeded() (bool, error) {
	if status, err := version.GetGlobalMigrationStatus("gateway"); err == nil {
		_status = status
		if status.LastMigratedVersion != "" {
			if r, err := version.Compare(status.LastMigratedVersion, common.Version); err == nil {
				return r < 0, nil
			}
		}
	}

	if _, err := os.Stat(version.LegacyCasaOSConfigFilePath); err != nil {
		_logger.Info("`%s` not found, migration is not needed.", version.LegacyCasaOSConfigFilePath)
		return false, nil
	}

	majorVersion, minorVersion, patchVersion, err := version.DetectLegacyVersion()
	if err != nil {
		if err == version.ErrLegacyVersionNotFound {
			return false, nil
		}

		return false, err
	}

	if majorVersion > 0 {
		return false, nil
	}

	if minorVersion > 3 {
		return false, nil
	}

	if minorVersion == 3 && patchVersion > 5 {
		return false, nil
	}

	_logger.Info("Migration is needed for a CasaOS version 0.3.5 and older...")
	return true, nil
}

func (u *migrationTool) PreMigrate() error {
	if _, err := os.Stat(gatewayConfigDirPath); os.IsNotExist(err) {
		_logger.Info("Creating %s since it doesn't exists...", gatewayConfigDirPath)
		if err := os.Mkdir(gatewayConfigDirPath, 0o755); err != nil {
			return err
		}
	}

	if _, err := os.Stat(gatewayConfigFilePath); os.IsNotExist(err) {
		_logger.Info("Creating %s since it doesn't exist...", gatewayConfigFilePath)

		f, err := os.Create(gatewayConfigFilePath)
		if err != nil {
			return err
		}

		defer f.Close()

		if _, err := f.WriteString(_gatewayINISample); err != nil {
			return err
		}
	}

	extension := "." + time.Now().Format("20060102") + ".bak"

	_logger.Info("Creating a backup %s if it doesn't exist...", version.LegacyCasaOSConfigFilePath+extension)
	if err := file.CopySingleFile(version.LegacyCasaOSConfigFilePath, version.LegacyCasaOSConfigFilePath+extension, "skip"); err != nil {
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

	newConfigFile, err := common.LoadConfig()
	if err != nil {
		return err
	}

	if err := migrateAppSection(legacyConfigFile, newConfigFile); err != nil {
		return err
	}

	if err := migrateServerSection(legacyConfigFile, newConfigFile); err != nil {
		return err
	}

	_logger.Info("Saving new config ...")
	return newConfigFile.WriteConfig()
}

func (u *migrationTool) PostMigrate() error {
	return _status.Done(common.Version)
}

func NewMigrationToolFor035AndOlder() interfaces.MigrationTool {
	return &migrationTool{}
}

func migrateAppSection(legacyConfigFile *ini.File, newConfigFile *viper.Viper) error {
	// LogPath
	if logPath, err := legacyConfigFile.Section("app").GetKey("LogPath"); err == nil {
		_logger.Info("[app] LogPath = %s", logPath.Value())
		newConfigFile.Set(common.ConfigKeyLogPath, logPath.Value())
	}

	if logPath, err := legacyConfigFile.Section("app").GetKey("LogSavePath"); err == nil {
		_logger.Info("[app] LogPath = %s", logPath.Value())
		newConfigFile.Set(common.ConfigKeyLogPath, logPath.Value())
	}

	// LogFileExt
	if logFileExt, err := legacyConfigFile.Section("app").GetKey("LogFileExt"); err == nil {
		_logger.Info("[app] LogFileExt = %s", logFileExt.Value())
		newConfigFile.Set(common.ConfigKeyLogFileExt, logFileExt.Value())
	}

	return nil
}

func migrateServerSection(legacyConfigFile *ini.File, newConfigFile *viper.Viper) error {
	if httpPort, err := legacyConfigFile.Section("server").GetKey("HttpPort"); err == nil {
		_logger.Info("[server] HttpPort = %s", httpPort.Value())
		newConfigFile.Set(common.ConfigKeyGatewayPort, httpPort.Value())
	}

	return nil
}
