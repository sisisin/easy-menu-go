package pkg

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)


func LoadConfig(out *yaml.Node) (err error) {
	wd, err := os.Getwd()
	Check(err)

	configFile := filepath.Join(wd, "easy-menu.yml")

	f, err := os.Open(configFile)
	Check(err)
	defer f.Close()

	buf, err := os.ReadFile(configFile)
	Check(err)

	return yaml.Unmarshal(buf, out)
}

