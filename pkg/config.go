package pkg

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadConfig(out *yaml.Node, configPathFromArg string) (string, error) {
	configFile := validConfigPathOrExit(configPathFromArg)
	println(configFile)
	buf, err := os.ReadFile(configFile)
	Check(err)

	return configFile, yaml.Unmarshal(buf, out)
}

const (
	defaultConfigPath1 = "easy-menu.yml"
	defaultConfigPath2 = "easy-menu.yaml"
)

func validConfigPathOrExit(configPathFromArg string) string {
	wd, err := os.Getwd()
	Check(err)

	var configFile string
	if configPathFromArg == "" {
		configFile = filepath.Join(wd, defaultConfigPath1)

		f, err := os.Open(configFile)
		if err != nil {
			configFile = filepath.Join(wd, defaultConfigPath2)
			f, err = os.Open(configFile)
		}

		switch err := err.(type) {
		case *fs.PathError:
			fmt.Printf("Error: cannot read config file `%s`\n", err.Path)
			fmt.Printf("%s\n", err.Err)

			os.Exit(1)
		}
		defer f.Close()

		return configFile
	} else {
		configFile = filepath.Join(wd, configPathFromArg)
		f, err := os.Open(configFile)

		switch err := err.(type) {
		case *fs.PathError:
			fmt.Printf("Error: cannot read config file `%s`\n", err.Path)
			fmt.Printf("%s\n", err.Err)

			os.Exit(1)
		}

		defer f.Close()
		return configFile
	}
}
