package line

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"github.com/zuodaotech/line-translator/config"
	"github.com/zuodaotech/line-translator/core"
	"github.com/zuodaotech/line-translator/handler/render"
	"gorm.io/gorm"
)

const Usage = "欢迎使用本机器人，请发送消息设置翻译模式：\n🐼中文 ↔️ 🗽英文，发送「/中英」\n🐼中文 ↔️ 🌸日文，发送「/中日」\n"

func HandleWebhook(syscfg *config.Config, conversations core.ConversationStore, taskz core.TaskService) http.HandlerFunc {
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
							slog.Info("[handler.line] group join event", "group.id", src.GroupId)
							// create a conversation for the group
							if _, err := conversations.CreateConversation(ctx, &core.Conversation{
								Channel:   "line",
								ChannelID: src.GroupId,
							}); err != nil {
								slog.Error("[handler.line] failed to create conversation", "error", err)
								render.Error(w, http.StatusInternalServerError, err)
								return
							}

							quickRespond(ctx, Usage, src.GroupId, e.ReplyToken, taskz)
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

									conv, err := conversations.GetConversationByChannel(ctx, "line", src.GroupId)
									if err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
										slog.Error("[handler.line] failed to get conversation", "error", err)
										render.Error(w, http.StatusInternalServerError, err)
										return
									}
									if conv.ID == 0 {
										slog.Warn("[handler.line] conversation not found, let's create one", "group_id", src.GroupId)
										// create a conversation for the group
										conv.Channel = "line"
										conv.ChannelID = src.GroupId
										convID, err := conversations.CreateConversation(ctx, conv)
										if err != nil {
											slog.Error("[handler.line] failed to create conversation", "error", err)
											render.Error(w, http.StatusInternalServerError, err)
											return
										}
										conv.ID = convID
									}

									text := strings.TrimSpace(msg.Text)
									if text == "" {
										continue
									}

									// handle command
									if msg.Text[0] == '/' || strings.HasPrefix(msg.Text, "／") {
										cmd := strings.TrimLeft(strings.TrimLeft(msg.Text, "/"), "／")
										// update conversation strategy
										strategy := ""
										hint := ""
										if cmd == "中英" || cmd == "英中" || cmd == "中日" || cmd == "日中" {
											if cmd == "中英" || cmd == "英中" {
												strategy = core.ConversationStrategyZhEn
												hint = "翻译模式：🐼中文 ↔️ 🗽英文"
											} else if cmd == "中日" || cmd == "日中" {
												strategy = core.ConversationStrategyZhJa
												hint = "翻译模式：🐼中文 ↔️ 🌸日文"
											}
											conv.Strategy = strategy
											if err := conversations.UpdateConversationStrategy(ctx, conv.ID, strategy); err != nil {
												slog.Error("[handler.line] failed to update conversation", "error", err)
												render.Error(w, http.StatusInternalServerError, err)
												return
											}
											quickRespond(ctx, hint, src.GroupId, e.ReplyToken, taskz)
										} else {
											quickRespond(ctx, Usage, src.GroupId, e.ReplyToken, taskz)
										}
										continue
									}

									// if no strategy, send usage
									if conv.Strategy == "" {
										if err := quickRespond(ctx, Usage, src.GroupId, e.ReplyToken, taskz); err != nil {
											slog.Error("[handler.line] failed to quick respond", "error", err)
											render.Error(w, http.StatusInternalServerError, err)
											return
										}
										continue
									}

									// create a task to translate the message
									data := &core.Task{
										UserID: 0,
										Action: core.TaskActionQuoteAndTranslate,
										Params: map[string]any{
											"group_id":    src.GroupId,
											"reply_token": e.ReplyToken,
											"quote_token": msg.QuoteToken,
											"text":        msg.Text,
											"strategy":    conv.Strategy,
										},
										Status: core.TaskStatusInit,
									}

									if _, err := taskz.CreateTask(ctx, data); err != nil {
										slog.Error("[handler.line] failed to create task", "error", err)
										render.Error(w, http.StatusInternalServerError, err)
										return
									}
								}
							case webhook.AudioMessageContent:
								{
									if err := quickRespond(ctx, "暂不支持语音消息", src.GroupId, e.ReplyToken, taskz); err != nil {
										slog.Error("[handler.line] failed to quick respond", "error", err)
										render.Error(w, http.StatusInternalServerError, err)
										return
									}

									// @TODO: fetch audio and transcript, disabled.
									//
									// slog.Info("[handler.line] audio message", "audio", msg)
									// fmt.Printf("[handler.line] msg: %v\n", msg)
									// data := &core.Task{
									// 	UserID: 0,
									// 	Action: core.TaskActionFetchAudioAndTranscript,
									// 	Params: map[string]interface{}{
									// 		"group_id":    src.GroupId,
									// 		"reply_token": e.ReplyToken,
									// 		"message_id":  msg.Id,
									// 	},
									// 	Status: core.TaskStatusInit,
									// }

									// if _, err := taskz.CreateTask(ctx, data); err != nil {
									// 	slog.Error("[handler.line] failed to create task", "error", err)
									// 	render.Error(w, http.StatusInternalServerError, err)
									// 	return
									// }
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

func quickRespond(ctx context.Context, msg string, groupID, replyToken string, taskz core.TaskService) error {
	data := &core.Task{
		UserID: 0,
		Action: core.TaskActionSendMessage,
		Params: map[string]interface{}{
			"group_id":    groupID,
			"reply_token": replyToken,
			"text":        msg,
		},
		Status: core.TaskStatusInit,
	}
	if _, err := taskz.CreateTask(ctx, data); err != nil {
		slog.Error("[handler.line] failed to create task", "error", err)
		return err
	}
	return nil
}
