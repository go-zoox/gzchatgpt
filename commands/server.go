package commands

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/gzchatgpt/server"
)

func RegistryServer(app *cli.MultipleProgram) {
	app.Register("server", &cli.Command{
		Name:  "server",
		Usage: "start chatgpt api server",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Usage:   "server port",
				Aliases: []string{"p"},
				EnvVars: []string{"PORT"},
				Value:   8080,
			},
			&cli.StringFlag{
				Name:     "chatgpt-api-key",
				Usage:    "ChatGPT API Key",
				EnvVars:  []string{"CHATGPT_API_KEY"},
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "feishu-challenge",
				Usage: "Enable feishu challenge",
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			return server.Serve(&server.Config{
				Port:                  ctx.Int64("port"),
				ChatGPTAPIKey:         ctx.String("chatgpt-api-key"),
				EnableFeishuChallenge: ctx.Bool("feishu-challenge"),
			})
		},
	})
}
