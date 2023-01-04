package main

import (
	"flag"
	"os"

	"github.com/sisisin/easy-menu-go/pkg"
)

func main() {
	var configFlag string
	flag.StringVar(&configFlag, "config", "", "config file path")
	flag.Parse()

	configFile, document := pkg.LoadConfig(configFlag)
	pkg.Run(document, configFile)

	os.Exit(0)
}
