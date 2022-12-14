package pkg

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/sisisin/easy-menu-go/pkg/args"
	"github.com/sisisin/easy-menu-go/pkg/collection"
	"github.com/sisisin/easy-menu-go/pkg/command"
	m "github.com/sisisin/easy-menu-go/pkg/menu"
	"github.com/sisisin/easy-menu-go/pkg/ui"
)

func Run() {
	flags := args.GetFlags()

	if flags.Version {
		fmt.Println(args.GetVersionStr())
		os.Exit(0)
	}

	configFile, document := LoadConfig(flags.Config)
	m := m.ParseMenu(document.Content[0])
	// v, _ := json.Marshal(*m)
	// println(string(v))
	props := getViewProps(*m, nil)
	ui.RenderMenu(props)

	processEasyMenu(*m, configFile)
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
				title := v.Name
				if title == "" && v.Kind == m.Command {
					title = v.Command.Command
				}
				return title
			}),
		}
	case m.Command:
		title := target.Name
		if title == "" {
			title = target.Command.Command
		}
		return ui.ViewProps{
			ViewType: ui.Confirm,
			Title:    title,
			Command:  target.Command.Command,
		}
	default:
		return ui.ViewProps{
			ViewType: ui.Unsupported,
			Title:    target.Name,
		}
	}
}

func processEasyMenu(menu m.MenuItem, configFile string) {
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

		if currentViewProps.ViewType == ui.CommandResult {
			cursor = cursor[:len(cursor)-1]
			currentViewProps = getViewProps(getCurrent(menu, cursor), &state)
			ui.RenderMenu(currentViewProps)
			continue
		}

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
							state = command.GetSelectedCommandState(menu, cursor, configFile)
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
