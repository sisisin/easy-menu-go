package main

import (
	"bufio"
	"fmt"
	"os"

	"sisisin/easy-menu-go/pkg"

	"gopkg.in/yaml.v3"
)

func main() {
	var document yaml.Node

	err := pkg.LoadConfig(&document)
	pkg.Check(err)
	pkg.Run(&document)

	// visitNode(&document, []int{})
	os.Exit(0)
}

func run(node *yaml.Node, cursor []int) {
		switch node.Kind {
		case yaml.SequenceNode:
		case yaml.ScalarNode:
			var v any
			node.Decode(&v)
			fmt.Println("Scalar")
			fmt.Println(v)
		case yaml.MappingNode:
			if len(cursor) == 0 {
				fmt.Println(node.Content[0].Value)
			}
		}
}

func printMenu(header string, contents []string){
	fmt.Println("=========>")
	fmt.Println("----------")
	fmt.Println(header)
	fmt.Println("----------")
	for i, c:=range contents {
		fmt.Printf("[%d] %v", i,c)
	}
	fmt.Println("<=========")
}

func p(t *yaml.Node) {
	var topMenus []string
	for _, v := range t.Content {
		topMenus = append(topMenus, v.Value)
	}

	menuString := "---> \n"
	for _, m := range topMenus {
		menuString += m + "\n"
	}
	menuString += "<---\n"
	fmt.Printf("%v", menuString)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("> ")
	var in string
	for {
		scanner.Scan()
		in = scanner.Text()

		if in == "q" {
			fmt.Println("received `q`, exit.")
			os.Exit(0)
		} else {
			fmt.Printf("received `%v`\n", in)
		}
	}
}

