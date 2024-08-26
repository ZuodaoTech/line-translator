package tasker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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

	strategy := task.Params.GetString("strategy")
	if strategy == "" {
		return nil, fmt.Errorf("strategy is empty")
	}

	srcLang, found := w.detector.Detect(text)
	if !found {
		srcLang = "en"
	}

	dstLang := ""
	if strategy == core.ConversationStrategyZhEn {
		if srcLang == "zh" {
			dstLang = "en"
		} else {
			dstLang = "zh"
		}
	} else if strategy == core.ConversationStrategyZhJa {
		if srcLang == "zh" {
			dstLang = "ja"
		} else {
			dstLang = "zh"
		}
	}

	var err error
	result := text
	if srcLang != dstLang {
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

	return jsonMap, nil
}
