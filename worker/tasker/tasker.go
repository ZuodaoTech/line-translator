package tasker

import (
	"context"
	"log/slog"
	"time"

	"github.com/zuodaotech/line-translator/common/assistant"
	"github.com/zuodaotech/line-translator/core"
)

type (
	Config struct {
	}

	Worker struct {
		cfg   Config
		tasks core.TaskStore

		taskz core.TaskService
		assi  *assistant.Assistant
	}
)

func New(
	cfg Config,
	tasks core.TaskStore,
	taskz core.TaskService,
	assi *assistant.Assistant,
) *Worker {

	return &Worker{
		cfg:   cfg,
		tasks: tasks,
		taskz: taskz,
		assi:  assi,
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
	case core.TaskActionTranslate:
		output, err = w.ProcessTaskActionTranslate(ctx, task)
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
