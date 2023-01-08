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

var getwd = os.Getwd

func validConfigPathOrExit(configPathFromArg string) string {
	wd, err := getwd()
	checkError(err)

	if configPathFromArg == "" {
		configFile := filepath.Join(wd, defaultConfigPath1)

		f, err := os.Open(configFile)
		if err != nil {
			configFile = filepath.Join(wd, defaultConfigPath2)
		}

		exitIfPathError(err)
		defer f.Close()
		return configFile
	} else {
		configFile := filepath.Join(wd, configPathFromArg)
		f, err := os.Open(configFile)
		exitIfPathError(err)
		defer f.Close()
		return configFile
	}

}

func exitIfPathError(err error) {
	switch castedErr := err.(type) {
	case *fs.PathError:
		exit(1, fmt.Sprintf("Error: cannot read config file `%s`\n%s\n", castedErr.Path, castedErr.Err))
	}
}

var exit = func(exitCode int, msg string) {
	fmt.Printf("%s", msg)
	os.Exit(exitCode)
}

func checkError(e error) {
	if e != nil {
		log.Fatalln(e)
		panic(e)
	}
}
