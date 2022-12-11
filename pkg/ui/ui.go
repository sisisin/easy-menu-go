package ui

import (
	"fmt"
	"io"
	"os"
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
	// fmt.Print("\033[H\033[2J")
	fmt.Println()
	fmt.Println("*** *** *** *** *** *** ***")
	fmt.Println()

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
		renderExecuting(props)
	}
}

func renderExecuting(props ViewProps) {
	fmt.Println("> ====================== <")
	fmt.Println("Executing: ", props.Title)
	fmt.Println("---------------------")

	cmd := props.CommandState.Cmd

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		printFail(props, err)
		return
	}
	defer stdout.Close()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		printFail(props, err)
		return
	}
	defer stderr.Close()

	if err := cmd.Start(); err != nil {
		printFail(props, err)
		return
	}

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	if err := cmd.Wait(); err != nil {
		printFail(props, err)
		return
	}

	fmt.Println("---------------------")
	fmt.Printf("succeeded `%v`\n", props.Command)
	fmt.Println("> ====================== <")
	fmt.Println("press any key to back menu")
}

func printFail(props ViewProps, err error) {
	fmt.Println("---------------------")
	fmt.Println("error: '", err, "'")
	fmt.Println("> ====================== <")
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
