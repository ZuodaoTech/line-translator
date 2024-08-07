package line

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"github.com/zuodaotech/line-translator/config"
	"github.com/zuodaotech/line-translator/core"
	"github.com/zuodaotech/line-translator/handler/render"
)

func HandleWebhook(syscfg *config.Config, taskz core.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// parse webhook
		cb, err := webhook.ParseRequest(syscfg.Line.ChannelSecret, r)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				slog.Error("[handler.line] invalid signature", "error", err)
				render.Error(w, http.StatusBadRequest, err)
				return
			}
			slog.Error("[handler.line] failed to parse request", "error", err)
			render.Error(w, http.StatusInternalServerError, err)
			return
		}

		// handle events
		for _, event := range cb.Events {
			// fmt.Printf("webhook called %+v...\n\n", event)
			switch e := event.(type) {
			case webhook.JoinEvent:
				{
					fmt.Printf("Join: %+v\n\n", e)
					switch src := e.Source.(type) {
					case webhook.GroupSource:
						{
							fmt.Printf("Group: %s\n\n", src.GroupId)
							// cli, err := line.New(line.Config{
							// 	ChannelID:  syscfg.Line.ChannelID,
							// 	ChannelKey: syscfg.Line.ChannelKey,
							// 	PrivateKey: syscfg.Line.JWTPrivateKey,
							// })
							// if err != nil {
							// 	slog.Error("[handler.line] failed to create line client", "error", err)
							// 	render.Error(w, http.StatusInternalServerError, err)
							// 	return
							// }
							// _, _, err = cli.GenerateToken()
							// if err != nil {
							// 	slog.Error("[handler.line] failed to generate token", "error", err)
							// 	return
							// }

							// if _, err := cli.ReplyTextMessage(e.ReplyToken, "", fmt.Sprintf("Hello, here is your group ID:\n%s", src.GroupId)); err != nil {
							// 	slog.Error("[handler.line] failed to send text reply", "error", err, "groupID", src.GroupId)
							// } else {
							// 	slog.Info("[handler.line] sent text reply", "groupID", src.GroupId)
							// }
						}
					case webhook.RoomSource:
						{
							slog.Info("[handler.line] room join event", "roomID", src.RoomId)
						}
					}
				}
			case webhook.MessageEvent:
				{
					switch src := e.Source.(type) {
					case webhook.GroupSource:
						{
							slog.Info("[handler.line] message event", "event.message", e.Message)
							// send a reply message
							switch msg := e.Message.(type) {
							case webhook.TextMessageContent:
								{
									slog.Info("[handler.line] text message", "text", msg.Text)

									// create a task
									data := &core.Task{
										UserID: 0,
										Action: core.TaskActionQuoteAndTranslate,
										Params: map[string]interface{}{
											"group_id":    src.GroupId,
											"reply_token": e.ReplyToken,
											"quote_token": msg.QuoteToken,
											"text":        msg.Text,
										},
										Status: core.TaskStatusInit,
									}

									if _, err := taskz.CreateTask(ctx, data); err != nil {
										slog.Error("[handler.line] failed to create task", "error", err)
										render.Error(w, http.StatusInternalServerError, err)
										return
									}
								}
							}
						}
					}
				}
			}
		}
		render.JSON(w, nil)
	}
}
