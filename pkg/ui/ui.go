package ui

import (
	"fmt"
	"os"
)

type ViewType uint32

const (
	List ViewType = 1 << iota
	Confirm
)

type ViewProps struct {
	ViewType ViewType
	Title    string
	List     []string
	Command  string
}

func RenderMenu(props ViewProps) {
	// fmt.Print("\033[H\033[2J")
	fmt.Println("*** *** *** *** *** *** ***")
	switch props.ViewType {
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
	}
}

func RenderEtc(in string) {
	if in == "q" {
		fmt.Println("received `q`, exit.")
		os.Exit(0)
	} else {
		fmt.Printf("received `%v`\n", in)
	}
}
