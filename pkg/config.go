package pkg

import (
	"fmt"
	"log"
	"os"
	"path"
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
		configFile, ok := lookupConfig(wd)
		if !ok {
			exit(1, fmt.Sprintf("Error: cannot read config file: %s\n", configFile))
		}

		return configFile
	} else {
		configFile := filepath.Join(wd, configPathFromArg)
		if fileExists(configFile) {
			return configFile
		}
		exit(1, fmt.Sprintf("Error: cannot read config file: %s\n", configFile))
		// 到達しない
		return ""
	}
}

func lookupConfig(dir string) (string, bool) {
	configFile, ok := existsDefaultConfig(dir)
	if ok {
		return configFile, true
	} else {
		if dir == "/" {
			return "", false
		}
		nextDir, _ := path.Split(path.Clean(dir))
		return lookupConfig(nextDir)
	}
}

func existsDefaultConfig(basename string) (path string, ok bool) {
	f := filepath.Join(basename, defaultConfigPath1)
	if fileExists(f) {
		return f, true
	}

	f = filepath.Join(basename, defaultConfigPath2)
	if fileExists(f) {
		return f, true
	}

	return "", false
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
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
