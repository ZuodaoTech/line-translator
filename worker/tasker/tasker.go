package tasker

import (
	"context"
	"log/slog"
	"time"

	"github.com/lyricat/goutils/langdetect"
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
