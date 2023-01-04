package pkg

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadConfig(configPathFromArg string) (string, *yaml.Node) {
	configFile := validConfigPathOrExit(configPathFromArg)
	buf, err := os.ReadFile(configFile)
	checkError(err)

	var node yaml.Node
	err = yaml.Unmarshal(buf, &node)
	checkError(err)

	return configFile, &node
}

const (
	defaultConfigPath1 = "easy-menu.yml"
	defaultConfigPath2 = "easy-menu.yaml"
)

func validConfigPathOrExit(configPathFromArg string) string {
	wd, err := os.Getwd()
	checkError(err)

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

func checkError(e error) {
	if e != nil {
		log.Fatalln(e)
		panic(e)
	}
}
