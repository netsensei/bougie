package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/netsensei/bougie/config"
	"github.com/netsensei/bougie/tui"
)

func main() {
	filePath := flag.String("file", "", "path to a file to be opened in bougie")

	flag.Parse()

	err := config.Init()
	if err != nil {
		fmt.Println("Error: error during startup:", err)
		os.Exit(1)
	}

	if err := tui.Start(*filePath); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
