package tasker

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/zuodaotech/line-translator/core"
)

func (w *Worker) ProcessTaskActionTranslate(ctx context.Context, task *core.Task) (core.JSONMap, error) {
	text := task.Params.GetString("text")
	dstLang := task.Params.GetString("dst_lang")
	if text == "" {
		return nil, fmt.Errorf("source is empty")
	}
	if dstLang == "" {
		dstLang = "English"
	}

	// if the text is too long, split it into multiple parts
	// each part is less than 30 lines
	// then translate each part
	// finally, combine all parts together
	lines := strings.Split(text, "\n")
	translateParts := make([]string, 0)
	if len(lines) > 30 {
		for i := 0; i < len(lines); i += 30 {
			end := i + 30
			if end > len(lines) {
				end = len(lines)
			}
			part := strings.Join(lines[i:end], "\n")
			translateParts = append(translateParts, part)
		}
	} else {
		translateParts = append(translateParts, text)
	}

	translated := ""
	for _, part := range translateParts {
		result, err := w.assi.Translate(ctx, part)
		if err != nil {
			slog.Error("translate failed.", "error", err)
			break
		}
		translated += result + "\n"
	}

	jsonMap := core.NewJSONMap()
	jsonMap.SetValue("translated", translated)
	jsonMap.SetValue("improved", translated)

	return jsonMap, nil
}
