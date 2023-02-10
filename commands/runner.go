package commands

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/gzchatgpt/runner"
)

func RegistryRunner(app *cli.MultipleProgram) {
	app.Register("runner", &cli.Command{
		Name:  "runner",
		Usage: "local runner",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "api-key",
				Usage:    "OpenAI API Key",
				EnvVars:  []string{"API_KEY"},
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			apiKey := ctx.String("api-key")
			return runner.Run(apiKey)
		},
	})
}
