package main

import (
	"github.com/netsensei/bougie/config"
	"github.com/netsensei/bougie/tui"
)

func main() {
	config.Init()
	tui.Start()
}
