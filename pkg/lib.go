package pkg

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func Run(document *yaml.Node) {
	menuConfig := document.Content[0]

	var meta Meta
	var menu *yaml.Node
	var menuTitle string

	for i := 0; i < len(menuConfig.Content); i += 2 {
		key := menuConfig.Content[i]
		value := menuConfig.Content[i+1]

		if key.Value == "meta" {
			meta = parseMeta(value)
		} else {
			menu = value
			menuTitle = key.Value
		}
	}

	// todo
	_ = meta

	// currentMenu := visitMenu(menu, []int{})
	// debugPrint(menu, []int{})
	props := ViewProps{}
	props.title = menuTitle
	for _, v := range menu.Content {
		props.list = append(props.list, v.Content[0].Value)
	}
	props.breadcrumb = []string{}

	printMenu(props)
}

func toViewProps(node *yaml.Node) ViewProps {
	p := ViewProps{}
	p.title = node.Content[0].Value

	for _, v := range node.Content[1].Content {
		p.list = append(p.list, v.Value)
	}
	p.breadcrumb = []string{}

	return p
}

func printMenu(props ViewProps) {
	fmt.Println("> ====================== <")
	fmt.Println(props.title)
	fmt.Println("---------------------")
	// depth := strings.Trim(strings.Join(props.breadcrumb, " > "), "[]")
	// fmt.Println(depth)
	for i, v := range props.list {
		fmt.Printf("[%d] %v\n", i, v)
	}
	fmt.Println("> ====================== <")
}

type ViewProps struct {
	title      string
	list       []string
	breadcrumb []string
}
type VisitResult struct {
	props  ViewProps
	cursor []int
}

func visitMenu(node *yaml.Node, cursor []int) *yaml.Node {
	if len(cursor) == 0 {
		return node
	} else {
		index := cursor[0]
		next := cursor[1:]
		return visitMenu(node.Content[index], next)
	}
}

type Meta struct {
	workDir string
	env     map[string]string
	lock    bool
}

func parseMeta(node *yaml.Node) Meta {
	meta := Meta{}
	return meta
}

func debugPrint(node *yaml.Node, cursor []int) {
	depth := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(cursor)), " > "), "[]")
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
	fmt.Println(kind, node.Value, node.Tag, depth)
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func VisitNode(node *yaml.Node, cursor []int) {
	if len(node.Content) == 0 {
		return
	} else {
		for i, v := range node.Content {
			debugPrint(v, cursor)
			VisitNode(v, append(cursor, i))
		}
	}
}
