package main

import (
	"github.com/1llusion1st/go-transport-manager/cmd/commands"
	"github.com/alecthomas/kong"
	"github.com/sirupsen/logrus"
)

var CLI struct {
	Debug   bool                    `help:"Enable debug mode."`
	Forward commands.ForwardCmd     `cmd:"" help:"forward connection"`
	Reserve commands.ReserveCommand `cmd:"" help:"reserve service"`
}

func main() {
	ctx := kong.Parse(&CLI)
	if CLI.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugf("debug mode activated")
	}
	err := ctx.Run(&commands.Context{Debug: CLI.Debug})
	ctx.FatalIfErrorf(err)
}
