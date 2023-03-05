package ui

import (
	"fmt"
	"github.com/sisisin/easy-menu-go/pkg/args"
	"github.com/sisisin/easy-menu-go/pkg/command"
	"os"
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
	debug := args.GetEnvs().Debug

	if debug {
		fmt.Println()
		fmt.Println("*** *** *** *** *** *** ***")
		fmt.Println()
	} else {
		fmt.Print("\033[H\033[2J")
	}

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
		fmt.Println("execute [y/n]")
		fmt.Println("> ====================== <")
	case CommandResult:
		renderExecuting(props)
	}
}

func renderExecuting(props ViewProps) {
	if props.CommandState.Err != nil {
		printFail(props, props.CommandState.Err)
		return
	}
	fmt.Println("> ====================== <")
	fmt.Println("Executing: ", props.Title)
	fmt.Println("---------------------")

	cmd := props.CommandState.Cmd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		printFail(props, err)
		return
	}

	if err := cmd.Wait(); err != nil {
		printFail(props, err)
		return
	}

	fmt.Println("---------------------")
	fmt.Println("succeeded, exitCode: ", cmd.ProcessState.ExitCode())
	fmt.Println("> ====================== <")
	fmt.Println("press any key to back menu")
}

func printFail(props ViewProps, err error) {
	fmt.Println("---------------------")
	fmt.Println("error: '", err, "'")
	fmt.Println("> ====================== <")
}
