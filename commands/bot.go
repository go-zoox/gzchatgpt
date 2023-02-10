package commands

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/gzchatgpt/robot/feishu"
)

func RegistrFeishuBot(app *cli.MultipleProgram) {
	app.Register("feishu-bot", &cli.Command{
		Name:  "feishu-bot",
		Usage: "feishu bot",
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
			&cli.StringFlag{
				Name:     "app-id",
				Usage:    "Feishu App ID",
				EnvVars:  []string{"APP_ID"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "app-secret",
				Usage:    "Feishu App SECRET",
				EnvVars:  []string{"APP_SECRET"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "encrypt-key",
				Usage:   "enable encryption if you need",
				EnvVars: []string{"ENCRYPT_KEY"},
			},
			&cli.StringFlag{
				Name:    "verification-token",
				Usage:   "enable token verification if you need",
				EnvVars: []string{"VERIFICATION_TOKEN"},
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			return feishu.ServeFeishuBot(&feishu.FeishuBotConfig{
				Port:              ctx.Int64("port"),
				ChatGPTAPIKey:     ctx.String("chatgpt-api-key"),
				AppID:             ctx.String("app-id"),
				AppSecret:         ctx.String("app-secret"),
				EncryptKey:        ctx.String("encrypt-key"),
				VerificationToken: ctx.String("verification-token"),
			})
		},
	})
}
