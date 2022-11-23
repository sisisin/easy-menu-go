package pkg

import (
	"bufio"
	"fmt"
	"os"
	"sisisin/easy-menu-go/pkg/collection"
	m "sisisin/easy-menu-go/pkg/menu"
	"sisisin/easy-menu-go/pkg/ui"
	"strconv"

	"gopkg.in/yaml.v3"
)

func Run(document *yaml.Node) {
	m := m.ParseMenu(document.Content[0])
	cursor := []int64{}
	props := getViewProps(*m, cursor)
	ui.RenderMenu(props)

	processEasyMenu(*m, &cursor)
}

func getViewProps(menu m.MenuItem, cursor []int64) ui.ViewProps {
	target := menu
	for _, v := range cursor {
		target = target.SubMenu.Items[v]
	}
	switch target.Kind {
	case m.SubMenu:
		return ui.ViewProps{
			ViewType: ui.List,
			Title:    target.Name,
			List: collection.Map(target.SubMenu.Items, func(v m.MenuItem, _ int) string {
				return v.Name
			}),
		}
	case m.Command:
		return ui.ViewProps{
			ViewType: ui.Confirm,
			Title:    target.Name,
			Command:  target.Command.Command,
		}
	default:
		return ui.ViewProps{
			ViewType: ui.Unsupported,
			Title:    target.Name,
		}
	}
}
func processEasyMenu(menu m.MenuItem, cursor *[]int64) {
	scanner := bufio.NewScanner(os.Stdin)
	var currentViewProps ui.ViewProps

	for {
		fmt.Print("> ")
		scanner.Scan()
		in := scanner.Text()
		num, err := strconv.ParseInt(in, 10, 0)

		if err != nil {
			if enum, ok := err.(*strconv.NumError); ok {
				switch enum.Err {
				case strconv.ErrRange:
					// no-op
				case strconv.ErrSyntax:
					switch in {
					case "q":
						fmt.Println("received `q`, exit.")
						os.Exit(0)
					}
				}
			}
		} else {
			*cursor = append(*cursor, num-1)
			currentViewProps = getViewProps(menu, *cursor)
			ui.RenderMenu(currentViewProps)
		}
	}
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
