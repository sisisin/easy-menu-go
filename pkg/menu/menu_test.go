package menu

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestParseMenu(t *testing.T) {
	var data = `
menu:
  - run: 'ls -al'
name: 'main menu'
env:
  ENV1: val_global_env
work_dir: w_dir
`

	node := loadNode(t, data)
	actual := ParseMenu(node)
	expected := MenuItem{
		Kind:    SubMenu,
		Name:    "main menu",
		WorkDir: "w_dir",
		Env: map[string]string{
			"ENV1": "val_global_env",
		},
		SubMenu: &MenuConfiguration{
			Items: []MenuItem{
				{
					Kind:    Command,
					Name:    "",
					WorkDir: "",
					Command: &CommandSpec{
						Command: "ls -al",
					},
				},
			},
		},
	}

	assert.Equal(t, expected, *actual)
}

func loadNode(t *testing.T, data string) *yaml.Node {
	var node yaml.Node

	if err := yaml.Unmarshal([]byte(data), &node); err != nil {
		t.Fatalf("%v", err)
	}

	return node.Content[0]
}
