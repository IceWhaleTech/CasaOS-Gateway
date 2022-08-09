package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/IceWhaleTech/CasaOS-Gateway/common"
)

func init() {
	versionFlag := flag.Bool("v", false, "version")
	flag.Parse()

	if *versionFlag {
		fmt.Println(common.Version)
		os.Exit(0)
	}
}

func main() {
	// TODO check if migration is needed
}
