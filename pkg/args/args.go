package args

import (
	"flag"
	"fmt"
	"os"
)

// set in build step
var (
	Version  = "in_development"
	Revision = "in_development"
)

func GetFlags() flagDefinition {
	var configFlag string
	flag.StringVar(&configFlag, "config", "", "config file path")

	var versionFlag bool
	flag.BoolVar(&versionFlag, "version", false, "show current version")

	flag.Parse()

	return flagDefinition{
		Config:  configFlag,
		Version: versionFlag,
	}
}

func GetVersionStr() string {
	return fmt.Sprintf("easy-menu-go version: %s, rev: %s", Version, Revision)
}

type flagDefinition struct {
	Config  string
	Version bool
}

type envDefinition struct {
	Debug bool
}

func GetEnvs() envDefinition {
	return envDefinition{
		Debug: os.Getenv("DEBUG") == "true",
	}
}
