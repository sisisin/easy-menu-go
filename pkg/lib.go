package pkg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sisisin/easy-menu-go/pkg/collection"
	"sisisin/easy-menu-go/pkg/command"
	m "sisisin/easy-menu-go/pkg/menu"
	"sisisin/easy-menu-go/pkg/ui"
	"strconv"

	"gopkg.in/yaml.v3"
)

func Run(document *yaml.Node) {
	m := m.ParseMenu(document.Content[0])
	props := getViewProps(*m, nil)
	ui.RenderMenu(props)

	processEasyMenu(*m)
}

func getCurrent(rootMenu m.MenuItem, cursor []int64) m.MenuItem {
	target := rootMenu
	for _, v := range cursor {
		target = target.SubMenu.Items[v]
	}
	return target
}

func getViewProps(target m.MenuItem, state *command.CommandState) ui.ViewProps {
	if state != nil && state.ProcessState != command.NotExecuting {
		viewType := ui.CommandResult
		return ui.ViewProps{
			ViewType:     viewType,
			Title:        target.Name,
			Command:      target.Command.Command,
			CommandState: *state,
		}
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

func processEasyMenu(menu m.MenuItem) {
	var currentViewProps ui.ViewProps
	cursor := []int64{}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		state := command.CommandState{
			ProcessState: command.NotExecuting,
			Err:          nil,
		}

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
					if currentViewProps.CommandState.ProcessState == command.Failed || currentViewProps.CommandState.ProcessState == command.Succeeded {
						cursor = cursor[:len(cursor)-1]
						break
					}
					switch in {
					case "q":
						fmt.Println("received `q`, exit.")
						os.Exit(0)
					case "b":
						if len(cursor) > 0 {
							cursor = cursor[:len(cursor)-1]
						}
					case "n":
						if currentViewProps.ViewType == ui.Confirm {
							cursor = cursor[:len(cursor)-1]
						}
					case "y":
						if currentViewProps.ViewType == ui.Confirm {
							current := getCurrent(menu, cursor)
							state = command.ExecuteCommand(current)
						}
					}
				}
			}
		} else {
			current := getCurrent(menu, cursor)
			idx := int(num - 1)
			if current.Kind == m.SubMenu && 0 <= idx && idx < len(current.SubMenu.Items) {
				cursor = append(cursor, num-1)
			}
		}

		currentViewProps = getViewProps(getCurrent(menu, cursor), &state)
		ui.RenderMenu(currentViewProps)
	}
}

func Check(e error) {
	if e != nil {
		log.Fatalln(e)
		panic(e)
	}
}
