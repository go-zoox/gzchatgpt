package feishu

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/go-lark/lark"
	"github.com/go-zoox/core-utils/regexp"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/retry"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type FeishuBotConfig struct {
	ChatGPTAPIKey     string
	AppID             string
	AppSecret         string
	EncryptKey        string
	VerificationToken string
}

func ServeFeishuBot(cfg *FeishuBotConfig) error {
	client := gpt3.NewClient(cfg.ChatGPTAPIKey)
	bot := lark.NewChatBot(cfg.AppID, cfg.AppSecret)
	_, _ = bot.GetTenantAccessTokenInternal(true)

	// 注册消息处理器
	handler := dispatcher.NewEventDispatcher(cfg.VerificationToken, cfg.EncryptKey).
		OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			// 处理消息 event，这里简单打印消息的内容
			fmt.Println(larkcore.Prettify(event))
			fmt.Println("OnP2MessageReceiveV1", event.RequestId())

			message := event.Event.Message.Content
			if message != nil {
				type Content struct {
					Text string `json:"text"`
				}
				var content Content
				if err := json.Unmarshal([]byte(*message), &content); err != nil {
					return err
				}

				textMessage := content.Text
				if textMessage != "" {
					fmt.Println("textMessage:", textMessage)
					if ok := regexp.Match("^@_user_1", textMessage); ok {
						question := textMessage[len("@_user_1 "):]
						fmt.Println("question:", question)
						for _, metion := range event.Event.Message.Mentions {
							if *metion.Key == "@_user_1" && *metion.Name == "ChatGPT" {
								go func() {
									logger.Infof("问题：%s", question)
									var err error
									var response *gpt3.CompletionResponse
									err = retry.Retry(func() error {
										response, err = client.CompletionWithEngine(context.Background(), gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
											Prompt: []string{
												question,
											},
											MaxTokens:   gpt3.IntPtr(3000),
											Temperature: gpt3.Float32Ptr(0),
										})
										if err != nil {
											logger.Errorf("failed to request answer: %v", err)
											return fmt.Errorf("failed to request answer: %v", err)
										}

										return nil
									}, 5, 3*time.Second)
									if err != nil {
										logger.Errorf("failed to get answer: %v", err)
										return
									}

									answer := strings.TrimSpace(response.Choices[0].Text)
									logger.Infof("回答：%s", answer)

									//
									msg := lark.NewMsgBuffer(lark.MsgPost)
									postContent := lark.NewPostBuilder().
										// Title("asdaads").
										TextTag(answer, 1, true).
										Render()
									om := msg.BindOpenChatID(*event.Event.Message.ChatId).Post(postContent).Build()
									resp, err := bot.PostMessage(om)
									if err != nil {
										logger.Errorf("failed to post message: %v", err)
										return
									}

									logger.Infof("robot response: %v", resp)
								}()

								return nil
							}
						}
					}
				}
			}

			return nil
		}).
		OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
			// 处理消息 event，这里简单打印消息的内容
			fmt.Println(larkcore.Prettify(event))
			fmt.Println("OnP2MessageReadV1", event.RequestId())
			return nil
		})

	// 注册 http 路由
	// http.HandleFunc("/webhook/event", httpserverext.NewEventHandlerFunc(handler, larkevent.WithLogLevel(larkcore.LogLevelDebug)))
	http.HandleFunc("/bot/feishu", httpserverext.NewEventHandlerFunc(handler, larkevent.WithLogLevel(larkcore.LogLevelDebug)))

	// 启动 http 服务
	return http.ListenAndServe(":8080", nil)
}
