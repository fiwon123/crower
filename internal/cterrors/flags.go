package cterrors

import "fmt"

func PrintFileAndFolderFlagsNotUsed() {
	fmt.Println("file and folder flag not used")
}

func PrintNotFileAndOutputPath() {
	fmt.Println("needs to specify file path and out folder")
}

func PrintNotArgs(msg string) {
	if msg == "" {
		fmt.Println("need to pass arguments")
	} else {
		fmt.Printf("need to pass arguments: %s \n", msg)
	}
}

func PrintEmptyPaths() {
	fmt.Println("empty paths")
}
