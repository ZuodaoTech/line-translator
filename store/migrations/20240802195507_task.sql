-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS "tasks" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT NOT NULL,
  "action" VARCHAR(255) NOT NULL,
  "params" JSONB NOT NULL DEFAULT '{}',
  "result" JSONB NOT NULL DEFAULT '{}',
  "status" INT NOT NULL DEFAULT 0,
  "trace_id" VARCHAR(36) NOT NULL,
  "scheduled_at" timestamptz,

  "created_at" timestamptz,
  "updated_at" timestamptz
);
CREATE INDEX IF NOT EXISTS "idx_tasks_user" ON "tasks" ("user_id");
CREATE INDEX IF NOT EXISTS "idx_tasks_status" ON "tasks" ("status");
CREATE INDEX IF NOT EXISTS "idx_tasks_trace_id" ON "tasks" ("trace_id");


CREATE TABLE IF NOT EXISTS "conversations" (
  "id" BIGSERIAL PRIMARY KEY,
  "channel" VARCHAR(32) NOT NULL,
  "channel_id" VARCHAR(255) NOT NULL,
  "strategy" VARCHAR(32) NOT NULL DEFAULT '',

  "created_at" timestamptz,
  "updated_at" timestamptz
);
CREATE UNIQUE INDEX IF NOT EXISTS "uidx_convs_channel" ON "conversations" ("channel", "channel_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS "tasks";
-- +goose StatementEnd
