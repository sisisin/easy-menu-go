package args

import (
	"flag"
	"os"
)

func GetFlags() flagDefinition {
	var configFlag string
	flag.StringVar(&configFlag, "config", "", "config file path")
	flag.Parse()

	return flagDefinition{
		Config: configFlag,
	}
}

type flagDefinition struct {
	Config string
}

type envDefinition struct {
	Debug bool
}

func GetEnvs() envDefinition {
	return envDefinition{
		Debug: os.Getenv("DEBUG") == "true",
	}
}
