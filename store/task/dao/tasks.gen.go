// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/zuodaotech/line-translator/core"
)

func newTask(db *gorm.DB, opts ...gen.DOOption) task {
	_task := task{}

	_task.taskDo.UseDB(db, opts...)
	_task.taskDo.UseModel(&core.Task{})

	tableName := _task.taskDo.TableName()
	_task.ALL = field.NewAsterisk(tableName)
	_task.ID = field.NewUint64(tableName, "id")
	_task.UserID = field.NewUint64(tableName, "user_id")
	_task.Action = field.NewString(tableName, "action")
	_task.Params = field.NewField(tableName, "params")
	_task.Result = field.NewField(tableName, "result")
	_task.Status = field.NewInt(tableName, "status")
	_task.TraceID = field.NewString(tableName, "trace_id")
	_task.ScheduledAt = field.NewTime(tableName, "scheduled_at")
	_task.CreatedAt = field.NewTime(tableName, "created_at")
	_task.UpdatedAt = field.NewTime(tableName, "updated_at")

	_task.fillFieldMap()

	return _task
}

type task struct {
	taskDo

	ALL         field.Asterisk
	ID          field.Uint64
	UserID      field.Uint64
	Action      field.String
	Params      field.Field
	Result      field.Field
	Status      field.Int
	TraceID     field.String
	ScheduledAt field.Time
	CreatedAt   field.Time
	UpdatedAt   field.Time

	fieldMap map[string]field.Expr
}

func (t task) Table(newTableName string) *task {
	t.taskDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t task) As(alias string) *task {
	t.taskDo.DO = *(t.taskDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *task) updateTableName(table string) *task {
	t.ALL = field.NewAsterisk(table)
	t.ID = field.NewUint64(table, "id")
	t.UserID = field.NewUint64(table, "user_id")
	t.Action = field.NewString(table, "action")
	t.Params = field.NewField(table, "params")
	t.Result = field.NewField(table, "result")
	t.Status = field.NewInt(table, "status")
	t.TraceID = field.NewString(table, "trace_id")
	t.ScheduledAt = field.NewTime(table, "scheduled_at")
	t.CreatedAt = field.NewTime(table, "created_at")
	t.UpdatedAt = field.NewTime(table, "updated_at")

	t.fillFieldMap()

	return t
}

func (t *task) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *task) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 10)
	t.fieldMap["id"] = t.ID
	t.fieldMap["user_id"] = t.UserID
	t.fieldMap["action"] = t.Action
	t.fieldMap["params"] = t.Params
	t.fieldMap["result"] = t.Result
	t.fieldMap["status"] = t.Status
	t.fieldMap["trace_id"] = t.TraceID
	t.fieldMap["scheduled_at"] = t.ScheduledAt
	t.fieldMap["created_at"] = t.CreatedAt
	t.fieldMap["updated_at"] = t.UpdatedAt
}

func (t task) clone(db *gorm.DB) task {
	t.taskDo.ReplaceConnPool(db.Statement.ConnPool)
	return t
}

func (t task) replaceDB(db *gorm.DB) task {
	t.taskDo.ReplaceDB(db)
	return t
}

type taskDo struct{ gen.DO }

type ITaskDo interface {
	gen.SubQuery
	Debug() ITaskDo
	WithContext(ctx context.Context) ITaskDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ITaskDo
	WriteDB() ITaskDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ITaskDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ITaskDo
	Not(conds ...gen.Condition) ITaskDo
	Or(conds ...gen.Condition) ITaskDo
	Select(conds ...field.Expr) ITaskDo
	Where(conds ...gen.Condition) ITaskDo
	Order(conds ...field.Expr) ITaskDo
	Distinct(cols ...field.Expr) ITaskDo
	Omit(cols ...field.Expr) ITaskDo
	Join(table schema.Tabler, on ...field.Expr) ITaskDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ITaskDo
	RightJoin(table schema.Tabler, on ...field.Expr) ITaskDo
	Group(cols ...field.Expr) ITaskDo
	Having(conds ...gen.Condition) ITaskDo
	Limit(limit int) ITaskDo
	Offset(offset int) ITaskDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ITaskDo
	Unscoped() ITaskDo
	Create(values ...*core.Task) error
	CreateInBatches(values []*core.Task, batchSize int) error
	Save(values ...*core.Task) error
	First() (*core.Task, error)
	Take() (*core.Task, error)
	Last() (*core.Task, error)
	Find() ([]*core.Task, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*core.Task, err error)
	FindInBatches(result *[]*core.Task, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*core.Task) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ITaskDo
	Assign(attrs ...field.AssignExpr) ITaskDo
	Joins(fields ...field.RelationField) ITaskDo
	Preload(fields ...field.RelationField) ITaskDo
	FirstOrInit() (*core.Task, error)
	FirstOrCreate() (*core.Task, error)
	FindByPage(offset int, limit int) (result []*core.Task, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ITaskDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	CreateTask(ctx context.Context, item *core.Task) (result uint64, err error)
	GetTasksByStatus(ctx context.Context, status int, limit int) (result []*core.Task, err error)
	GetTaskByTraceID(ctx context.Context, traceID string) (result *core.Task, err error)
	CountPendingTasks(ctx context.Context) (result int, err error)
	UpdateTaskStatusWithResult(ctx context.Context, id uint64, status int, result any) (err error)
	UpdateTaskStatus(ctx context.Context, id uint64, status int) (err error)
}

// INSERT INTO @@table (
// user_id, action, params, status,
// trace_id,
// {{if item.ScheduledAt != nil}} scheduled_at, {{end}}
// created_at, updated_at
// ) VALUES (
// @item.UserID, @item.Action, @item.Params,
// {{if item.ScheduledAt != nil}} 1, {{else}} 0, {{end}}
// @item.TraceID,
// {{if item.ScheduledAt != nil}} @item.ScheduledAt, {{end}}
// NOW(), NOW()
// )
// RETURNING id;
func (t taskDo) CreateTask(ctx context.Context, item *core.Task) (result uint64, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("INSERT INTO tasks ( user_id, action, params, status, trace_id, ")
	if item.ScheduledAt != nil {
		generateSQL.WriteString("scheduled_at, ")
	}
	params = append(params, item.UserID)
	params = append(params, item.Action)
	params = append(params, item.Params)
	generateSQL.WriteString("created_at, updated_at ) VALUES ( ?, ?, ?, ")
	if item.ScheduledAt != nil {
		generateSQL.WriteString("1, ")
	} else {
		generateSQL.WriteString("0, ")
	}
	params = append(params, item.TraceID)
	generateSQL.WriteString("?, ")
	if item.ScheduledAt != nil {
		params = append(params, item.ScheduledAt)
		generateSQL.WriteString("?, ")
	}
	generateSQL.WriteString("NOW(), NOW() ) RETURNING id; ")

	var executeSQL *gorm.DB
	executeSQL = t.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// SELECT * FROM @@table
// WHERE
//
//	status = @status
//
// ORDER BY created_at DESC
// LIMIT @limit
func (t taskDo) GetTasksByStatus(ctx context.Context, status int, limit int) (result []*core.Task, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, status)
	params = append(params, limit)
	generateSQL.WriteString("SELECT * FROM tasks WHERE status = ? ORDER BY created_at DESC LIMIT ? ")

	var executeSQL *gorm.DB
	executeSQL = t.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// SELECT * FROM @@table
// WHERE
//
//	trace_id = @traceID
//
// LIMIT 1;
func (t taskDo) GetTaskByTraceID(ctx context.Context, traceID string) (result *core.Task, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, traceID)
	generateSQL.WriteString("SELECT * FROM tasks WHERE trace_id = ? LIMIT 1; ")

	var executeSQL *gorm.DB
	executeSQL = t.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// SELECT COUNT(*) FROM @@table
// WHERE
//
//	status = 0 OR status = 2
//
// LIMIT 10;
func (t taskDo) CountPendingTasks(ctx context.Context) (result int, err error) {
	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT COUNT(*) FROM tasks WHERE status = 0 OR status = 2 LIMIT 10; ")

	var executeSQL *gorm.DB
	executeSQL = t.UnderlyingDB().Raw(generateSQL.String()).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// UPDATE @@table
// SET
// result = @result,
// status = @status,
// updated_at = NOW()
// WHERE id = @id;
func (t taskDo) UpdateTaskStatusWithResult(ctx context.Context, id uint64, status int, result any) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, result)
	params = append(params, status)
	params = append(params, id)
	generateSQL.WriteString("UPDATE tasks SET result = ?, status = ?, updated_at = NOW() WHERE id = ?; ")

	var executeSQL *gorm.DB
	executeSQL = t.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

// UPDATE @@table
// SET
// status = @status,
// updated_at = NOW()
// WHERE id = @id;
func (t taskDo) UpdateTaskStatus(ctx context.Context, id uint64, status int) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, status)
	params = append(params, id)
	generateSQL.WriteString("UPDATE tasks SET status = ?, updated_at = NOW() WHERE id = ?; ")

	var executeSQL *gorm.DB
	executeSQL = t.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (t taskDo) Debug() ITaskDo {
	return t.withDO(t.DO.Debug())
}

func (t taskDo) WithContext(ctx context.Context) ITaskDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t taskDo) ReadDB() ITaskDo {
	return t.Clauses(dbresolver.Read)
}

func (t taskDo) WriteDB() ITaskDo {
	return t.Clauses(dbresolver.Write)
}

func (t taskDo) Session(config *gorm.Session) ITaskDo {
	return t.withDO(t.DO.Session(config))
}

func (t taskDo) Clauses(conds ...clause.Expression) ITaskDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t taskDo) Returning(value interface{}, columns ...string) ITaskDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t taskDo) Not(conds ...gen.Condition) ITaskDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t taskDo) Or(conds ...gen.Condition) ITaskDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t taskDo) Select(conds ...field.Expr) ITaskDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t taskDo) Where(conds ...gen.Condition) ITaskDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t taskDo) Order(conds ...field.Expr) ITaskDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t taskDo) Distinct(cols ...field.Expr) ITaskDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t taskDo) Omit(cols ...field.Expr) ITaskDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t taskDo) Join(table schema.Tabler, on ...field.Expr) ITaskDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t taskDo) LeftJoin(table schema.Tabler, on ...field.Expr) ITaskDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t taskDo) RightJoin(table schema.Tabler, on ...field.Expr) ITaskDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t taskDo) Group(cols ...field.Expr) ITaskDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t taskDo) Having(conds ...gen.Condition) ITaskDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t taskDo) Limit(limit int) ITaskDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t taskDo) Offset(offset int) ITaskDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t taskDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ITaskDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t taskDo) Unscoped() ITaskDo {
	return t.withDO(t.DO.Unscoped())
}

func (t taskDo) Create(values ...*core.Task) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t taskDo) CreateInBatches(values []*core.Task, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t taskDo) Save(values ...*core.Task) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t taskDo) First() (*core.Task, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*core.Task), nil
	}
}

func (t taskDo) Take() (*core.Task, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*core.Task), nil
	}
}

func (t taskDo) Last() (*core.Task, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*core.Task), nil
	}
}

func (t taskDo) Find() ([]*core.Task, error) {
	result, err := t.DO.Find()
	return result.([]*core.Task), err
}

func (t taskDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*core.Task, err error) {
	buf := make([]*core.Task, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t taskDo) FindInBatches(result *[]*core.Task, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t taskDo) Attrs(attrs ...field.AssignExpr) ITaskDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t taskDo) Assign(attrs ...field.AssignExpr) ITaskDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t taskDo) Joins(fields ...field.RelationField) ITaskDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t taskDo) Preload(fields ...field.RelationField) ITaskDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t taskDo) FirstOrInit() (*core.Task, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*core.Task), nil
	}
}

func (t taskDo) FirstOrCreate() (*core.Task, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*core.Task), nil
	}
}

func (t taskDo) FindByPage(offset int, limit int) (result []*core.Task, count int64, err error) {
	result, err = t.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = t.Offset(-1).Limit(-1).Count()
	return
}

func (t taskDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t taskDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t taskDo) Delete(models ...*core.Task) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *taskDo) withDO(do gen.Dao) *taskDo {
	t.DO = *do.(*gen.DO)
	return t
}
