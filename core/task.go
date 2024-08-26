package core

import (
	"context"
	"time"
)

const (
	TaskStatusInit      = iota
	TaskStatusScheduled // 1
	TaskStatusRunning   // 2
	TaskStatusCompleted // 3
	TaskStatusError     // 4
)

const (
	TaskActionSendMessage             = "send_message"
	TaskActionQuoteAndTranslate       = "quote_and_translate"
	TaskActionFetchAudioAndTranscript = "fetch_audio_and_transcript"
)

type (
	TaskResultWrapper struct {
		Code    int     `json:"code"`
		Data    JSONMap `json:"data"`
		Message string  `json:"message"`
	}

	Task struct {
		ID          uint64     `json:"id"`
		UserID      uint64     `json:"user_id"`
		Action      string     `json:"action"`
		Params      JSONMap    `gorm:"jsonb" json:"params"`
		Result      JSONMap    `gorm:"jsonb" json:"result"`
		Status      int        `json:"status"`
		TraceID     string     `json:"trace_id"`
		ScheduledAt *time.Time `json:"scheduled_at"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   time.Time  `json:"updated_at"`

		PendingCount int `gorm:"-" json:"pending_count"`
	}

	TaskStore interface {
		// INSERT INTO @@table (
		//  user_id, action, params, status,
		//  trace_id,
		//  {{if item.ScheduledAt != nil}} scheduled_at, {{end}}
		//  created_at, updated_at
		// ) VALUES (
		//  @item.UserID, @item.Action, @item.Params,
		//  {{if item.ScheduledAt != nil}} 1, {{else}} 0, {{end}}
		//  @item.TraceID,
		//  {{if item.ScheduledAt != nil}} @item.ScheduledAt, {{end}}
		//  NOW(), NOW()
		// )
		// RETURNING id;
		CreateTask(ctx context.Context, item *Task) (uint64, error)

		// SELECT * FROM @@table
		// WHERE
		// 	status = @status
		// ORDER BY created_at DESC
		// LIMIT @limit
		GetTasksByStatus(ctx context.Context, status int, limit int) ([]*Task, error)

		// SELECT * FROM @@table
		// WHERE
		// 	trace_id = @traceID
		// LIMIT 1;
		GetTaskByTraceID(ctx context.Context, traceID string) (*Task, error)

		// SELECT COUNT(*) FROM @@table
		// WHERE
		// 	status = 0 OR status = 2
		// LIMIT 10;
		CountPendingTasks(ctx context.Context) (int, error)

		// UPDATE @@table
		// SET
		//  result = @result,
		//  status = @status,
		//  updated_at = NOW()
		// WHERE id = @id;
		UpdateTaskStatusWithResult(ctx context.Context, id uint64, status int, result any) error

		// UPDATE @@table
		// SET
		//  status = @status,
		//  updated_at = NOW()
		// WHERE id = @id;
		UpdateTaskStatus(ctx context.Context, id uint64, status int) error
	}

	TaskService interface {
		CreateTask(ctx context.Context, item *Task) (*Task, error)
		GetTasksByStatus(ctx context.Context, status int, limit int) ([]*Task, error)
		GetTaskByTraceID(ctx context.Context, traceID string) (*Task, error)
		UpdateTaskStatus(ctx context.Context, id uint64, status int) error
		UpdateTaskStatusWithResult(ctx context.Context, id uint64, status int, data JSONMap) error
		UpdateTaskStatusWithError(ctx context.Context, id uint64, status int, message string) error

		BulkCreateTasks(ctx context.Context, items []*Task) ([]string, error)
	}
)

func (t *Task) GetResult() JSONMap {
	return t.Result
}
