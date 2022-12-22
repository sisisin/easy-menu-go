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

	err := pkg.LoadConfig(&document, configFlag)
	pkg.Check(err)
	pkg.Run(&document)

	os.Exit(0)
}
