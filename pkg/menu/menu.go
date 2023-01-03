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
	menu    *mapEntry
	run     *mapEntry
	env     mapEntry
	workDir string
	eval    *mapEntry
}

func validatedNode(node *yaml.Node) mappedNode {
	if node.Kind != yaml.MappingNode {
		panic("node must be MappingNode")
	}

	var env *mapEntry = nil
	var workDir *mapEntry = nil
	var menu *mapEntry = nil
	var run *mapEntry = nil
	var eval *mapEntry = nil

	nodeEntries := toEntries(node)

	for i := 0; i < len(nodeEntries); i++ {
		v := nodeEntries[i]
		switch v.key.Value {
		case "env":
			if env == nil {
				env = &v
			} else {
				fmt.Printf("error: duplicate keys. key `env` must be one.")
				os.Exit(1)
			}
		case "work_dir":
			if workDir == nil {
				workDir = &v
			} else {
				fmt.Printf("error: duplicate keys. key `work_dir` must be one.")
				os.Exit(1)
			}
		case "run":
			if run == nil {
				run = &v
			} else {
				fmt.Printf("error: duplicate keys. key `run` must be one.")
				os.Exit(1)
			}
		case "eval":
			if eval == nil {
				eval = &v
			} else {
				fmt.Printf("error: duplicate keys. key `eval` must be one.")
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

	if menu != nil && run != nil {
		fmt.Printf("error: menu definition or key `run` must be exist one, but exist both.\n")
		os.Exit(1)
	}
	/*
		todo: validation
			- eval/runは片方だけ有効
	*/

	ret := mappedNode{}
	if menu != nil {
		ret.menu = menu
	}
	if run != nil {
		ret.run = run
	}
	if eval != nil {
		ret.eval = eval
	}
	if workDir != nil {
		ret.workDir = (*workDir).value.Value
	}
	if env != nil {
		ret.env = *env
	}

	return ret
}

func ParseMenu(node *yaml.Node) *MenuItem {
	n := validatedNode(node)
	menuName := n.menu.key.Value
	value := n.menu.value

	switch value.Kind {
	case yaml.SequenceNode:
		return &MenuItem{
			Kind:    SubMenu,
			Name:    menuName,
			Env:     parseEnv(n.env.value),
			WorkDir: n.workDir,
			SubMenu: &MenuConfiguration{
				Items: collection.Map(value.Content, func(v *yaml.Node, _ int) MenuItem {
					return *ParseMenu(v)
				}),
			},
		}
	case yaml.ScalarNode:
		// todo: validate if Env,WorkDir are not empty
		c := factory.newSimpleCommand(menuName, value.Value)
		return &c

	case yaml.MappingNode:
		child := validatedNode(&value)

		if child.run == nil && child.eval == nil {
			fmt.Printf("error: run menu must have key `run` or `eval`\n")
			os.Exit(1)
		}

		if child.eval != nil {
			return &MenuItem{
				Kind: EvalMenu,
				Name: menuName,
			}
		}

		if child.run != nil {
			return &MenuItem{
				Kind:    Command,
				Name:    menuName,
				Env:     parseEnv(child.env.value),
				WorkDir: child.workDir,
				SubMenu: nil,
				Command: &CommandSpec{
					Command: child.run.value.Value,
				},
			}
		}

		fmt.Printf("error: command definition must have `eval` or `run` value")
		os.Exit(1)
	default:
		println("????????")
	}

	return nil
}

func parseEnv(envNode yaml.Node) map[string]string {
	// todo: validation

	var env map[string]string = nil
	if err := envNode.Decode(&env); err != nil {
		// todo: error message
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	return env
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

type factoryT struct{}

var factory factoryT

func (f factoryT) newSimpleCommand(name string, command string) MenuItem {
	return MenuItem{
		Kind:    Command,
		Name:    name,
		WorkDir: "",
		Env:     nil,
		Command: &CommandSpec{
			Command: command,
		},
	}
}
