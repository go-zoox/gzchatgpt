package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/go-zoox/zoox"
	"github.com/go-zoox/zoox/defaults"
)

type Config struct {
	Port                  int64
	ChatGPTAPIKey         string
	EnableFeishuChallenge bool
}

func Serve(cfg *Config) error {
	app := defaults.Application()

	client := gpt3.NewClient(cfg.ChatGPTAPIKey)

	app.Post("/ask", func(ctx *zoox.Context) {
		type DTO struct {
			Question string `json:"question"`
		}
		var dto DTO
		if err := ctx.BindJSON(&dto); err != nil {
			ctx.Error(http.StatusBadRequest, "invalid request")
			return
		}
		c := context.Background()

		ctx.Logger.Infof("问题：%s", dto.Question)
		response, err := client.CompletionWithEngine(c, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
			Prompt: []string{
				dto.Question,
			},
			MaxTokens:   gpt3.IntPtr(3000),
			Temperature: gpt3.Float32Ptr(0),
		})
		if err != nil {
			ctx.Logger.Errorf("failed to request answer: %v", err)
			ctx.Error(http.StatusBadRequest, "failed to answer")
			return
		}

		ctx.Logger.Infof("回答：%s", response.Choices[0].Text)

		ctx.JSON(200, zoox.H{
			"answer": response.Choices[0].Text,
		})
	})

	// if cfg.EnableFeishuChallenge {
	// 	app.Logger.Infof("enable feishu challenge")
	// 	middleware := larkgin.NewLarkMiddleware()
	// 	app.Post("/robot/feishu", func(ctx *zoox.Context) {
	// 		ctx.Logger.Infof("feishu challenge request")

	// 		type Message struct {
	// 			Challenge string `json:"challenge"`
	// 			//
	// 			Schema string `json:"schema"`
	// 			Header struct {
	// 				EventID    string `json:"event_id"`
	// 				Token      string `json:"token"`
	// 				CreateTime string `json:"create_time"`
	// 				EventType  string `json:"event_type"`
	// 				TenantKey  string `json:"tenant_key"`
	// 				AppID      string `json:"app_id"`
	// 			} `json:"header"`
	// 			Event struct {
	// 				AppID    string `json:"app_id"`
	// 				ChatID   string `json:"chat_id"`
	// 				Operator struct {
	// 					OpenID string `json:"open_id"`
	// 					UserID string `json:"user_id"`
	// 				} `json:"operator"`
	// 				TenantKey string `json:"tenant_key"`
	// 				Type      string `json:"type"`
	// 				User      struct {
	// 					Name   string `json:"name"`
	// 					OpenID string `json:"open_id"`
	// 					UserID string `json:"user_id"`
	// 				} `json:"user"`
	// 			} `json:"event"`
	// 		}
	// 		var message Message
	// 		if err := ctx.BindJSON(&message); err != nil {
	// 			ctx.Logger.Errorf("bad request")
	// 			ctx.String(400, "bad request")
	// 			return
	// 		}

	// 		if message.Challenge != "" {
	// 			ctx.JSON(200, message)
	// 			return
	// 		}

	// 		if evt, ok := middleware.GetEvent(ctx); ok { // => GetEvent instead of GetMessage
	// 			if evt.Header.EventType == lark.EventTypeMessageReceived {
	// 				if msg, err := evt.GetMessageReceived(); err == nil {
	// 					fmt.Println(msg.Message.Content)
	// 				}
	// 				// you may have to parse other events
	// 			}
	// 		}
	// 	})
	// }

	app.Get("/", func(ctx *zoox.Context) {
		ctx.String(200, "hello world")
	})

	return app.Run(fmt.Sprintf(":%d", cfg.Port))
}
