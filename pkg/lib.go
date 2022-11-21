package pkg

import (
	"bufio"
	"fmt"
	"os"
	"sisisin/easy-menu-go/pkg/collection"
	"sisisin/easy-menu-go/pkg/ui"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func Run(document *yaml.Node) {
	config := document.Content[0]

	var meta Meta
	var menuMap *yaml.Node
	var menuTitle string

	for i := 0; i < len(config.Content); i += 2 {
		key := config.Content[i]
		value := config.Content[i+1]

		if key.Value == "meta" {
			meta = parseMeta(value)
		} else {
			menuMap = value
			menuTitle = key.Value
		}
	}

	// todo
	_ = meta

	props := ui.ViewProps{
		ViewType: ui.List,
		Title:    menuTitle,
		List: collection.Map(menuMap.Content, func(v *yaml.Node, _ int) string {
			return v.Content[0].Value
		}),
	}

	ui.RenderMenu(props)

	cursor := []int64{}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		in := scanner.Text()

		num, err := strconv.ParseInt(in, 10, 0)

		if err != nil {
			if enum, ok := err.(*strconv.NumError); ok {
				switch enum.Err {
				case strconv.ErrRange:
					// no-op
				case strconv.ErrSyntax:
					ui.RenderEtc(in)
				}
			}
		}

		cursor = append(cursor, num)
		selected := getSelectedNodeByCursor(menuMap, cursor)
		props := toViewProps(selected)
		ui.RenderMenu(props)
	}
}

func getSelectedNodeByCursor(parent *yaml.Node, cursor []int64) *yaml.Node {
	result := parent
	for i := 0; i < len(cursor); i++ {
		currentCursor := cursor[i]
		v := result.Content[currentCursor]
		if v.Kind == yaml.MappingNode {
			return v
		} else if v.Kind == yaml.SequenceNode {
			next := append(cursor, int64(i))
			return getSelectedNodeByCursor(parent, next)
		}
	}
	return nil
}

func toViewProps(node *yaml.Node) ui.ViewProps {
	key := node.Content[0]
	val := node.Content[1]

	switch val.Kind {

	case yaml.SequenceNode:
		/*
			# ∨∨∨∨∨∨∨∨∨∨∨∨∨∨ - node(argument)
			- Menu List Item:
				- command1: echo 1
				# ∨∨∨∨∨∨∨∨ - node.Content[1].Content[1].Content[0].Value
				- command2: echo 2
		*/
		return ui.ViewProps{
			ViewType: ui.List,
			Title:    key.Value,
			List: collection.Map(val.Content, func(v *yaml.Node, _ int) string {
				return v.Content[0].Value
			}),
		}
	case yaml.MappingNode:
		return ui.ViewProps{
			ViewType: ui.Confirm,
			Title:    val.Content[0].Value,
			Command:  val.Content[1].Value,
		}
	case yaml.ScalarNode:
		return ui.ViewProps{
			ViewType: ui.Confirm,
			Title:    key.Value,
			Command:  val.Value,
		}
	}
	return ui.ViewProps{}
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
	fmt.Println(kind, node.Value, node.Tag, depth, v)
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
