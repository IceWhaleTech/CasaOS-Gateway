package main

import (
	"flag"
	"fmt"
	"os"

	interfaces "github.com/IceWhaleTech/CasaOS-Common"
	"github.com/IceWhaleTech/CasaOS-Gateway/common"
)

var _logger *Logger

func init() {
	versionFlag := flag.Bool("v", false, "version")
	debugFlag := flag.Bool("d", true, "debug")
	flag.Parse()

	if *versionFlag {
		fmt.Println(common.Version)
		os.Exit(0)
	}

	if os.Getuid() != 0 {
		_logger.Info("Root privileges are required to run this program.")
		os.Exit(1)
	}

	_logger = NewLogger()

	if *debugFlag {
		_logger.DebugMode = true
	}
}

func main() {
	updaters := []interfaces.Updater{
		NewUpdater033to035(),
	}

	for _, updater := range updaters {
		migrationNeeded, err := updater.IsMigrationNeeded()
		if err != nil {
			panic(err)
		}

		if !migrationNeeded {
			continue
		}

		if err := updater.PreMigrate(); err != nil {
			panic(err)
		}

		if err := updater.Migrate(); err != nil {
			panic(err)
		}

		if err := updater.PostMigrate(); err != nil {
			panic(err)
		}

	}
}
