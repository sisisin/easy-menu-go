package menu

import (
	"fmt"
	"os"

	"github.com/sisisin/easy-menu-go/pkg/collection"

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

type mappedNode struct {
	menu    mapEntry
	env     mapEntry
	workDir mapEntry
}

func validatedNode(node *yaml.Node) mappedNode {
	if node.Kind != yaml.MappingNode {
		panic("node must be MappingNode")
	}

	var env *mapEntry = nil
	var workDir *mapEntry = nil
	var menu *mapEntry = nil

	nodeEntries := toEntries(node)

	for i := 0; i < len(nodeEntries); i++ {
		v := nodeEntries[i]
		switch v.key.Value {
		case "env":
			if env == nil {
				env = &v
			} else {
				fmt.Printf("error: duplicate env keys. env must be one.")
				os.Exit(1)
			}
		case "work_dir":
			if workDir == nil {
				workDir = &v
			} else {
				fmt.Printf("error: duplicate work_dir keys. work_dir must be one.")
				os.Exit(1)
			}
		default:
			if menu == nil {
				menu = &v
			} else {
				fmt.Printf("error: menu definition must be one. keys: `%s`, `%s`\n", menu.key.Value, v.key.Value)
				os.Exit(1)
			}
		}
	}

	if menu == nil {
		fmt.Printf("error: menu definition must be exist\n")
		os.Exit(1)
	}
	ret := mappedNode{
		menu: *menu,
	}
	if workDir != nil {
		ret.workDir = *workDir
	}
	if env != nil {
		ret.env = *env
	}

	return ret
}

func ParseMenu(node *yaml.Node) *MenuItem {
	n := validatedNode(node)

	key := n.menu.key
	value := n.menu.value

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
		entries := toEntries(&value)
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

	var entries []mapEntry
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
