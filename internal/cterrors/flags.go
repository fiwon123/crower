package cterrors

import "fmt"

func PrintFileAndFolderFlagsNotUsed() {
	fmt.Println("file and folder flag not used")
}

func PrintNotFileAndOutputPath() {
	fmt.Println("needs to specify file path and out folder")
}

func PrintEmptyPaths() {
	fmt.Println("empty paths")
}
