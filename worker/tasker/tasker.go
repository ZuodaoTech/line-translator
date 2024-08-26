package tasker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/lyricat/goutils/langdetect"
	"github.com/lyricat/goutils/social/line"
	"github.com/lyricat/goutils/speech"
	"github.com/patrickmn/go-cache"
	"github.com/zuodaotech/line-translator/common/assistant"
	"github.com/zuodaotech/line-translator/core"
)

type (
	Config struct {
		LineChannelID     string
		LineChannelKey    string
		LineJWTPrivateKey string

		AzureEndpoint string
		AzureAPIKey   string
	}

	Worker struct {
		cfg   Config
		assi  *assistant.Assistant
		tasks core.TaskStore
		taskz core.TaskService

		detector   *langdetect.Detector
		speechCli  *speech.Client
		tokenCache *cache.Cache
	}
)

func New(
	cfg Config,
	assi *assistant.Assistant,
	tasks core.TaskStore,
	taskz core.TaskService,
) *Worker {

	detector := langdetect.New()
	speechCli := speech.New(speech.Config{
		AzureEndpoint: cfg.AzureEndpoint,
		AzureAPIKey:   cfg.AzureAPIKey,
	})

	c := cache.New(5*time.Minute, 10*time.Minute)

	return &Worker{
		cfg:   cfg,
		assi:  assi,
		tasks: tasks,
		taskz: taskz,

		detector:   detector,
		speechCli:  speechCli,
		tokenCache: c,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	dur := time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := w.run(ctx); err == nil {
				dur = 1 * time.Second
			} else {
				dur = 10 * time.Second
			}
		}
	}
}

func (w *Worker) run(ctx context.Context) error {
	tks, err := w.tasks.GetTasksByStatus(ctx, core.TaskStatusInit, 100)
	if err != nil {
		slog.Warn("[auxilia.tasker] failed to get init tasks", "error", err)
		return err
	}

	for _, tk := range tks {
		if err := w.ProcessTask(ctx, tk); err != nil {
			slog.Warn("[auxilia.tasker] failed to process task", "task.ID", tk.ID, "error", err)
		}
	}

	return nil
}

func (w *Worker) ProcessTask(ctx context.Context, task *core.Task) error {
	var err error
	var output core.JSONMap
	switch task.Action {
	case core.TaskActionQuoteAndTranslate:
		output, err = w.ProcessTaskActionQuoteAndTranslate(ctx, task)
	case core.TaskActionFetchAudioAndTranscript:
		output, err = w.ProcessTaskActionFetchAudioAndTranscript(ctx, task)
	case core.TaskActionSendMessage:
		output, err = w.ProcessTaskActionSendMessage(ctx, task)
	}

	if err != nil {
		slog.Warn("[tasker] failed to process task", "task.ID", task.ID, "error", err)
		if err := w.taskz.UpdateTaskStatusWithError(ctx, task.ID, core.TaskStatusError, err.Error()); err != nil {
			slog.Warn("[tasker] failed to update task status to error", "task.ID", task.ID, "error", err)
			return err
		}
		return err
	}

	if err := w.taskz.UpdateTaskStatusWithResult(ctx, task.ID, core.TaskStatusCompleted, output); err != nil {
		slog.Warn("[tasker] failed to update task status to complete", "task.ID", task.ID, "error", err)
		return err
	}
	return nil
}

func (w *Worker) ProcessTaskActionSendMessage(ctx context.Context, task *core.Task) (core.JSONMap, error) {
	text := task.Params.GetString("text")
	if text == "" {
		return nil, fmt.Errorf("source is empty")
	}

	groupID := task.Params.GetString("group_id")
	if groupID == "" {
		return nil, fmt.Errorf("group_id is empty")
	}

	replyToken := task.Params.GetString("reply_token")
	quoteToken := task.Params.GetString("quote_token")

	cli, err := w.GetLineClient(groupID)
	if err != nil {
		slog.Error("[worker.tasker] failed to generate token", "error", err)
		return nil, err
	}
	if _, err := cli.ReplyTextMessage(replyToken, quoteToken, text); err != nil {
		slog.Error("[worker.tasker] failed to send text reply", "error", err)
	} else {
		slog.Info("[worker.tasker] sent text reply")
	}

	return nil, nil
}

func (w *Worker) GetLineClient(groupId string) (*line.Client, error) {
	var cli *line.Client
	var err error
	val, found := w.tokenCache.Get(groupId)

	if found {
		slog.Info("[worker.tasker] found token in cache", "groupID", groupId)
		item, ok := val.(*TokenCacheItem)
		if ok && item.ExpireAt.After(time.Now()) && item.AccessToken != "" {
			slog.Info("[worker.tasker] token is valid", "groupID", groupId, "expireAt", item.ExpireAt)
			fmt.Printf("item.AccessToken: %v\n", item.AccessToken)
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
			slog.Error("[worker.tasker] failed to create line client", "error", err)
			return nil, err
		}
		token, expired, err := cli.GenerateToken()
		if err != nil {
			slog.Error("[worker.tasker] failed to generate token", "error", err)
			return nil, err
		}
		expiredDur := time.Until(*expired)
		slog.Info("[worker.tasker] generated token", "groupID", groupId, "expiredDur", expiredDur)
		w.tokenCache.Set(groupId, &TokenCacheItem{
			AccessToken: token,
			ExpireAt:    *expired,
		}, expiredDur)
	}

	return cli, nil
}
