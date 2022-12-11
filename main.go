package main

import (
	"os"

	"github.com/sisisin/easy-menu-go/pkg"

	"gopkg.in/yaml.v3"
)

func main() {
	var document yaml.Node

	err := pkg.LoadConfig(&document)
	pkg.Check(err)
	pkg.Run(&document)

	os.Exit(0)
}
