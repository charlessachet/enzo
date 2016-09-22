package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func ls(path string) {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read dir: %s\n", err.Error())
		os.Exit(1)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
}

func main() {
	path := "."

	if len(os.Args) > 1 {
		if os.Args[1][0] != '-' {
			path = os.Args[1]
		} else {
			fmt.Fprintf(os.Stderr, "%s not implemented\n", os.Args[1])
			os.Exit(1)
		}
	}

	ls(path)
}
