package menu

import (
	"fmt"
	"sisisin/easy-menu-go/pkg/collection"

	"gopkg.in/yaml.v3"
)

type MenuItemKind uint32

const (
	SubMenu MenuItemKind = 1 << iota
	Command
	EvalMenu
)

type MenuItem struct {
	Kind MenuItemKind

	Name    string
	WorkDir string
	Env     map[string]string

	SubMenu *MenuConfiguration
	Command *CommandSpec
}
type MenuConfiguration struct {
	Items []MenuItem
}
type CommandSpec struct {
	Command string
}

func ParseMenu(node *yaml.Node) *MenuItem {
	if node.Kind != yaml.MappingNode {
		panic("node must be MappingNode")
	}

	key := node.Content[0]
	value := node.Content[1]

	switch value.Kind {
	case yaml.SequenceNode:
		return &MenuItem{
			Kind:    SubMenu,
			Name:    key.Value,
			Env:     map[string]string{},
			WorkDir: "todo",
			SubMenu: &MenuConfiguration{
				Items: collection.Map(value.Content, func(v *yaml.Node, _ int) MenuItem {
					return *ParseMenu(v)
				}),
			},
		}
	case yaml.ScalarNode:
		return &MenuItem{
			Kind:    Command,
			Name:    key.Value,
			Env:     map[string]string{},
			WorkDir: "todo",
			Command: &CommandSpec{
				Command: value.Value,
			},
		}
	case yaml.MappingNode:
		entries := toEntries(value)
		idx := collection.FindIndex(entries, func(v mapEntry, _ int) bool {
			return v.key.Value == "eval"
		})
		if idx > -1 {
			// todo
			return &MenuItem{
				Kind: EvalMenu,
				Name: key.Value,
			}
		} else {
			item := &MenuItem{
				Kind: Command,
				Name: key.Value,
				Env:  map[string]string{},
			}
			for _, v := range entries {
				switch v.key.Value {
				case "env":
					// todo
				case "work_dir":
					item.WorkDir = v.value.Value
				case "run":
					item.Command = &CommandSpec{
						Command: v.value.Value,
					}
				}
			}
			return item
		}
	default:
		println("????????")
	}

	return nil
}

type mapEntry struct {
	key   yaml.Node
	value yaml.Node
}

func toEntries(node *yaml.Node) []mapEntry {
	if node.Kind != yaml.MappingNode {
		panic("node must be MappingNode")
	}

	entries := []mapEntry{}
	for i := 0; i < len(node.Content); i += 2 {
		entries = append(entries, mapEntry{
			key:   *node.Content[i],
			value: *node.Content[i+1],
		})
	}

	return entries
}

func debugPrint(node *yaml.Node) {
	var kind string
	switch node.Kind {
	case yaml.SequenceNode:
		kind = "sequence"
	case yaml.ScalarNode:
		kind = "scalar"
	case yaml.MappingNode:
		kind = "mapping"
	}

	var v any
	node.Decode(&v)
	y, err := yaml.Marshal(v)
	if err != nil {
		fmt.Println("failed marshal")
	}

	fmt.Println(kind, node.Value, node.Tag)
	fmt.Println(string(y))
}
