package tasker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/zuodaotech/line-translator/core"
)

func (w *Worker) ProcessTaskActionQuoteAndTranslate(ctx context.Context, task *core.Task) (core.JSONMap, error) {
	text := task.Params.GetString("text")
	if text == "" {
		return nil, fmt.Errorf("source is empty")
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
		_, _, err = w.lineCli.GenerateToken()
		if err != nil {
			slog.Error("[worker.tasker.translator] failed to generate token", "error", err)
			return nil, err
		}
		if _, err := w.lineCli.ReplyTextMessage(replyToken, quoteToken, result); err != nil {
			slog.Error("[worker.tasker.translator] failed to send text reply", "error", err)
		} else {
			slog.Info("[worker.tasker.translator] sent text reply")
		}
	}

	return jsonMap, nil
}
