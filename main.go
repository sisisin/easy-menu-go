package main

import (
	"flag"
	"os"

	"github.com/sisisin/easy-menu-go/pkg"

	"gopkg.in/yaml.v3"
)

func main() {
	var document yaml.Node

	var configFlag string
	flag.StringVar(&configFlag, "config", "", "config file path")
	flag.Parse()

	configFile, err := pkg.LoadConfig(&document, configFlag)
	pkg.Check(err)
	pkg.Run(&document, configFile)

	os.Exit(0)
}
