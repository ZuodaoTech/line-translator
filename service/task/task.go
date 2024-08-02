package task

import (
	"context"

	"github.com/zuodaotech/line-translator/common/uuid"
	"github.com/zuodaotech/line-translator/core"
	"github.com/zuodaotech/line-translator/store"
	"github.com/zuodaotech/line-translator/store/task"
)

func New(
	cfg Config,
	tasks core.TaskStore,
) *service {
	return &service{
		cfg:   cfg,
		tasks: tasks,
	}
}

type Config struct {
}

type service struct {
	cfg   Config
	tasks core.TaskStore
}

func (s *service) CreateTask(ctx context.Context, item *core.Task) (*core.Task, error) {
	traceID := uuid.New()
	item.TraceID = traceID
	tid, err := s.tasks.CreateTask(ctx, item)
	if err != nil {
		return nil, err
	}
	item.ID = tid
	item.TraceID = traceID
	return item, nil
}

func (s *service) BulkCreateTasks(ctx context.Context, items []*core.Task) ([]string, error) {
	traceIDs := make([]string, 0)
	if err := store.Transaction(func(tx *store.Handler) error {
		tasks := task.New(tx)
		for _, item := range items {
			traceID := uuid.New()
			item.TraceID = traceID
			_, err := tasks.CreateTask(ctx, item)
			if err != nil {
				return err
			}
			traceIDs = append(traceIDs, traceID)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return traceIDs, nil
}

func (s *service) GetTasksByStatus(ctx context.Context, status int, limit int) ([]*core.Task, error) {
	return s.tasks.GetTasksByStatus(ctx, status, limit)
}

func (s *service) GetTaskByTraceID(ctx context.Context, traceID string) (*core.Task, error) {
	tk, err := s.tasks.GetTaskByTraceID(ctx, traceID)
	if err != nil {
		return nil, err
	}
	pendingCount, err := s.tasks.CountPendingTasks(ctx)
	if err != nil {
		return nil, err
	}
	tk.PendingCount = pendingCount
	return tk, nil
}

func (s *service) UpdateTaskStatus(ctx context.Context, id uint64, status int) error {
	return s.tasks.UpdateTaskStatus(ctx, id, status)
}

func (s *service) UpdateTaskStatusWithResult(ctx context.Context, id uint64, status int, data core.JSONMap) error {
	result := &core.TaskResultWrapper{
		Code: status,
		Data: data,
	}
	return s.tasks.UpdateTaskStatusWithResult(ctx, id, status, result)
}

func (s *service) UpdateTaskStatusWithError(ctx context.Context, id uint64, status int, msg string) error {
	result := &core.TaskResultWrapper{
		Code:    status,
		Message: msg,
	}
	return s.tasks.UpdateTaskStatusWithResult(ctx, id, status, result)
}
