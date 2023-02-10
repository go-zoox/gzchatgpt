package main

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/gzchatgpt/commands"
)

func main() {
	app := cli.NewMultipleProgram(&cli.MultipleProgramConfig{
		Name:    "gzchatgpt",
		Usage:   "gzchatgpt is a portable chatgpt server",
		Version: Version,
	})

	commands.RegistryRunner(app)
	commands.RegistryServer(app)

	commands.RegistrFeishuBot(app)

	app.Run()
}
