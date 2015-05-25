package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "mutter"
	app.Version = Version
	app.Usage = ""
	app.Author = "adwd"
	app.Email = "masahiro.nishida@bizreach.co.jp"
	app.Commands = Commands

	app.Run(os.Args)
}
