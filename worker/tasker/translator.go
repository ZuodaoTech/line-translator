package tasker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/zuodaotech/line-translator/common/line"
	"github.com/zuodaotech/line-translator/core"
)

type (
	TokenCacheItem struct {
		AccessToken string
		ExpireAt    time.Time
	}
)

func (w *Worker) ProcessTaskActionQuoteAndTranslate(ctx context.Context, task *core.Task) (core.JSONMap, error) {
	text := task.Params.GetString("text")
	if text == "" {
		return nil, fmt.Errorf("source is empty")
	}

	groupID := task.Params.GetString("group_id")
	if groupID == "" {
		return nil, fmt.Errorf("group_id is empty")
	}

	srcLang, found := w.detector.Detect(text)
	if !found {
		srcLang = "en"
	}

	dstLang := "en"
	if srcLang == "zh" {
		dstLang = "ja"
	} else if srcLang == "ja" {
		dstLang = "zh"
	}

	fmt.Printf("srcLang: %v\n", srcLang)
	fmt.Printf("dstLang: %v\n", dstLang)

	var err error
	result := text
	if srcLang != "en" && srcLang != dstLang {
		result, err = w.assi.Translate(ctx, text, srcLang, dstLang)
		if err != nil {
			slog.Error("translate failed.", "error", err)
			return nil, err
		}
	}

	jsonMap := core.NewJSONMap()
	jsonMap.SetValue("translated", result)

	// send the message back
	replyToken := task.Params.GetString("reply_token")
	quoteToken := task.Params.GetString("quote_token")

	if replyToken != "" {
		cli, err := w.GetLineClient(groupID)
		if err != nil {
			slog.Error("[worker.tasker.translator] failed to generate token", "error", err)
			return nil, err
		}
		if _, err := cli.ReplyTextMessage(replyToken, quoteToken, result); err != nil {
			slog.Error("[worker.tasker.translator] failed to send text reply", "error", err)
		} else {
			slog.Info("[worker.tasker.translator] sent text reply")
		}
	}

	return jsonMap, nil
}

func (w *Worker) GetLineClient(groupId string) (*line.Client, error) {
	var cli *line.Client
	var err error
	val, found := w.tokenCache.Get(groupId)

	if found {
		slog.Info("[worker.tasker.translator] found token in cache", "groupID", groupId)
		item, ok := val.(*TokenCacheItem)
		if ok && item.ExpireAt.After(time.Now()) && item.AccessToken != "" {
			slog.Info("[worker.tasker.translator] token is valid", "groupID", groupId, "expireAt", item.ExpireAt)
			cli, err = line.NewFromAccessToken(item.AccessToken)
			if err != nil {
				return nil, err
			}
		}
	}

	if cli == nil {
		cli, err = line.New(line.Config{
			ChannelID:  w.cfg.LineChannelID,
			ChannelKey: w.cfg.LineChannelKey,
			PrivateKey: w.cfg.LineJWTPrivateKey,
		})
		if err != nil {
			slog.Error("[worker.tasker.translator] failed to create line client", "error", err)
			return nil, err
		}
		token, expired, err := cli.GenerateToken()
		if err != nil {
			slog.Error("[worker.tasker.translator] failed to generate token", "error", err)
			return nil, err
		}
		expiredDur := time.Until(*expired)
		slog.Info("[worker.tasker.translator] generated token", "groupID", groupId, "expiredDur", expiredDur)
		w.tokenCache.Set(groupId, &TokenCacheItem{
			AccessToken: token,
			ExpireAt:    *expired,
		}, expiredDur)
	}

	return cli, nil
}
