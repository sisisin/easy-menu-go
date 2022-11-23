package ui

import (
	"fmt"
	"sisisin/easy-menu-go/pkg/command"
)

type ViewType uint32

const (
	List ViewType = 1 << iota
	Confirm
	Unsupported
	CommandResult
)

type ViewProps struct {
	ViewType     ViewType
	Title        string
	List         []string
	Command      string
	CommandState command.CommandState
}

func RenderMenu(props ViewProps) {
	fmt.Print("\033[H\033[2J")
	// fmt.Println()
	// fmt.Println("*** *** *** *** *** *** ***")
	// fmt.Println()

	switch props.ViewType {
	case Unsupported:
		fmt.Println("> ====================== <")
		fmt.Println(props.Title)
		fmt.Println("---------------------")
		fmt.Println("Unsupported config")
		fmt.Println("> ====================== <")
	case List:
		fmt.Println("> ====================== <")
		fmt.Println(props.Title)
		fmt.Println("---------------------")
		for i, v := range props.List {
			fmt.Printf("[%d] %v\n", i+1, v)
		}
		fmt.Println("> ====================== <")

	case Confirm:
		fmt.Println("> ====================== <")
		fmt.Println(props.Title)
		fmt.Println("---------------------")
		fmt.Printf("execute `%v` [y/n]\n", props.Command)
		fmt.Println("> ====================== <")
	case CommandResult:
		if props.CommandState.ProcessState == command.Failed {
			renderError(props)
		} else {
			renderCommandSucceeded(props)
		}
	}
}

func renderError(props ViewProps) {
	fmt.Println("> ====================== <")
	fmt.Println(props.Title)
	fmt.Println("---------------------")
	fmt.Printf("failed `%v`\n", props.Command)
	fmt.Println(props.CommandState.Err)
	fmt.Println("> ====================== <")
}

func renderCommandSucceeded(props ViewProps) {
	fmt.Println("> ====================== <")
	fmt.Println("Executing: ", props.Title)
	fmt.Println("---------------------")
	fmt.Printf("succeeded `%v`\n", props.Command)
	fmt.Println("> ====================== <")
	fmt.Println("press any key to back menu")
}
