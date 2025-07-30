package main

import (
	"fmt"
	"os"

	"github.com/netsensei/bougie/config"
	"github.com/netsensei/bougie/tui"
)

func main() {
	err := config.Init()
	if err != nil {
		fmt.Println("Error during startup:", err)
		os.Exit(1)
	}

	tui.Start()
}
