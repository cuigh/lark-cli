package main

import (
	"github.com/cuigh/auxo/app"
	"github.com/cuigh/auxo/app/flag"
	"github.com/cuigh/auxo/config"
	"github.com/cuigh/lark-cli/cmd"
)

func main() {
	config.SetDefaultValue("banner", false)
	app.Name = "lark-cli"
	app.Version = "0.9"
	app.Desc = "A tool for developing lark based application"
	app.Action = cmd.Root
	app.AddCommand(cmd.New())
	app.AddCommand(cmd.Gen())
	app.Flags.Register(flag.Help | flag.Version)
	app.Start()
}
